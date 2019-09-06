// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package process

import (
    "fbt/config"
    "fbt/fileinfo"
    "fbt/store"
    "fbt/transport"
    "fbt/uri"
    "github.com/xfali/goutils/log"
    "sync"
)

func Process(srcDir string, trans transport.Transport, store store.MetaStore) error {
    if config.GConfig.Incremental {
        return incrementalProcess(srcDir, srcDir, trans, store)
    } else {
        return allProcess(srcDir, srcDir, trans, store)
    }
}

func process(list1, list2 []fileinfo.FileInfo, result *[]fileinfo.FileInfo, reverse bool) {
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

func processDiff(relDir string, result []fileinfo.FileInfo, trans transport.Transport, mstore store.MetaStore) (err error) {
    if len(result) > 0 {
        if config.GConfig.SyncTrans {
            err = syncProcessDiff(relDir, result, trans, mstore)
        } else {
            err = asyncProcessDiff(relDir, result, trans, mstore)
        }
        if err != nil {
            return err
        }
        mstore.Save()
    }

    return nil
}

func syncProcessDiff(relDir string, result []fileinfo.FileInfo, trans transport.Transport, mstore store.MetaStore) error {
    //prepare
    for i := range result {
        result[i].From = uri.Get(uri.File, result[i].FilePath)
        uri, err := trans.GetUri(relDir, result[i].FileName)
        if err != nil {
            return err
        }
        result[i].To = uri
        log.Debug("diff file : %v", result[i])
    }

    for i := range result {
        err := trans.Send(result[i])
        if err != nil {
            return err
        }

        errSave := store.SaveMeta(mstore, result[i])
        if errSave != nil {
            return errSave
        }
    }
    return nil
}

func asyncProcessDiff(relDir string, result []fileinfo.FileInfo, trans transport.Transport, mstore store.MetaStore) error {
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
