// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package store

import (
    "fbt/fileinfo"
)

type MetaStore interface {
    Open(url string, dataDir string) error

    Insert(info fileinfo.FileInfo) error

    Update(info fileinfo.FileInfo) error

    Read(uri string) error

    Query() ([]fileinfo.FileInfo, error)

    QueryByPath(uri string) (fileinfo.FileInfo, error)

    Delete(info fileinfo.FileInfo) error

    Save() error

    Close() error
}

func SaveMeta(store MetaStore, info fileinfo.FileInfo) error {
    switch info.State {
    case fileinfo.Create:
        return store.Insert(info)
    case fileinfo.Modified:
        return store.Update(info)
    case fileinfo.Deleted:
        return store.Delete(info)
    default:
        break
    }
    return nil
}
