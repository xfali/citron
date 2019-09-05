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
)

//增量更新
func incrementalProcess(rootDir, srcDir string, trans transport.Transport, store store.MetaStore) error {
    err := store.Read(srcDir)
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

    rel := io.SubPath(srcDir, rootDir)

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
