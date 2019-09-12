// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package filter

import (
    "citron/fileinfo"
    "github.com/xfali/goutils/io"
    "github.com/xfali/goutils/log"
    "os"
)

func KeepDelFiler(info fileinfo.FileInfo, fc FilterChain) error {
    if info.State != fileinfo.Deleted {
        return fc.Filter(info)
    }
    return nil
}

func RmSourceFilter(info fileinfo.FileInfo, fc FilterChain) error {
    err := fc.Filter(info)
    if err != nil {
        return err
    }

    path := info.FilePath
    if io.IsPathExists(path) {
        if info.IsDir {
            //FIXME: 由于可能先发送当前目录，如果目录被先删除，则该目录下的所有文件将不能被正确备份。
            //需要理清文件和目录的备份策略：是否需要发送目录，还是单纯发送文件。
            // 发送目录的好处：当目录被删除时，备份仓库可以一次性删除整个目录，而不用先删除目录下的文件，然后判断是否需要同时删除目录
            // 坏处：会残留空目录，如这个场景，源目录会残留，需要特殊处理
            //err := os.RemoveAll(path)
            //if err != nil {
            //    return err
            //}
        } else {
            err := os.Remove(path)
            if err != nil {
                return err
            }
        }
    } else {
        log.Info("file not found %s", path)
    }
    return nil
}
