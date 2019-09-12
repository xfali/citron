// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package fileinfo

import (
    "citron/uri"
    "github.com/xfali/goutils/log"
    "time"
)

const (
    Create = iota
    Modified
    Deleted
)

const (
    MD5  = "MD5"
    SHA1 = "SHA1"
)

type FileInfo struct {
    FileName string `json:"filename"`
    FilePath string `json:"filepath"`
    Parent   string `json:"parent"`

    From   uri.URI `json:"from"`
    To     uri.URI `json:"to"`
    Hidden bool `json:"hidden,omitempty"`

    State int `json:"state"`

    IsDir bool  `json:"isDir"`
    Size  int64 `json:"size"`

    ModTime time.Time `json:"modTime"`

    Checksum     string `json:"checksum,omitempty"`
    ChecksumType string `json:"checksumType,omitempty"`
}

func (f *FileInfo) Process(other FileInfo) FileInfo {
    if f.FilePath != other.FilePath {
        log.Fatal("%s is not match %s", f.FilePath, other.FilePath)
    }

    if f.Empty() && !other.Empty() {
        ret := other
        ret.State = Create
        log.Debug("new file %s", ret.FilePath)
        return ret
    }

    if !f.Empty() && other.Empty() {
        ret := *f
        ret.State = Deleted
        log.Debug("delete file %s", ret.FilePath)
        return ret
    }

    if !f.Empty() && !other.Empty() {
        //1、文件大小发生变化，直接备份
        //2、修改时间相同，不备份
        //3、修改时间不相同：a、如果未开启文件校验，则直接备份；b、如果开启文件检验，则看文件是否校验码相等
        if f.Size != other.Size || (!f.ModTime.Equal(other.ModTime) && !f.checksumEqual(other)) {
            ret := other
            ret.State = Modified
            log.Debug("modify file %s", ret.FilePath)
            return ret
        }
    }

    return FileInfo{}
}

func (f *FileInfo) checksumEqual(other FileInfo) bool {
    if f.IsDir {
        return true
    }
    if f.ChecksumType != "" {
        if other.ChecksumType == "" {
            return false
        } else {
            if f.ChecksumType != other.ChecksumType {
                return false
            }
        }
    }

    if f.Checksum != "" {
        if other.Checksum == "" {
            return false
        } else {
            if f.Checksum == other.Checksum {
                return true
            }
        }
    }
    return false
}

func (f *FileInfo) Empty() bool {
    return f.FilePath == ""
}
