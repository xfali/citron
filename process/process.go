// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package process

import (
    "fbt/config"
    "fbt/fileinfo"
    "fbt/statistic"
    "fbt/store"
    "fbt/transport"
    "fbt/uri"
    "github.com/xfali/goutils/log"
    "sync"
)

type ProcFunc func(string, []fileinfo.FileInfo) error

func Process(srcDir string, trans transport.Transport, store store.MetaStore, statis *statistic.Statistic) error {
    statisticFunc := func(relDir string, result []fileinfo.FileInfo) error {
        if len(result) > 0 {
            for _, info := range result {
                if !info.IsDir {
                    statis.AddFileCount(1)
                    statis.AddTotalSize(info.Size)
                }
            }
        }

        return nil
    }

    trancFunc := func(relDir string, result []fileinfo.FileInfo) error {
        var err error
        if len(result) > 0 {
            if config.GConfig.SyncTrans {
                err = syncProcessDiff(relDir, result, trans, store, statis)
            } else {
                err = asyncProcessDiff(relDir, result, trans, store, statis)
            }
            if err != nil {
                return err
            }
            store.Save()
        }

        return nil
    }

    if config.GConfig.Incremental {
        err := incrementalProcess(srcDir, srcDir, store, statisticFunc)
        if err != nil {
            return err
        }
        statis.ResetTime()
        return incrementalProcess(srcDir, srcDir, store, trancFunc)
    } else {
        err := allProcess(srcDir, srcDir, store, statisticFunc)
        if err != nil {
            return err
        }
        statis.ResetTime()
        return allProcess(srcDir, srcDir, store, trancFunc)
    }
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
    relDir string,
    result []fileinfo.FileInfo,
    trans transport.Transport,
    mstore store.MetaStore,
    statis *statistic.Statistic) error {
    //prepare
    for i := range result {
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
    relDir string,
    result []fileinfo.FileInfo,
    trans transport.Transport,
    mstore store.MetaStore,
    statis *statistic.Statistic) error {
    //prepare
    size := len(result)
    var wg sync.WaitGroup
    wg.Add(size)
    for i := range result {
        index := i
        go func() {
            defer wg.Done()

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
        }()
    }
    wg.Wait()

    wg = sync.WaitGroup{}
    wg.Add(size)
    for i := range result {
        index := i
        go func() {
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
        }()
    }
    wg.Wait()
    return nil
}
