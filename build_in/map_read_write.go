package main

import "sync"

// 博客 链接：https://www.jianshu.com/p/10a998089486

// Maps are not safe for concurrent use: it's not defined what happens when you read and write to them simultaneously.
// If you need to read from and write to a map from concurrently executing goroutines, the accesses must be mediated
// by some kind of synchronization mechanism. One common way to protect maps is with sync.RWMutex.

// After long discussion it was decided that the typical use of maps did not require safe access from multiple goroutines,
// and in those cases where it did, the map was probably part of some larger data structure or computation that
// was already synchronized. Therefore requiring that all map operations grab a mutex would slow down most programs
// and add safety to few. This was not an easy decision, however, since it means uncontrolled map access can crash the program.

// 大致意思就是说，并发访问map是不安全的，会出现未定义行为，导致程序退出。所以如果希望在多协程中并发访问map，
// 必须提供某种同步机制，一般情况下通过读写锁sync.RWMutex实现对map的并发访问控制，将map和sync.RWMutex封装一下，
// 可以实现对map的安全并发访问，示例代码如下：


type SafeMap struct {
	sync.RWMutex
	Map map[int]int
}

func main() {
	safeMap := newSafeMap(10)

	for i := 0; i < 10000; i++ {
		go safeMap.writeMap(i, i)
		go safeMap.readMap(i)
	}
}

func newSafeMap(size int) *SafeMap {
	sm := new(SafeMap)
	sm.Map = make(map[int]int, size)
	return sm
}

func (sm *SafeMap) readMap(key int) int {
	sm.RLock()
	value := sm.Map[key]
	sm.RUnlock()
	return value
}

func (sm *SafeMap) writeMap(key int, value int) {
	sm.Lock()
	sm.Map[key]	= value
	sm.Unlock()
}
