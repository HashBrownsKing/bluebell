package redis

import (
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var VoteTimeExpire = 60 * 60 * 24 * 7 // 新帖子一周内可进行投票
var (
	ErrorVoteTimeExpire = errors.New("投票时间已过")
	ErrorPostNotExist   = errors.New("帖子在redis中不存在")
)

/*
帖子排序的原则是越新、得分越高的帖子越上面
以创建帖子的时间戳为帖子的初始得分，多一票赞成+432分，多一票反对-432分
一天有86400s, 432分每票表示前一天创建的帖子要想分数高于下一天同一时间创建的帖子需要获得200个支持
*/
const ScorePerVote float64 = 432 // 投一票加432分

func VoteForPost(userID, postID string, voteInfo float64) (err error) {
	// 1. 获取帖子创发布时间并检查是否超过一周
	postCreateTime, err := rdb.ZScore(ctx, GetRedisKey(KeyPostTimeZSet), postID).Result()
	if err != nil {
		if err == redis.Nil {
			zap.L().Error("帖子不存在", zap.String("postID", postID))
			return ErrorPostNotExist
		}
		zap.L().Error("查询帖子创建时间失败", zap.Error(err))
		return err
	}
	if time.Now().Unix()-int64(postCreateTime) > int64(VoteTimeExpire) {
		return ErrorVoteTimeExpire
	}

	// 2. 查询当前用户给当前帖子的投票记录
	voteRecord, err := rdb.ZScore(ctx, GetRedisKey(KeyUserVoteRecordZSetPrefix+postID), userID).Result()
	if err != nil && err != redis.Nil { // 非空值错误才返回
		zap.L().Error("查询用户历史投票失败", zap.Error(err))
		return err
	}
	// 如果用户没投过票，voteRecord = 0
	if err == redis.Nil {
		voteRecord = 0
	}
	if voteInfo == voteRecord { // 投票状态没有改变直接返回
		return nil
	}

	// 3. 更新帖子分数
	scoreDiff := (voteInfo - voteRecord) * ScorePerVote // 计算分数差值
	pipeline := rdb.TxPipeline()
	pipeline.ZIncrBy(ctx, GetRedisKey(KeyPostScoreZSet), scoreDiff, postID)

	// 4.更新用户的投票数据
	if voteInfo == 0 {
		pipeline.ZRem(ctx, GetRedisKey(KeyUserVoteRecordZSetPrefix+postID), userID)
	} else {
		pipeline.ZAdd(ctx, GetRedisKey(KeyUserVoteRecordZSetPrefix+postID), redis.Z{
			Score:  voteInfo,
			Member: userID,
		})
	}

	_, err = pipeline.Exec(ctx)
	if err != nil {
		zap.L().Error("redis.VoteForPost failed: ", zap.Error(err))
	}
	return
}
