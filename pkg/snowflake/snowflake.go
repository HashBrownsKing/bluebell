package snowflake

import (
	"fmt"
	"time"

	"github.com/bwmarrin/snowflake"
)

var node *snowflake.Node // 单节点

func Init(startTime string, mechineID int64) (err error) {
	var st time.Time
	st, err = time.ParseInLocation("2006-01-02", startTime, time.Local)
	if err != nil {
		return fmt.Errorf("时间解析失败: %w", err)
	}
	// 设置雪花算法起始时间戳（毫秒级）
	snowflake.Epoch = st.UnixMilli()
	// 初始化机器节点
	node, err = snowflake.NewNode(mechineID)
	if err != nil {
		return fmt.Errorf("节点初始化失败: %w", err)
	}
	return
}

// GenID 生成全局唯一ID
func GetID() int64 {
	return node.Generate().Int64()
}
