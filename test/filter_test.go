// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package test

import (
    "citron/fileinfo"
    "citron/filter"
    "testing"
)

func TestFilter(t *testing.T) {
    fm := filter.FilterManager{}
    filter1 := func(info fileinfo.FileInfo, fc filter.FilterChain) error {
        t.Logf("filter1: %s\n", info.FilePath)
        fc.Filter(info)
        return nil
    }

    filter2 := func(info fileinfo.FileInfo, fc filter.FilterChain) error {
        t.Logf("filter2: %s\n", info.FilePath)
        fc.Filter(info)
        return nil
    }

    filter3 := func(info fileinfo.FileInfo, fc filter.FilterChain) error {
        t.Logf("filter3: %s\n", info.FilePath)
        fc.Filter(info)
        return nil
    }

    fm.Add(filter3, filter2, filter1)
    err1 := fm.RunFilter(fileinfo.FileInfo{FilePath:"./test"})
    t.Log(err1)
}

func TestFilter2(t *testing.T) {
    fm := filter.FilterManager{}
    filter1 := func(info fileinfo.FileInfo, fc filter.FilterChain) error {
        t.Logf("filter1: %s\n", info.FilePath)
        fc.Filter(info)
        return nil
    }

    filter2 := func(info fileinfo.FileInfo, fc filter.FilterChain) error {
        t.Logf("filter2: %s\n", info.FilePath)
        return nil
    }

    filter3 := func(info fileinfo.FileInfo, fc filter.FilterChain) error {
        t.Fatal("cannot be here!")
        fc.Filter(info)
        return nil
    }

    fm.Add(filter3, filter2, filter1)
    err1 := fm.RunFilter(fileinfo.FileInfo{FilePath:"./test"})
    t.Log(err1)

    err2 := fm.RunFilter(fileinfo.FileInfo{FilePath:"./test"})
    t.Log(err2)
}
