//+----------------------------------------------------------------------
// | Copyright (c) 2024 http://www.vuecmf.com All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( https://github.com/vuecmf/vuecmf-go/blob/master/LICENSE )
// +----------------------------------------------------------------------
// | Author: tulihua2004@126.com
// +----------------------------------------------------------------------

/*
//使用示例
//app.Cache().Set("hello", "123456")

login := &model.LoginForm{
	Username: "haha",
	Password: "123456",
}

_ = app.Cache().Set("user", login)

var loginRes model.LoginForm
_ = app.Cache().Get("user", &loginRes)

fmt.Println("loginRes = ", loginRes)

var str string
app.Cache().Get("hello", &str)
fmt.Println("str cache = ", str)*/

package app

import (
	"encoding/json"
	"github.com/allegro/bigcache/v3"
	"log"
	"sync"
	"time"
)

var bc *bigcache.BigCache

func init() {
	if bc == nil {
		config := bigcache.Config{
			// number of shards (must be a power of 2)
			Shards: 1024,

			// time after which entry can be evicted
			LifeWindow: 120 * time.Minute,

			// Interval between removing expired entries (clean up).
			// If set to <= 0 then no action is performed.
			// Setting to < 1 second is counterproductive — bigcache has a one second resolution.
			CleanWindow: 5 * time.Minute,

			// rps * lifeWindow, used only in initial memory allocation
			MaxEntriesInWindow: 1000 * 10 * 60,

			// max entry size in bytes, used only in initial memory allocation
			MaxEntrySize: 500,

			// prints information about additional memory allocation
			Verbose: true,

			// cache will not allocate more memory than this limit, value in MB
			// if value is reached then the oldest entries can be overridden for the new ones
			// 0 value means no size limit
			HardMaxCacheSize: 8192,

			// callback fired when the oldest entry is removed because of its expiration time or no space left
			// for the new entry, or because delete was called. A bitmask representing the reason will be returned.
			// Default value is nil which means no callback and it prevents from unwrapping the oldest entry.
			OnRemove: nil,

			// OnRemoveWithReason is a callback fired when the oldest entry is removed because of its expiration time or no space left
			// for the new entry, or because delete was called. A constant representing the reason will be passed through.
			// Default value is nil which means no callback and it prevents from unwrapping the oldest entry.
			// Ignored if OnRemove is specified.
			OnRemoveWithReason: nil,
		}
		var err error
		bc, err = bigcache.NewBigCache(config)
		if err != nil {
			log.Fatal(err)
		}
	}
}

type cache struct {
}

// Set 添加缓存
//
//	参数：
//		key 缓存的键
//		content 缓存的内容
func (c *cache) Set(key string, content interface{}) error {
	cb, err := json.Marshal(content)
	if err != nil {
		return err
	}
	err = bc.Set(key, cb)
	return err
}

// Get 获取缓存内容
//
//	参数：
//		key 缓存的key
//		res 存放缓存内容的容器
func (c *cache) Get(key string, res interface{}) error {
	rb, err := bc.Get(key)
	if err != nil {
		return err
	}
	err = json.Unmarshal(rb, res)
	return err
}

// Del 删除缓存内容
//
//	参数：
//		key 缓存的key
func (c *cache) Del(key string) error {
	return bc.Delete(key)
}

var cacheOnce sync.Once
var c *cache

// Cache 获取缓存组件实例
func Cache() *cache {
	cacheOnce.Do(func() {
		c = &cache{}
	})

	return c
}
