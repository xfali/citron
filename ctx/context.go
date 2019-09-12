// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package ctx

import (
    "citron/config"
    "citron/fileinfo"
    "citron/filter"
    "citron/io"
    "citron/statistic"
    "citron/store"
    "citron/transport"
    "citron/uri"
    "github.com/xfali/goutils/log"
    "path/filepath"
)

type Context struct {
    Transport transport.Transport
    Store     store.MetaStore
    Statistic *statistic.Statistic
    FilterMgr *filter.FilterManager
}

func (ctx *Context) ConfigFilter(cfg config.Config) {
    ctx.FilterMgr = &filter.FilterManager{}
    ctx.FilterMgr.Add(ctx.SendFile)

    if cfg.RmSrc {
        ctx.FilterMgr.Add(filter.RmSourceFilter)
    }

    if cfg.RegexpHidden != "" {
        ctx.FilterMgr.Add(filter.NewRegexp(cfg.RegexpHidden).HideFiler)
    }

    if !cfg.RmDel {
        ctx.FilterMgr.Add(filter.KeepDelFiler)
    }

    if cfg.RegexpBackup != "" {
        ctx.FilterMgr.Add(filter.NewRegexp(cfg.RegexpBackup).BackupFiler)
    }
}

func (ctx *Context) GetUri(info *fileinfo.FileInfo, root string) error {
    relDir := io.SubPath(filepath.Dir(info.FilePath), root)
    info.From = uri.Get(uri.File, info.FilePath)
    uri, err := ctx.Transport.GetUri(relDir, info.FileName)
    if err != nil {
        info.To = ""
        log.Error(err.Error())
        ctx.Statistic.AddFailedFile(*info)
        return err
    }
    info.To = uri
    log.Debug("diff file : %v", *info)
    return nil
}

func (ctx *Context) SendFile(info fileinfo.FileInfo, fc filter.FilterChain) error {
    if info.To == "" {
        return nil
    }
    err := ctx.Transport.Send(info)
    if err != nil {
        log.Error(err.Error())
        ctx.Statistic.AddFailedFile(info)
        return err
    }

    errSave := store.SaveMeta(ctx.Store, info)
    if errSave != nil {
        log.Error(errSave.Error())
        return errSave
    }
    return nil
}
