// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package statistic

import (
    "math"
    "time"
)

type Limiter struct {
    interval time.Duration

    rate float64

    start time.Time

    statistic *Statistic
}

func NewLimiter(st *Statistic, rate int64, interval time.Duration) *Limiter {
    return &Limiter{
        interval: interval,
        rate:     calcRate(rate, interval),
        start:    time.Now(),
        statistic: st,
    }
}

func calcRate(count int64, t time.Duration) float64 {
    return float64(count) * (float64(time.Second) / float64(t))
}

func (l *Limiter) ResetTime() {
    l.start = time.Now()
}

func (l *Limiter) Check(startTime time.Time, count int64) time.Duration {
    if l.rate == math.MaxInt64 {
        return 0
    }

    now := time.Since(startTime)
    if now == 0 {
        now = 1
    }
    var t time.Duration
    if now < l.interval {
        t = now
    } else {
        t = now - (now % l.interval)
    }

    rate := float64(count) / (float64(t) / float64(time.Second))

    if rate <= l.rate {
        return 0
    }

    return time.Duration(float64(count)/l.rate)*time.Second - t
}

func (l *Limiter) OnRead(size int64) {
    count := l.statistic.AddReadSize(size)
    waitTime := l.Check(l.start, count)
    if waitTime > 0 {
        time.Sleep(waitTime)
    }
}

func (l *Limiter) OnWrite(size int64) {
    l.statistic.AddWriteSize(size)
}
