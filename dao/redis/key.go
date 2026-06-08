package redis

// 定义投票功能相关的redis key

const (
	KeyPrefix                   = "bluebell:"
	KeyPostTimeZSet             = "post:time"       // zset; 帖子及发布时间
	KeyPostScoreZSet            = "post:score"      // zset; 帖子及投票分数
	KeyUserVoteRecordZSetPrefix = "post:voted:"     // zset; 存储给某篇帖子投票的用户和用户投票状态,还需要加上post_id
	KeyCommunitySetPrefix       = "community:post:" // set; 存储某个社区下的帖子，还需要加上community_id
)

func GetRedisKey(key string) string {
	return KeyPrefix + key
}
