// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package process

import (
    "citron/cmd"
    "citron/config"
    "citron/ctx"
    "citron/fileinfo"
    "github.com/xfali/executor"
    "path/filepath"
    "sync"
)

type ProcFunc func([]fileinfo.FileInfo) error

func Process(srcDir string, ctx *ctx.Context) error {
    srcDir = filepath.Clean(srcDir)
    var diffFile []fileinfo.FileInfo
    statisticFunc := func(result []fileinfo.FileInfo) error {
        if len(result) > 0 {
            for _, info := range result {
                diffFile = append(diffFile, info)
                if !info.IsDir {
                    ctx.Statistic.AddFileCount(1)
                    ctx.Statistic.AddTotalSize(info.Size)
                }
            }
        }

        return nil
    }

    if config.GConfig.Incremental {
        err := incrementalProcess(srcDir, ctx.Store, statisticFunc)
        if err != nil {
            return err
        }
    } else {
        err := allProcess(srcDir, ctx.Store, statisticFunc)
        if err != nil {
            return err
        }
    }

    p := cmd.NewProgress(ctx.Statistic)
    ctx.Statistic.ResetTime()
    ctx.Limiter.ResetTime()
    p.Start()
    defer p.Stop()

    var err error
    if len(diffFile) > 0 {
        if config.GConfig.MultiTaskNum <= 1 {
            err = syncProcessDiff(srcDir, diffFile, ctx)
        } else {
            err = asyncProcessDiff(srcDir, diffFile, ctx)
        }
        if err != nil {
            return err
        }
        ctx.Store.Save()
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
    ctx *ctx.Context) error {
    //prepare
    for i := range result {
        err := ctx.GetUri(&result[i], root)
        if err != nil {
            return err
        }
    }

    for i := range result {
        err := ctx.FilterMgr.RunFilter(result[i])
        if err != nil {
            return err
        }
    }
    return nil
}

func asyncProcessDiff(
    root string,
    result []fileinfo.FileInfo,
    ctx *ctx.Context) error {
    //prepare
    size := len(result)
    var wg sync.WaitGroup
    wg.Add(size)

    exec := executor.NewFixedExecutor(config.GConfig.MultiTaskNum, size - config.GConfig.MultiTaskNum)
    defer exec.Stop()

    for i := range result {
        index := i
        err := exec.Run(func() {
            defer wg.Done()
            ctx.GetUri(&result[index], root)
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
            ctx.FilterMgr.RunFilter(result[index])
        }, 0, nil)
        if err != nil {
            return err
        }
    }
    wg.Wait()
    return nil
}
