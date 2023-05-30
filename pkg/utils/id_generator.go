package utils

import (
	"errors"
	"github.com/sony/sonyflake"
)

var generator *sonyflake.Sonyflake

func Init() {
	settings := sonyflake.Settings{}
	generator = sonyflake.NewSonyflake(settings)
}

func GenerateID() (int64, error) {
	if generator == nil {
		return 0, errors.New("generator not init")
	}
	id, err := generator.NextID()
	if err != nil {
		return 0, err
	}
	return int64(id), nil
}
