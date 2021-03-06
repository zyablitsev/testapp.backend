package model

import "sync"

var (
	once     sync.Once
	instance *Data
)

func GetDataInstance() *Data {
	once.Do(func() {
		instance = new(Data)
		instance.users = make(map[int]map[int]int8)
	})
	return instance
}
