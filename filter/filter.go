// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package filter

import "citron/fileinfo"

//中断处理返回true，继续处理返回false
type Filter func(info fileinfo.FileInfo, fc FilterChain) error

type FilterChain []Filter
type FilterManager FilterChain

func (fc *FilterManager) Add(filter ...Filter) {
    *fc = append(*fc, filter...)
}

func (fc FilterManager) RunFilter(info fileinfo.FileInfo) error {
    return FilterChain(fc).Filter(info)
}

func (fc FilterChain) Filter(info fileinfo.FileInfo) error {
    if len(fc) > 0 {
        filter := fc[len(fc)-1]
        chain := fc.next()
        return filter(info, chain)
    }
    return nil
}

func (fc FilterChain) next() FilterChain {
    if len(fc) > 0 {
        return fc[:len(fc)-1]
    } else {
        return FilterChain{}
    }
}
