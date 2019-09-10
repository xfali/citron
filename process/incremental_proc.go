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

//增量备份
func incrementalProcess(rootDir, srcDir string, store store.MetaStore, proc ProcFunc) error {
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

    allInfo, err := store.Query()
    if err != nil {
        return err
    }
    var result []fileinfo.FileInfo
    findDiffFiles(allInfo, files, &result, false)
    findDiffFiles(files, allInfo, &result, true)

    errDiff := proc(rel, result)
    if errDiff != nil {
        return errDiff
    }

    for _, dir := range files {
        if dir.IsDir {
            if dir.State == fileinfo.Create || dir.State == fileinfo.Modified {
                err := incrementalProcess(rootDir, dir.FilePath, store, proc)
                if err != nil {
                    return err
                }
            }
        }
    }
    return nil
}
