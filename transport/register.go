// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package transport

import (
    "citron/errors"
    "github.com/xfali/goutils/log"
    "time"
)

var TransportCache = map[string]Transport{
    "file": NewDefaultTransport(),
}

func Open(transType, url string, incremental, newRepo bool, l Listener) (Transport, error) {
    if t, ok := TransportCache[transType]; ok {
        err := t.Open(url, incremental, newRepo, time.Now(), l)
        if err != nil {
            log.Warn(errors.TransportOpenError.Error())
            return nil, err
        }
        return t, nil
    }
    return nil, errors.StoreNotFound
}
