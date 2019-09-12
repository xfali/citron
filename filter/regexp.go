// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package filter

import (
    "citron/fileinfo"
    "fmt"
    "github.com/xfali/goutils/log"
    "regexp"
)

var regexp_exp = map[string] string {
    `以.go后缀结尾` : `^\S+\.go$`,
    `包含go` : `^\S+go\S+$`,
    `不以.go后缀结尾` : `^((?!\.go$).)*$`,
    `所有` : `^\S+$`,
}

type reg regexp.Regexp

func PrintRegexp() {
    for k, v := range regexp_exp {
        fmt.Printf("%s :\n", k)
        fmt.Printf("\t%s\n", v)
    }
}

func NewRegexp(regstr string) *reg {
    return (*reg)(regexp.MustCompile(regstr))
}

func (reg *reg) HideFiler(info fileinfo.FileInfo, fc FilterChain) error {
    if (*regexp.Regexp)(reg).MatchString(info.FilePath) {
        log.Info("Match! set %s hidden", info.To)
        info.Hidden = true
    }
    err := fc.Filter(info)
    if err != nil {
        return err
    }
    return nil
}

func (reg *reg) BackupFiler(info fileinfo.FileInfo, fc FilterChain) error {
    if (*regexp.Regexp)(reg).MatchString(info.FilePath) {
        log.Info("Match! backup %s", info.From)
        return fc.Filter(info)
    }
    return nil
}
