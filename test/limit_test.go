// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package test

import (
    "citron/statistic"
    "testing"
    "time"
)

func TestLimiter(t *testing.T) {
    t.Run("1 per second", func(t *testing.T) {
        l := statistic.NewLimiter(1, time.Second)

        now := time.Now()
        var i int64 = 0
        for ; i < 10; i++ {
            st := l.Check(now, i)
            if st > 0 {
                t.Logf("i %d sleep %d\n", i, st)
                time.Sleep(st)
            }
        }
    })

    t.Run("10 per second", func(t *testing.T) {
        l := statistic.NewLimiter(10, time.Second)

        now := time.Now()
        var i int64 = 0
        for ; i < 11; i++ {
            st := l.Check(now, i)
            if st > 0 {
                t.Logf("i %d sleep %d\n", i, st)
                time.Sleep(st)
            }
        }
    })

    t.Run("5 per millisecond", func(t *testing.T) {
        l := statistic.NewLimiter(5, time.Millisecond)

        now := time.Now()
        var i int64 = 0
        for ; i < 51000; i++ {
            st := l.Check(now, i)
            if st > 0 {
                t.Logf("i %d sleep %d\n", i, st)
                time.Sleep(st)
            }
        }
    })
}
