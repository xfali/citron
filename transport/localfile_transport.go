// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package transport

import (
    "citron/checksum"
    "citron/config"
    "citron/errors"
    "citron/fileinfo"
    "citron/history"
    "citron/store"
    "citron/uri"
    uio "github.com/xfali/goutils/io"
    "github.com/xfali/goutils/log"
    "io"
    "os"
    "path/filepath"
    "time"
)

const (
    version     = "v0.0.1"
    historyFile = "history"
)

type LocalFileTransport struct {
    //仓库入口地址，备份目录的根目录
    target string
    //时间戳
    timestamp time.Time
    //单次备份详细日志
    store store.MetaStore
    //是否每次都使用新仓储
    newRepo bool
    //是否增量备份
    incremental bool
    //版本
    version string
    //备份日志
    record *history.Recorder
    //实际备份目标目录
    backupDir string
    //监听器
    listener Listener
}

func NewDefaultTransport() Transport {
    l := LocalFileTransport{
        version:  version,
        record:   history.New(),
        listener: FakeListener(0),
    }
    return &l
}

func (t *LocalFileTransport) Open(uri string, incremental, newRepo bool, timestamp time.Time, listener Listener) error {
    t.target = uri
    if !uio.IsPathExists(t.target) {
        err := uio.Mkdir(t.target)
        if err != nil {
            return err
        }
    }
    err := t.record.Open(filepath.Join(t.target, historyFile))
    if err != nil {
        return err
    }

    t.incremental = incremental
    t.newRepo = newRepo
    t.timestamp = timestamp

    errP := t.prepareBackupDir()
    if errP != nil {
        return errP
    }

    s := store.NewDefaultStore()
    errO := s.Open(filepath.Join(t.backupDir, config.InfoDir, t.storeFile()))
    if errO != nil {
        t.record.Close()
        return errO
    }
    t.store = s

    t.listener = listener

    return t.record.Append(history.History{
        Timestamp:   timestamp,
        Path:        t.backupDir,
        Version:     t.version,
        Incremental: t.incremental,
    })
}

func (t *LocalFileTransport) prepareBackupDir() error {
    dir := "backup"
    if t.newRepo {
        dir = t.timeStr()
    }

    dir = filepath.Join(t.target, dir)
    if uio.IsPathExists(dir) {
        if t.newRepo {
            return errors.TransportBackupDirError
        }
    } else {
        uio.Mkdir(dir)
    }
    t.backupDir = dir

    return nil
}

func (t *LocalFileTransport) storeFile() string {
    name := ""
    if t.newRepo {
        name = "root"
    } else {
        name = t.timeStr()
    }
    return name + ".meta"
}

func (t *LocalFileTransport) timeStr() string {
    return t.timestamp.Format("20060102150405")
}

func (t *LocalFileTransport) Send(info fileinfo.FileInfo) error {
    log.Info("Send from %s to %s", info.From, info.To)
    switch info.State {
    case fileinfo.Create, fileinfo.Modified:
        err := t.create(&info)
        if err != nil {
            return err
        }
        break
    case fileinfo.Deleted:
        err := t.remove(&info)
        if err != nil {
            return err
        }
        break
    }
    return t.store.Insert(info)
}

func (t *LocalFileTransport) remove(info *fileinfo.FileInfo) error {
    path := GetPath(info.To)
    if uio.IsPathExists(path) {
        if info.IsDir {
            err := os.RemoveAll(path)
            if err != nil {
                return err
            }
        } else {
            err := os.Remove(path)
            if err != nil {
                return err
            }
        }
    } else {
        log.Info("file not found %s", path)
    }
    info.FilePath = ""
    return nil
}

func (t *LocalFileTransport) create(info *fileinfo.FileInfo) error {
    src := GetPath(info.From)
    dest := GetPath(info.To)
    if info.IsDir {
        //FIXME: unnecessary
        err := uio.Mkdir(dest)
        if err != nil {
            return err
        }
        info.FilePath = dest
        return nil
    } else {
        dir := filepath.Dir(dest)
        if !uio.IsPathExists(dir) {
            if err := uio.Mkdir(dir); err != nil {
                return err
            }
        }
        err := t.copyFileByPath(src, dest)
        if err != nil {
            return err
        }

        errC := checkFile(*info, dest)
        if errC != nil {
            return errC
        }

        info.FilePath = dest
        return nil
    }
}

func checkFile(info fileinfo.FileInfo, targetFile string) error {
    if info.Checksum != "" && info.ChecksumType != "" {
        hash := checksum.New(info.ChecksumType)
        sum, err := checksum.GetFileCheckSum(hash, targetFile)
        if err != nil {
            return err
        }
        if sum != info.Checksum {
            return errors.TransportChecksumNotMatch
        }
    }
    return nil
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

func (t *LocalFileTransport) copyFileByPath(src, dest string) error {
    sourceFileStat, err := os.Stat(src)
    if err != nil {
        return err
    }

    if !sourceFileStat.Mode().IsRegular() {
        log.Error("%s is not a regular file", src)
        return errors.TransportReadSourceFileError
    }

    source, err := os.Open(src)
    if err != nil {
        return err
    }
    defer source.Close()

    destination, err := os.Create(dest)
    if err != nil {
        return err
    }
    defer destination.Close()

    _, errCpy := t.copyFile(destination, source)
    return errCpy
}

func (t *LocalFileTransport) copyFile(dst io.Writer, src io.Reader) (written int64, err error) {
    // If the reader has a WriteTo method, use it to do the copy.
    // Avoids an allocation and a copy.
    if wt, ok := src.(io.WriterTo); ok {
        return wt.WriteTo(dst)
    }
    // Similarly, if the writer has a ReadFrom method, use it to do the copy.
    if rt, ok := dst.(io.ReaderFrom); ok {
        return rt.ReadFrom(src)
    }
    var buf []byte
    size := 32 * 1024
    if l, ok := src.(*io.LimitedReader); ok && int64(size) > l.N {
        if l.N < 1 {
            size = 1
        } else {
            size = int(l.N)
        }
    }
    buf = make([]byte, size)

    for {
        nr, er := src.Read(buf)
        if nr > 0 {
            t.listener.AddReadSize(int64(nr))
            nw, ew := dst.Write(buf[0:nr])
            if nw > 0 {
                t.listener.AddWriteSize(int64(nw))
                written += int64(nw)
            }
            if ew != nil {
                err = ew
                break
            }
            if nr != nw {
                err = io.ErrShortWrite
                break
            }
        }
        if er != nil {
            if er != io.EOF {
                err = er
            }
            break
        }
    }
    return written, err
}

type FakeListener int64

func (s FakeListener) AddReadSize(delta int64) int64 {
    return 0
}

func (s FakeListener) AddWriteSize(delta int64) int64 {
    return 0
}
