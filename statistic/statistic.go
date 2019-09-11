// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package statistic

import (
    "citron/fileinfo"
    "fmt"
    "sync"
    "sync/atomic"
    "time"
)

const (
    B = 1
    KB = 1024
    MB = 1024 * KB
    GB = 1024 * MB
)

type Listener interface {
    OnWrite(int64, int64)
}

type Statistic struct {
    startTime time.Time
    totalSize int64
    totalFile int64
    readSize  int64
    writeSize int64
    failList  []fileinfo.FileInfo
    mutex     sync.Mutex
    listeners []Listener
}

func New() *Statistic {
    return &Statistic{
        startTime: time.Now(),
    }
}

func (s *Statistic) ResetTime() {
    s.startTime = time.Now()
}

func (s *Statistic) AddListener(l Listener) {
    s.listeners = append(s.listeners, l)
}

func (s *Statistic) notifyWrite() {
    for _, l := range s.listeners {
        l.OnWrite(s.writeSize, s.totalSize)
    }
}

func (s *Statistic) AddTotalSize(delta int64) int64 {
    return logAdd("AddTotalSize", &s.totalSize, delta)
}

func (s *Statistic) AddFileCount(delta int64) int64 {
    return logAdd("AddFileCount", &s.totalFile, delta)
}

func (s *Statistic) AddReadSize(delta int64) int64 {
    ret := atomic.AddInt64(&s.readSize, delta)
    return ret
}

func (s *Statistic) AddWriteSize(delta int64) int64 {
    ret := atomic.AddInt64(&s.writeSize, delta)
    s.notifyWrite()
    return ret
}

func (s *Statistic) ReadSize() int64 {
    return s.readSize
}

func (s *Statistic) WriteSize() int64 {
    return s.writeSize
}

func (s *Statistic) TotalSize() int64 {
    return s.totalSize
}

func logAdd(tag string, addr *int64, delta int64) int64 {
    ret := atomic.AddInt64(addr, delta)
    //log.Debug("%s: size %d", tag, ret)
    return ret
}

func (s *Statistic) String() string {
    return fmt.Sprintf("Statistic info - Total Size: %d , Total File Count: %d , Read Rate: %.2f MB/S Write Rate: %.2f MB/S, Read: %d , Write: %d , Use time: %d ms",
        s.totalSize, s.totalFile, s.ReadRate(MB, time.Second),s.WriteRate(MB, time.Second), s.readSize, s.writeSize, time.Since(s.startTime)/time.Millisecond)
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

//每毫秒读取速度
func (s *Statistic) ReadRate(sizeMeasure int64, timeMeasure time.Duration) float64 {
    if sizeMeasure <= 0 {
        sizeMeasure = 1
    }
    if timeMeasure <= 0 {
        timeMeasure = 1
    }
    useTime := time.Since(s.startTime)
    return float64(s.readSize) / float64(useTime) * float64(timeMeasure) / float64(sizeMeasure)
}

//每毫秒写入速度
func (s *Statistic) WriteRate(sizeMeasure int64, timeMeasure time.Duration) float64 {
    if sizeMeasure <= 0 {
        sizeMeasure = 1
    }
    if timeMeasure <= 0 {
        timeMeasure = 1
    }
    useTime := time.Since(s.startTime)
    return float64(s.writeSize) / float64(useTime) * float64(timeMeasure) / float64(sizeMeasure)
}

//预计完成时间(单位：秒)
func (s *Statistic) PredictTime() int64 {
    var (
        rTime int64 = 0
        wTime int64 = 0
    )
    rRate := s.ReadRate(B, time.Second)
    wRate := s.WriteRate(B, time.Second)
    if rRate == 0 {
        rTime = -1
    } else {
        rTime = int64(float64(s.totalSize - s.readSize) / rRate)
    }

    if wRate == 0 {
        wTime = -1
    } else {
        wTime = int64(float64(s.totalSize - s.writeSize) / wRate)
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
