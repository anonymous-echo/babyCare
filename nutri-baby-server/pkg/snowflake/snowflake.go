package snowflake

import (
	"sync"

	"github.com/bwmarrin/snowflake"
)

var (
	node *snowflake.Node
	once sync.Once
)

// Init 初始化雪花算法生成器
// nodeID: 节点ID (0-1023), 用于分布式环境区分不同服务实例
func Init(nodeID int64) error {
	var err error
	once.Do(func() {
		node, err = snowflake.NewNode(nodeID)
	})
	return err
}

// Generate 生成雪花ID
func Generate() int64 {
	if node == nil {
		// 如果未初始化,使用默认节点ID 1
		_ = Init(1)
	}
	return node.Generate().Int64()
}

// GenerateString 生成字符串格式的雪花ID
func GenerateString() string {
	return node.Generate().String()
}

// ParseID 解析雪花ID,返回时间戳等信息
func ParseID(id int64) snowflake.ID {
	return snowflake.ID(id)
}
