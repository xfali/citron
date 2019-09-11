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

//增量备份
func incrementalProcess(srcDir string, store store.MetaStore, proc ProcFunc) error {
    files, err := io.GetDirFiles(srcDir)
    if err != nil {
        return err
    }

    allInfo, err := store.QueryByPath(srcDir)
    if err != nil {
        return err
    }
    var result []fileinfo.FileInfo
    findDiffFiles(allInfo, files, &result, false)
    findDiffFiles(files, allInfo, &result, true)

    errDiff := proc(result)
    if errDiff != nil {
        return errDiff
    }

    for _, dir := range files {
        if dir.IsDir {
            if dir.State == fileinfo.Create || dir.State == fileinfo.Modified {
                err := incrementalProcess(dir.FilePath, store, proc)
                if err != nil {
                    return err
                }
            }
        }
    }
    return nil
}
