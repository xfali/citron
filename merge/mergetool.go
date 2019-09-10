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
    err = srcStore.Open(srcPath)
    if err != nil {
        return errors.MergeInfoNotFound
    }
    defer srcStore.Close()

    destStore := store.NewDefaultStore()
    err = destStore.Open(destPath)
    if err != nil {
        return errors.MergeInfoNotFound
    }
    defer srcStore.Close()

    return nil
}
