// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package store

import (
    "fbt/errors"
    "github.com/xfali/goutils/log"
)

var StoreCache = map[string]MetaStore{
    "file": NewDefaultStore(),
}

func Open(storeType, storeUri, dataDir string) (MetaStore, error) {
    if s, ok := StoreCache[storeType]; ok {
        err := s.Open(storeUri, dataDir)
        if err != nil {
            log.Warn(errors.StoreOpenError.Error())
            return nil, err
        }
        return s, nil
    }
    return nil, errors.StoreNotFound
}
