package list

import "sync"

var instance *WhiteList
var once sync.Once

func GetInstance() *WhiteList {
	once.Do(func() {
		instance = &WhiteList{}
	})
	return instance
}
