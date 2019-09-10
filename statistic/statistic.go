// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package statistic

import (
    "fbt/fileinfo"
    "github.com/xfali/goutils/log"
    "sync"
    "sync/atomic"
    "time"
)

type Statistic struct {
    startTime time.Time
    totalSize int64
    totalFile int64
    readSize  int64
    writeSize int64
    failList  []fileinfo.FileInfo
    mutex     sync.Mutex
}

func New() *Statistic {
    return &Statistic{
        startTime: time.Now(),
    }
}

func (s *Statistic) ResetTime() {
    s.startTime = time.Now()
}

func (s *Statistic) AddTotalSize(delta int64) int64 {
    return logAdd("AddTotalSize", &s.totalSize, delta)
}

func (s *Statistic) AddFileCount(delta int64) int64 {
    return logAdd("AddFileCount", &s.totalFile, delta)
}

func (s *Statistic) AddReadSize(delta int64) int64 {
    return logAdd("AddReadSize", &s.readSize, delta)
}

func (s *Statistic) AddWriteSize(delta int64) int64 {
    return logAdd("AddWriteSize", &s.writeSize, delta)
}

func logAdd(tag string, addr *int64, delta int64) int64 {
    ret := atomic.AddInt64(addr, delta)
    log.Debug("%s: size %d", tag, ret)
    return ret
}

func (s *Statistic) AddFailedFile(file fileinfo.FileInfo) {
    s.mutex.Lock()
    defer s.mutex.Unlock()

    s.failList = append(s.failList, file)
}

func (s *Statistic) GetFailedFile() []fileinfo.FileInfo {
    s.mutex.Lock()
    defer s.mutex.Unlock()

    return s.failList
}

//每秒读取速度
func (s *Statistic) ReadRate() int64 {
    return s.readSize / (int64(time.Since(s.startTime)) / int64(time.Second))
}

//每秒写入速度
func (s *Statistic) WriteRate() int64 {
    return s.writeSize / (int64(time.Since(s.startTime)) / int64(time.Second))
}

//预计完成时间(单位：秒)
func (s *Statistic) PredictTime() int64 {
    rRate := s.ReadRate()
    wRate := s.WriteRate()
    var (
        rTime int64 = 0
        wTime int64 = 0
    )
    if rRate == 0 {
        rTime = -1
    } else {
        rTime = (s.totalSize - s.readSize) / rRate
    }
    if wRate == 0 {
        wTime = -1
    } else {
        wTime = (s.totalSize - s.writeSize) / wRate
    }

    if rTime == -1 && wTime == -1 {
        return -1
    }
    if rTime > wTime {
        return rTime
    } else {
        return wTime
    }
}
