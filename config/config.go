// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package config

import "encoding/json"

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

    Regexp
}

var GConfig Config

func (c *Config) String() string {
    b, err := json.Marshal(GConfig)
    if err != nil {
        return ""
    }
    return string(b)
}
