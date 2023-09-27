package transfer

import (
	"encoding/json"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/mq"
)

type Task struct {
	storageID uint64
}

func (t Task) transfer() error {
	data, err := json.Marshal(&t)
	if err != nil {
		return err
	}

	return mq.Publish("", "transfer-task", data)
}
