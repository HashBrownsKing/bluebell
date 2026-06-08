package redis

import (
	"bluebell/models"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func CreatePost(postID, communityID string) (err error) {
	pipeline := rdb.TxPipeline() // 开启一个事务
	// 设置创建时间
	pipeline.ZAdd(ctx, GetRedisKey(KeyPostTimeZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})
	// 设置初始分数
	pipeline.ZAdd(ctx, GetRedisKey(KeyPostScoreZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})
	pipeline.SAdd(ctx, GetRedisKey(KeyCommunitySetPrefix+communityID), postID)
	_, err = pipeline.Exec(ctx)
	return err
}

func getIDsFormKey(page, size int64, key string) (ids []string, err error) {
	start := (page - 1) * size
	stop := start + size - 1
	ids = make([]string, 0, size)
	ids, err = rdb.ZRangeArgs(ctx, redis.ZRangeArgs{
		Key:   key,
		Start: start,
		Stop:  stop,
		Rev:   true,
	}).Result()
	if err != nil {
		zap.L().Error("rdb.ZRangeArgs failed", zap.Error(err))
		return nil, err
	}
	return ids, nil
}

func GetPostIDInOrder(p *models.ParamPostList) (ids []string, err error) {
	key := GetRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		key = GetRedisKey(KeyPostScoreZSet)
	}
	return getIDsFormKey(p.Page, p.Size, key)
}

func GetPostVoteNum(ids []string) (voteNum []int64, err error) {
	pipeline := rdb.Pipeline()
	for _, id := range ids {
		key := GetRedisKey(KeyUserVoteRecordZSetPrefix + id)
		// 查找分数在1-1范围内的个数
		pipeline.ZCount(ctx, key, "1", "1")
	}
	cmders, err := pipeline.Exec(ctx)
	if err != nil {
		zap.L().Error("redis.GetPostVoteNum failed: ", zap.Error(err))
		return nil, err
	}
	voteNum = make([]int64, 0, len(ids))
	for _, cmder := range cmders {
		num, err := cmder.(*redis.IntCmd).Result()
		if err != nil {
			zap.L().Error("解析投票数失败", zap.Error(err))
			voteNum = append(voteNum, 0)
			continue
		}
		voteNum = append(voteNum, num)
	}
	return voteNum, nil
}

func GetCommunityPostIDInOrder(p *models.ParamPostList) (ids []string, err error) {
	// 社区帖子集合key
	communityKey := GetRedisKey(KeyCommunitySetPrefix + strconv.FormatInt(p.CommunityID, 10))
	// 排序key
	orderKey := GetRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		orderKey = GetRedisKey(KeyPostScoreZSet)
	}
	// 创建临时缓存key
	cacheKey := orderKey + strconv.FormatInt(p.CommunityID, 10)
	// 如果缓存不存在
	if rdb.Exists(ctx, cacheKey).Val() < 1 {
		pipeline := rdb.Pipeline()
		// 将communityKey和orderKey对应的数据取交集存入cacheKey
		pipeline.ZInterStore(ctx, cacheKey, &redis.ZStore{
			Aggregate: "MAX",
			Keys:      []string{communityKey, orderKey},
		})
		pipeline.Expire(ctx, cacheKey, 60*time.Second)
		_, err := pipeline.Exec(ctx)
		if err != nil {
			return nil, err
		}
	}
	return getIDsFormKey(p.Page, p.Size, cacheKey)
}
