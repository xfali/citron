// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package transport

import (
    "citron/fileinfo"
    "citron/uri"
    "time"
)

type Listener interface {
    OnRead(int64)
    OnWrite(int64)
}

type Transport interface {
    Open(uri string, incremental, newRepo bool, timestamp time.Time, listener Listener) error
    GetUri(relDir, file string) (uri.URI, error)
    Send(info fileinfo.FileInfo) error
    Close() error
}
