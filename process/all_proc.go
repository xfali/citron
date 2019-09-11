// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package process

import (
    "citron/fileinfo"
    "citron/io"
    "citron/store"
)

//全量备份
func allProcess(srcDir string, store store.MetaStore, proc ProcFunc) error {
    files, err := io.GetDirFiles(srcDir)
    if err != nil {
        return err
    }

    var result = files

    errDiff := proc(result)
    if errDiff != nil {
        return errDiff
    }

    for _, dir := range result {
        if dir.IsDir {
            if dir.State == fileinfo.Create || dir.State == fileinfo.Modified {
                err := allProcess(dir.FilePath, store, proc)
                if err != nil {
                    return err
                }
            }
        }
    }
    return nil
}
