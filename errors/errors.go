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
var SourceAndTargetSame = &ErrCode{Code: "5003", Msg: "source path is the same as target uri"}
var SourceOrTargetError = &ErrCode{Code: "5004", Msg: "source dir or target uri error"}
var StoreNotFound = &ErrCode{Code: "10001", Msg: "store not found"}
var StoreOpenError = &ErrCode{Code: "10002", Msg: "store open error"}
var StoreFileNotFound = &ErrCode{Code: "10003", Msg: "store file not found"}
var TransportNotFound = &ErrCode{Code: "20001", Msg: "transport not found"}
var TransportOpenError = &ErrCode{Code: "20002", Msg: "transport open error"}
var TransportChecksumNotMatch = &ErrCode{Code: "20003", Msg: "transport checksum not match"}
var TransportBackupDirError = &ErrCode{Code: "20004", Msg: "backup dir exists!"}
var TransportReadSourceFileError = &ErrCode{Code: "20005", Msg: "transport read source file error"}
var MergeInfoNotFound = &ErrCode{Code: "30001", Msg: "merge dir not found info files"}

type ErrCode struct {
    Code string `json:"code"`
    Msg  string `json:"message"`
}

func (r *ErrCode) Error() string {
    return fmt.Sprintf("{\"code\": \"%s\", \"msg\": \"%s\"}", r.Code, r.Msg)
}
