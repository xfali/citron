// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package test

import (
    "citron/checksum"
    "citron/fileinfo"
    "testing"
)

var file = "./checksum_test.go"

func TestDummy(t *testing.T) {
    s, e := checksum.GetFileCheckSum(checksum.Get(), file)
    if e != nil {
        t.Fatal(e)
    }
    t.Log(string(s))
}

func TestMd5(t *testing.T) {
    hash := checksum.New(fileinfo.MD5)
    s, e := checksum.GetFileCheckSum(hash, file)
    if e != nil {
        t.Fatal(e)
    }
    t.Log(string(s))
}

func TestSha1(t *testing.T) {
    hash := checksum.New(fileinfo.SHA1)
    s, e := checksum.GetFileCheckSum(hash, file)
    if e != nil {
        t.Fatal(e)
    }
    t.Log(string(s))
}

func BenchmarkDummy(t *testing.B) {
    hash := checksum.Get()
    for i := 0; i < t.N; i++ {
        _, e := checksum.GetFileCheckSum(hash, file)
        if e != nil {
            t.Fatal(e)
        }
    }
}

func BenchmarkMd5(t *testing.B) {
    hash := checksum.New(fileinfo.MD5)
    for i := 0; i < t.N; i++ {
        _, e := checksum.GetFileCheckSum(hash, file)
        if e != nil {
            t.Fatal(e)
        }
    }
}

func BenchmarkSha1(t *testing.B) {
    hash := checksum.New(fileinfo.SHA1)
    for i := 0; i < t.N; i++ {
        _, e := checksum.GetFileCheckSum(hash, file)
        if e != nil {
            t.Fatal(e)
        }
    }
}
