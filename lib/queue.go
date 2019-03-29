package lib

import (
	"sync"
	"time"
)

// Queue Queue
type Queue struct {
	data   []byte
	length int
	sync.Mutex
}

// QueueInstance get instance
func QueueInstance() *Queue {
	queue := &Queue{
		data:   make([]byte, 0),
		length: 0,
	}
	return queue
}

// Push Push
func (q *Queue) Push(data []byte) {
	q.Lock()
	defer q.Unlock()
	for _, b := range data {
		q.data = append(q.data, b)
		q.length++
	}
}

// Pop Pop getCount: 获取byte的数量，如果数量不足，将进入阻塞
func (q *Queue) Pop(getCount int) []byte {
	result := make([]byte, getCount)
	for {
		if getCount < q.length {
			q.Lock()
			result = q.data[:getCount]
			q.data = q.data[getCount:]
			q.length = q.length - getCount
			q.Unlock()
			break
		} else {
			time.Sleep(1e8) // 阻塞时间请自定
		}
	}
	return result
}
