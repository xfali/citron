// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package checksum

import (
    "crypto/md5"
    "crypto/sha1"
    "encoding/hex"
    "fbt/config"
    "fbt/fileinfo"
    "github.com/xfali/goutils/log"
    "hash"
    "io"
    "os"
)

var gDummyHash DummyHash = ""

func New(typename string) hash.Hash {
    switch typename {
    case fileinfo.MD5:
        return md5.New()
    case fileinfo.SHA1:
        return sha1.New()
    }
    return &gDummyHash
}

func Get() hash.Hash {
    return New(config.GConfig.ChecksumType)
}

func Format(sum []byte) string {
    return hex.EncodeToString(sum)
}

func GetFileCheckSum(hash hash.Hash, path string) (string, error) {
    if _, ok := hash.(*DummyHash); ok {
        return "", nil
    }

    f, err := os.Open(path)
    if err != nil {
        log.Error("Open %v", err)
        return "", err
    }
    defer f.Close()

    if _, err := io.Copy(hash, f); err != nil {
        log.Error("Copy %v", err)
        return "", err
    }

    return Format(hash.Sum(nil)), nil
}

var dummy_hash = []byte("DUMMY_HASH")

type DummyHash string

func (h *DummyHash) Write(b []byte) (int, error) {
    return len(b), nil
}

func (h *DummyHash) Sum(b []byte) []byte {
    return []byte("")
}

func (h *DummyHash) Reset() {

}

func (h *DummyHash) Size() int {
    return len(dummy_hash)
}

func (h *DummyHash) BlockSize() int {
    return 0
}
