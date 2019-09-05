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

type Transport interface {
    Open(uri string, incremental bool, timestamp time.Time) error
    GetUri(relDir, file string) (uri.URI, error)
    Send(info fileinfo.FileInfo) error
    Close() error
}
