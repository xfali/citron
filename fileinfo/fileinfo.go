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

    From uri.URI `json:"from"`
    To   uri.URI `json:"to"`

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
        if !f.ModTime.Equal(other.ModTime) || f.Size != other.Size {
            if f.checksumEqual(other) {
                return FileInfo{}
            }
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
