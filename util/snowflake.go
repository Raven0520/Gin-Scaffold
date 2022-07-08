package util

import (
	"errors"
	"sync"
	"time"
)

/**
 * 配置
 */
const (
	epoch int64 = 1526285084373

	numWorkerBits = 10

	numSequenceBits = 12

	MaxWorkID = -1 ^ (-1 << numWorkerBits)

	MaxSequence = -1 ^ (-1 << numSequenceBits)
)

type SnowFlake struct {
	lastTimestamp uint64
	sequence      uint32
	workerID      uint32
	lock          sync.Mutex // 互斥锁
}

func (sf *SnowFlake) pack() uint64 {
	uuid := (sf.lastTimestamp << (numWorkerBits + numSequenceBits)) | (uint64(sf.workerID) << numSequenceBits) | (uint64(sf.sequence))
	return uuid
}

// NewSnowFlake returns a new snowflake node that can be used to generate snowflake
func NewSnowFlake(workerID uint32) (*SnowFlake, error) {
	if workerID <= 0 || workerID > MaxWorkID {
		return nil, errors.New("InvalidWorkerId") // 无效的进程ID
	}
	return &SnowFlake{workerID: workerID}, nil
}

// Generate Next creates and returns a unique snowflake ID
func (sf *SnowFlake) Generate() (uint64, error) {
	// 锁处理
	sf.lock.Lock()
	defer sf.lock.Unlock()

	ts := timestamp()
	if ts == sf.lastTimestamp {
		sf.sequence = (sf.sequence + 1) & MaxSequence
		if sf.sequence == 0 {
			ts = sf.waitNextMilli(ts)
		}
	} else {
		sf.sequence = 0
	}

	if ts < sf.lastTimestamp {
		return 0, errors.New("InvalidSystemClock") // 锁失效
	}

	sf.lastTimestamp = ts
	return sf.pack(), nil
}

// waitNextMilli if that microsecond is full
// wait for the next microsecond
func (sf *SnowFlake) waitNextMilli(ts uint64) uint64 {
	for ts == sf.lastTimestamp {
		time.Sleep(100 * time.Microsecond)
		ts = timestamp()
	}
	return ts
}

// timestamp
func timestamp() uint64 {
	return uint64(time.Now().UnixNano()/int64(1000000) - epoch)
}
