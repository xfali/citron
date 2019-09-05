// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package transport

import (
    "errors"
    "fbt/config"
    "fbt/fileinfo"
    "fbt/history"
    myio "fbt/io"
    "fbt/store"
    "fbt/uri"
    "github.com/xfali/goutils/io"
    "github.com/xfali/goutils/log"
    "os"
    "path/filepath"
    "time"
)

const (
    version     = "v0.0.1"
    historyFile = "history"
)

type LocalFileTransport struct {
    target      string
    timestamp   time.Time
    store       store.MetaStore
    incremental bool
    version     string
    record      *history.Recorder
    backupDir   string
}

func NewDefaultTransport() Transport {
    l := LocalFileTransport{
        version: version,
        record:  history.New(),
    }
    return &l
}

func (t *LocalFileTransport) Open(uri string, incremental bool, timestamp time.Time) error {
    t.target = uri
    if !io.IsPathExists(t.target) {
        err := io.Mkdir(t.target)
        if err != nil {
            return err
        }
    }
    err := t.record.Open(filepath.Join(t.target, historyFile))
    if err != nil {
        return err
    }

    if !incremental {

    }

    t.incremental = incremental
    t.timestamp = timestamp

    errP := t.prepareBackupDir()
    if errP != nil {
        return errP
    }

    s := store.NewDefaultStore()
    errO := s.Open( filepath.Join(t.backupDir, config.InfoDir), t.backupDir)
    if errO != nil {
        return errO
    }
    errR := s.Read(t.backupDir)
    if errR != nil {
        return errR
    }
    t.store = s

    return t.record.Append(history.History{
        Timestamp:   timestamp,
        Path:        t.backupDir,
        Version:     t.version,
        Incremental: t.incremental,
    })
}

func (t *LocalFileTransport) prepareBackupDir() error {
    dir := t.timestamp.Format("20060102150405")
    dir = filepath.Join(t.target, dir)
    if io.IsPathExists(dir) {
        return errors.New("backup dir exists! ")
    } else {
        io.Mkdir(dir)
        t.backupDir = dir
    }
    return nil
}

func (t *LocalFileTransport) Send(info fileinfo.FileInfo) error {
    switch info.State {
    case fileinfo.Create, fileinfo.Modified:
        err := create(info)
        if err != nil {
            return err
        }
        break
    case fileinfo.Deleted:
        err := remove(info)
        if err != nil {
            return err
        }
        break
    }
    return t.store.Insert(info)
}

func remove(info fileinfo.FileInfo) error {
    path := GetPath(info.To)
    if io.IsPathExists(path) {
        if info.IsDir {
            return os.RemoveAll(path)
        } else {
            return os.Remove(path)
        }
    } else {
        log.Info("file not found %s", path)
    }
    return nil
}

func create(info fileinfo.FileInfo) error {
    src := GetPath(info.From)
    dest := GetPath(info.To)
    if info.IsDir {
        return io.Mkdir(dest)
    } else {
        return myio.CopyFile(src, dest)
    }
}

func (t *LocalFileTransport) GetUri(relDir, file string) (uri.URI, error) {
    return uri.URI(uri.File + filepath.Join(t.backupDir, relDir, file)), nil
}

func GetPath(path uri.URI) string {
    if len(path) >= len(uri.File) {
        if path[:len(uri.File)] == uri.File {
            return string(path[len(uri.File):])
        } else {
            return ""
        }
    }
    return ""
}

func (t *LocalFileTransport) Close() error {
    defer t.record.Close()
    if t.store != nil {
        return t.store.Close()
    }
    return nil
}
