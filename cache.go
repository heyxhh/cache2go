/*
 * Simple caching library with expiration capabilities
 *     Copyright (c) 2012, Radu Ioan Fericean
 *                   2013-2017, Christian Muehlhaeuser <muesli@gmail.com>
 *
 *   For license see LICENSE.txt
 */

package cache2go

import (
	"sync"
)

var (
	cache = make(map[string]*CacheTable) // cache记录着所有的CacheTable
	mutex sync.RWMutex                   // mutex是cache的读写锁
)

// Cache returns the existing cache table with given name or creates a new one
// if the table does not exist yet.
func Cache(table string) *CacheTable {
	mutex.RLock()
	t, ok := cache[table]
	mutex.RUnlock()

	if !ok {
		mutex.Lock()
		t, ok = cache[table] // 注意: 加锁成功后必须要二次检查table是否已存在；因为有可能被别的并发Cache函数创建了CacheTable
		// Double check whether the table exists or not.
		if !ok {
			t = &CacheTable{
				name:  table,
				items: make(map[interface{}]*CacheItem),
			}
			cache[table] = t // 注意t是指针，若传递的是值类型，相当于再次复制了一份。而且后续修改数据时也要获取指针。
		}
		mutex.Unlock()
	}

	return t
}
