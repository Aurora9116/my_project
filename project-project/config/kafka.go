package config

import (
	"context"
	"go.uber.org/zap"
	"test.com/project-common/kk"
)

var kw *kk.KafkaWriter

func InitKafkaWriter() func() {
	kw = kk.GetWriter("localhost:9092")
	return kw.Close
}
func SendLog(data []byte) {
	kw.Send(kk.LogData{
		Data:  data,
		Topic: "msproject_log",
	})
}

type KafkaCache struct {
	R *kk.KafkaReader
}

func (c *KafkaCache) DeleteCache() {
	for {
		message, err := c.R.R.ReadMessage(context.Background())
		if err != nil {
			zap.L().Error("DeleteCache err", zap.Error(err))
			continue
		}
		if string(message.Value) == "task" {

		}
	}
}

func NewCacheReader() *KafkaCache {
	reader := kk.GetReader([]string{"localhost:9092"}, "cache_group", "msproject_cache")
	return &KafkaCache{R: reader}
}
