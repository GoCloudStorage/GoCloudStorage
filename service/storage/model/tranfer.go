package model

import (
	"encoding/json"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/mq"
)

type Task struct {
	StorageID uint64
}

func (t Task) Transfer() error {
	data, err := json.Marshal(&t)
	if err != nil {
		return err
	}

	return mq.Publish("", "transfer-task", data)
}
