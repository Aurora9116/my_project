package kk

import (
	"encoding/json"
	"testing"
	"time"
)

func TestProducer(t *testing.T) {
	w := GetWriter("localhost:9092")
	m := make(map[string]string)
	m["projectCode"] = "120w"
	bytes, _ := json.Marshal(m)
	w.Send(LogData{
		Data:  bytes,
		Topic: "log",
	})
	time.Sleep(time.Second * 2)
}
func TestConsumer(t *testing.T) {
	GetReader([]string{"localhost:9092"}, "group1", "log")
	for {

	}
}
