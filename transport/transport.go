// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package transport

import (
    "fbt/fileinfo"
    "fbt/uri"
    "time"
)

type Listener interface {
    AddReadSize(int64) int64
    AddWriteSize(int64) int64
}

type Transport interface {
    Open(uri string, incremental, newRepo bool, timestamp time.Time, listener Listener) error
    GetUri(relDir, file string) (uri.URI, error)
    Send(info fileinfo.FileInfo) error
    Close() error
}
