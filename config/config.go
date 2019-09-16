// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package config

import (
    "citron/statistic"
    "encoding/json"
    "github.com/xfali/goutils/log"
    "strconv"
    "strings"
    "time"
)

const (
    InfoDir = ".citronmeta"
)

type Regexp struct {
    RegexpHidden string
    RegexpBackup string
}

type Config struct {
    SourceDir    string
    DestUri      string
    ChecksumType string
    Incremental  bool
    NewRepo      bool
    MultiTaskNum int
    RmSrc        bool
    RmDel        bool
    Limit        string

    Regexp
}

var GConfig Config

func (c *Config) ParseLimit() (int64, time.Duration) {
    if c.Limit == "" {
        return 0, 0
    } else {
        strs := strings.Split(c.Limit, "/")
        first := strs[0]
        if first == "" {
            return 0, 0
        }

        m := strings.ToUpper(first[len(first)-1:])
        var rate int64 = 1
        if m == "M" {
            rate = int64(statistic.MB)
        } else if m == "K" {
            rate = int64(statistic.KB)
        } else if m == "G" {
            rate = int64(statistic.GB)
        }

        ret, err := strconv.ParseInt(first[:len(first)-1], 10, 64)
        if err != nil {
            log.Warn("parse limit error %s", err.Error())
            return 0, 0
        }

        if len(strs) < 2 {
            return ret * rate, time.Second
        }
        timeM := strings.ToUpper(strs[1])
        if timeM == "S" {
            return ret * rate, time.Second
        } else if timeM == "MS" {
            return ret * rate, time.Millisecond
        }
        return ret * rate, time.Second
    }
}

func (c *Config) String() string {
    b, err := json.Marshal(GConfig)
    if err != nil {
        return ""
    }
    return string(b)
}
