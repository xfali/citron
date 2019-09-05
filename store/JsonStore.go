// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package store

import (
    "encoding/json"
    "fbt/fileinfo"
    "fbt/io"
    utilio "github.com/xfali/goutils/io"
    "io/ioutil"
    "os"
    "path/filepath"
    "strings"
)

const(
    MetaFileName = ".meta"
)

type JsonStore struct {
    data    []fileinfo.FileInfo
    file    *os.File
    root    string
    metaDir string
}

func NewDefaultStore() MetaStore {
    i := JsonStore{}
    return &i
}

func (s *JsonStore) Open(storeDir string, dir string) error {
    metaPath := storeDir
    if !utilio.IsPathExists(metaPath) {
        err := utilio.Mkdir(metaPath)
        utilio.SetInvisible(metaPath)
        if err != nil {
            return err
        }
    }
    s.metaDir= filepath.Clean(metaPath)
    s.root = filepath.Clean(dir)

    return nil
}

func (s *JsonStore) resetData(dir string) error {
    //close at first
    s.Close()

    dir = filepath.Clean(dir)

    rel := io.SubPath(dir, s.root)
    rel = strings.Replace(rel, string(filepath.Separator), "_", -1)
    if rel == "" {
        rel = "root"
    }
    rel += MetaFileName
    path := filepath.Join(s.metaDir, rel)
    if utilio.IsPathExists(path) {
        file, err := os.OpenFile(path, os.O_RDONLY, 0644)
        if err != nil {
            return err
        }
        defer file.Close()
        b, err := ioutil.ReadAll(file)
        if err == nil && len(b) > 0 {
            var ret []fileinfo.FileInfo
            err := json.Unmarshal(b, &ret)
            if err != nil {
                return err
            }
            s.data = ret
        }
    }
    file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
    if err != nil {
        return err
    }
    s.file = file

    return nil
}

func (s *JsonStore) Insert(info fileinfo.FileInfo) error {
    i := s.find(info.FilePath)
    if i != -1 {
        return s.update(i, info)
    }

    return s.insert(info)
}

func (s *JsonStore) Update(info fileinfo.FileInfo) error {
    i := s.find(info.FilePath)
    if i == -1 {
        return s.insert(info)
    }

    return s.update(i, info)
}

func (s *JsonStore) update(index int, info fileinfo.FileInfo) error {
    s.data[index] = info
    return nil
}

func (s *JsonStore) insert(info fileinfo.FileInfo) error {
    s.data = append(s.data, info)
    return nil
}

func (s *JsonStore) find(filepath string) int {
    for i := range s.data {
        if s.data[i].FilePath == filepath {
            return i
        }
    }
    return -1
}

func (s *JsonStore) Read(dir string) error {
    err := s.resetData(dir)
    if err != nil {
        return err
    }

    return nil
}

func (s *JsonStore) Query() ([]fileinfo.FileInfo, error) {
    return s.data, nil
}

func (s *JsonStore) QueryByPath(uri string) (fileinfo.FileInfo, error) {
    i := s.find(uri)
    if i == -1 {
        return fileinfo.FileInfo{}, nil
    }
    return s.data[i], nil
}

func (s *JsonStore) Delete(info fileinfo.FileInfo) error {
    i := s.find(info.FilePath)
    if i == -1 {
        return nil
    }

    s.data = append(s.data[:i], s.data[i+1:]...)
    return nil
}

func (s *JsonStore) Save() error {
    if s.file != nil && len(s.data) > 0 {
        b, err := json.Marshal(s.data)
        if err != nil {
            return err
        }

        err2 := io.SaveWrite(s.file, b)
        if err2 != nil {
            return err
        }
    }
    return nil
}

func (s *JsonStore) Close() error {
    s.Save()
    s.data = []fileinfo.FileInfo{}
    if s.file != nil {
        err := s.file.Close()
        s.file = nil
        return err
    }
    return nil
}
