package model

import "sync"

var (
	once     sync.Once
	instance *Data
)

func GetDataInstance() *Data {
	once.Do(func() {
		instance = new(Data)
		instance.Records = make(map[int]map[[4]byte]bool)
	})
	return instance
}
