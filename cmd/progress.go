// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package cmd

import (
    "fbt/statistic"
    "fmt"
    "github.com/schollz/progressbar/v2"
    "sync/atomic"
    "time"
)

type Progress struct {
    current  int32
    bar      *progressbar.ProgressBar
    stopChan chan bool
    st       *statistic.Statistic
}

func NewProgress(st *statistic.Statistic) *Progress {
    r := &Progress{
        current: 0,
        bar: progressbar.NewOptions(
            100,
            progressbar.OptionSetDescription("备份进度"),
            progressbar.OptionEnableColorCodes(true),
            progressbar.OptionOnCompletion(func() {
                fmt.Println()
            })),
        stopChan: make(chan bool),
        st: st,
    }
    return r
}

func (p *Progress) Start() {
    go func() {
        timer := time.NewTicker(10 * time.Millisecond)
        for {
            select {
            case <-p.stopChan:
                timer.Stop()
                return
            case <-timer.C:
                p.move()
            }
        }
    }()
}

func (p *Progress) Stop() {
    close(p.stopChan)
    p.bar.Finish()
}

func (p *Progress) move() {
    cur := int32(p.st.WriteSize() * 100 / p.st.TotalSize())
    tmp := atomic.LoadInt32(&p.current)
    if cur > tmp {
       if atomic.CompareAndSwapInt32(&p.current, tmp, cur) {
           p.bar.Add(int(cur - tmp))
       }
    }
}

func (p *Progress) OnWrite(write int64, total int64) {
    //cur := int32(write * 100 / total)
    //tmp := atomic.LoadInt32(&p.current)
    //if cur > tmp {
    //    if atomic.CompareAndSwapInt32(&p.current, tmp, cur) {
    ////        p.bar.Add(int(cur - tmp))
    //    }
    //}
}

func (p *Progress) Finish() {
    p.bar.Finish()
}