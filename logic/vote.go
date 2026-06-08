package logic

import (
	"bluebell/dao/redis"
	"bluebell/models"
	"strconv"

	"go.uber.org/zap"
)

func VoteForPost(userID int64, p *models.ParamPostVote) (err error) {
	zap.L().Debug("VoteForPost",
		zap.Int64("userID", userID),
		zap.Int64("postID", p.PostID),
		zap.Int8("voteInfo", p.VoteInfo))
	// redis中的键都是字符串，zset中的分数是浮点数，提前进行类型转换
	userIDStr := strconv.Itoa(int(userID))
	postIDStr := strconv.Itoa(int(p.PostID))
	voteInfo := float64(p.VoteInfo)
	return redis.VoteForPost(userIDStr, postIDStr, voteInfo)
}
