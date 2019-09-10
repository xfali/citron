// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package process

import (
    "fbt/fileinfo"
    "fbt/io"
    "fbt/store"
    "path/filepath"
    "strings"
)

//全量备份
func allProcess(rootDir, srcDir string, store store.MetaStore, proc ProcFunc) error {
    rel := io.SubPath(srcDir, rootDir)
    metaFile := strings.Replace(rel, string(filepath.Separator), "_", -1)
    if metaFile == "" {
        metaFile = "root"
    }

    err := store.Read(metaFile)
    if err != nil {
        return err
    }

    files, err := io.GetDirFiles(srcDir)
    if err != nil {
        return err
    }

    var result = files

    errDiff := proc(rel, result)
    if errDiff != nil {
        return errDiff
    }

    for _, dir := range result {
        if dir.IsDir {
            if dir.State == fileinfo.Create || dir.State == fileinfo.Modified {
                err := allProcess(rootDir, dir.FilePath, store, proc)
                if err != nil {
                    return err
                }
            }
        }
    }
    return nil
}
