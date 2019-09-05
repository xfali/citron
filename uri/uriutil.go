// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package uri

type URI string

const (
    File = "file://"
)

func Get(uriType, path string) URI {
    return URI(uriType + path)
}

func (uri URI) String() string {
    return string(uri)
}
