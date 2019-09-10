// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package process

import (
    "fbt/cmd"
    "fbt/config"
    "fbt/fileinfo"
    "fbt/io"
    "fbt/statistic"
    "fbt/store"
    "fbt/transport"
    "fbt/uri"
    "github.com/xfali/executor"
    "github.com/xfali/goutils/log"
    "path/filepath"
    "sync"
)

type ProcFunc func([]fileinfo.FileInfo) error

func Process(srcDir string, trans transport.Transport, store store.MetaStore, statis *statistic.Statistic) error {
    srcDir = filepath.Clean(srcDir)
    var diffFile []fileinfo.FileInfo
    statisticFunc := func(result []fileinfo.FileInfo) error {
        if len(result) > 0 {
            for _, info := range result {
                diffFile = append(diffFile, info)
                if !info.IsDir {
                    statis.AddFileCount(1)
                    statis.AddTotalSize(info.Size)
                }
            }
        }

        return nil
    }

    if config.GConfig.Incremental {
        err := incrementalProcess(srcDir, store, statisticFunc)
        if err != nil {
            return err
        }
    } else {
        err := allProcess(srcDir, store, statisticFunc)
        if err != nil {
            return err
        }
    }

    p := cmd.NewProgress(statis)
    p.Start()
    defer p.Stop()

    var err error
    if len(diffFile) > 0 {
        if config.GConfig.SyncTrans {
            err = syncProcessDiff(srcDir, diffFile, trans, store, statis)
        } else {
            err = asyncProcessDiff(srcDir, diffFile, trans, store, statis)
        }
        if err != nil {
            return err
        }
        store.Save()
    }

    return nil
}

func findDiffFiles(list1, list2 []fileinfo.FileInfo, result *[]fileinfo.FileInfo, reverse bool) {
    for _, info := range list1 {
        found := false
        for _, file := range list2 {
            if info.FilePath == file.FilePath {
                ret := info.Process(file)
                if !ret.Empty() && !reverse {
                    *result = append(*result, ret)
                }
                found = true
                break
            }
        }
        if !found {
            ret := info
            if reverse {
                ret.State = fileinfo.Create
            } else {
                ret.State = fileinfo.Deleted
            }
            *result = append(*result, ret)
        }
    }
}

func syncProcessDiff(
    root string,
    result []fileinfo.FileInfo,
    trans transport.Transport,
    mstore store.MetaStore,
    statis *statistic.Statistic) error {
    //prepare
    for i := range result {
        relDir := io.SubPath(root, filepath.Dir(result[i].FilePath))
        result[i].From = uri.Get(uri.File, result[i].FilePath)
        uri, err := trans.GetUri(relDir, result[i].FileName)
        if err != nil {
            statis.AddFailedFile(result[i])
            return err
        }
        result[i].To = uri
        log.Debug("diff file : %v", result[i])
    }

    for i := range result {
        err := trans.Send(result[i])
        if err != nil {
            statis.AddFailedFile(result[i])
            return err
        }

        errSave := store.SaveMeta(mstore, result[i])
        if errSave != nil {
            return errSave
        }
    }
    return nil
}

func asyncProcessDiff(
    root string,
    result []fileinfo.FileInfo,
    trans transport.Transport,
    mstore store.MetaStore,
    statis *statistic.Statistic) error {
    //prepare
    size := len(result)
    var wg sync.WaitGroup
    wg.Add(size)

    exec := executor.NewFixedExecutor(4, size - 4)
    defer exec.Stop()

    for i := range result {
        index := i
        err := exec.Run(func() {
            defer wg.Done()

            relDir := io.SubPath(filepath.Dir(result[index].FilePath), root)
            result[index].From = uri.Get(uri.File, result[index].FilePath)
            uri, err := trans.GetUri(relDir, result[index].FileName)
            if err != nil {
                result[index].To = ""
                log.Error(err.Error())
                statis.AddFailedFile(result[index])
                return
            }
            result[index].To = uri
            log.Debug("diff file : %v", result[index])
        }, 0, nil)
        if err != nil {
            return err
        }
    }
    wg.Wait()

    wg = sync.WaitGroup{}
    wg.Add(size)
    for i := range result {
        index := i
        err := exec.Run(func() {
            defer wg.Done()
            if result[index].To == "" {
                return
            }
            err := trans.Send(result[index])
            if err != nil {
                log.Error(err.Error())
                statis.AddFailedFile(result[index])
                return
            }

            errSave := store.SaveMeta(mstore, result[index])
            if errSave != nil {
                log.Error(errSave.Error())
                return
            }
        }, 0, nil)
        if err != nil {
            return err
        }
    }
    wg.Wait()
    return nil
}
