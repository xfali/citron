// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package errors

import (
    "fmt"
)

var SourceDirNotExists = &ErrCode{Code: "5001", Msg: "source dir not exits"}
var TargetUriEmpty = &ErrCode{Code: "5002", Msg: "target uri is empty"}
var StoreNotFound = &ErrCode{Code: "10001", Msg: "store not found"}
var StoreOpenError = &ErrCode{Code: "10002", Msg: "store open error"}
var StoreFileNotFound = &ErrCode{Code: "10003", Msg: "store file not found"}
var TransportNotFound = &ErrCode{Code: "20001", Msg: "transport not found"}
var TransportOpenError = &ErrCode{Code: "20002", Msg: "transport open error"}
var MergeInfoNotFound = &ErrCode{Code: "30001", Msg: "merge dir not found info files"}

type ErrCode struct {
    Code string `json:"code"`
    Msg  string `json:"message"`
}

func (r *ErrCode) Error() string {
    return fmt.Sprintf("{\"code\": \"%s\", \"msg\": \"%s\"}", r.Code, r.Msg)
}