// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package merge

import (
    "fbt/config"
    "fbt/errors"
    "fbt/store"
    "path/filepath"
)

func Merge(src, dest, save string) (err error) {
    srcPath := filepath.Join(src, config.InfoDir)
    destPath := filepath.Join(dest, config.InfoDir)

    srcStore := store.NewDefaultStore()
    err = srcStore.Open(srcPath, src)
    if err != nil {
        return errors.MergeInfoNotFound
    }
    defer srcStore.Close()

    err = srcStore.Read(src)
    if err != nil {
        return errors.MergeInfoNotFound
    }

    destStore := store.NewDefaultStore()
    err = destStore.Open(destPath, dest)
    if err != nil {
        return errors.MergeInfoNotFound
    }
    defer srcStore.Close()

    err = destStore.Read(dest)
    if err != nil {
        return errors.MergeInfoNotFound
    }

    return nil
}
