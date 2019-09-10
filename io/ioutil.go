// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package io

import (
    "fbt/checksum"
    "fbt/config"
    "fbt/fileinfo"
    "fmt"
    "io"
    "io/ioutil"
    "os"
    "path/filepath"
    "strings"
)

func SaveWrite(file *os.File, b []byte) error {
    if file != nil {
        _, err := file.WriteAt(b, 0)
        if err != nil {
            return err
        }
    }
    return nil
}

func GetDirFiles(path string) ([]fileinfo.FileInfo, error) {
    finfos, err := ioutil.ReadDir(path)
    if err != nil {
        return nil, err
    }

    ret := make([]fileinfo.FileInfo, len(finfos))
    i := 0
    for _, f := range finfos {
        filePath := filepath.Join(path, f.Name())
        ct, cs := "", ""
        if !f.IsDir() {
            ct = config.GConfig.ChecksumType
            cs, err = checksum.GetFileCheckSum(checksum.Get(), filePath)
            if err != nil {
                return nil, err
            }
        }

        info := fileinfo.FileInfo{
            FileName: f.Name(),
            FilePath: filePath,
            Parent:   path,
            IsDir:    f.IsDir(),
            ModTime:  f.ModTime(),
            Size:     f.Size(),
            //default create
            State:        fileinfo.Create,
            ChecksumType: ct,
            Checksum:     cs,
        }
        ret[i] = info
        i++
    }
    return ret, nil
}

func SubPath(src, root string) string {
    src = filepath.Clean(src)
    root = filepath.Clean(root)

    rel := strings.Replace(src, root, "", 1)
    if len(rel) > 0 && rel[0:1] == string(filepath.Separator) {
        rel = rel[1:]
    }
    return rel
}

func CopyFile(src, dest string) error {
    sourceFileStat, err := os.Stat(src)
    if err != nil {
        return err
    }

    if !sourceFileStat.Mode().IsRegular() {
        return fmt.Errorf("%s is not a regular file", src)
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
    _, errCpy := io.Copy(destination, source)
    return errCpy
}
