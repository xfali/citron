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
    "fbt/transport"
    "path/filepath"
    "strings"
)

//增量更新
func incrementalProcess(rootDir, srcDir string, trans transport.Transport, store store.MetaStore) error {
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
    process(allInfo, files, &result, false)
    process(files, allInfo, &result, true)

    errDiff := processDiff(rel, result, trans, store)
    if errDiff != nil {
        return errDiff
    }

    for _, dir := range files {
        if dir.IsDir {
            if dir.State == fileinfo.Create || dir.State == fileinfo.Modified {
                err := incrementalProcess(rootDir, dir.FilePath, trans, store)
                if err != nil {
                    return err
                }
            }
        }
    }
    return nil
}
