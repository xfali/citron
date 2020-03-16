// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package store

import (
    "citron/errors"
    "github.com/xfali/goutils/log"
)

var StoreCache = map[string]MetaStore{
    "file": NewDefaultStore(),
}

func Register(storeType string, store MetaStore) {
    StoreCache[storeType] = store
}

func Open(storeType, storeUri string) (MetaStore, error) {
    if s, ok := StoreCache[storeType]; ok {
        err := s.Open(storeUri)
        if err != nil {
            log.Warn(errors.StoreOpenError.Error())
            return nil, err
        }
        return s, nil
    }
    return nil, errors.StoreNotFound
}
