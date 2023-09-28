package model

import (
	"github.com/GoCloudstorage/GoCloudstorage/pkg/mq"
	"testing"
)

func TestTask(t *testing.T) {
	mq.Init("162.14.115.114:5672", "cill", "12345678")
	task := Task{StorageID: 196168787014189057}
	err := task.transfer()
	if err != nil {
		t.Fatal(err)
	}
}
