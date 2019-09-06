// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package main

import (
    "fbt/config"
    "fbt/errors"
    "fbt/merge"
    "fbt/process"
    "fbt/store"
    "fbt/transport"
    "flag"
    "github.com/xfali/goutils/log"
    "path/filepath"
    "strings"
)

func main() {
    sourceDir := flag.String("s", "", "source dir")
    destUri := flag.String("d", "", "dest uri")
    checksumType := flag.String("check-type", "", "checksum type: MD5 | SHA256")
    incremental := flag.String("incr", "true", "incremental backup")
    newRepo := flag.String("n", "true", "creating a new backup repository every time")
    sync := flag.Bool("sync", false, "synchronous transport")
    mergeSrc := flag.String("merge-src", "", "path of src merge dir")
    mergeDest := flag.String("merge-dest", "", "path of dest merge dir")
    mergeSave := flag.String("merge-save", "", "dir save merge result")

    flag.Parse()

    config.GConfig.SourceDir = *sourceDir
    config.GConfig.DestUri = *destUri
    config.GConfig.ChecksumType = *checksumType
    config.GConfig.Incremental = *incremental == "true"
    config.GConfig.NewRepo = *newRepo == "true"
    config.GConfig.SyncTrans = *sync

    log.Info("config: %s\n", config.GConfig.String())

    if *mergeDest != "" || *mergeSrc != "" || *mergeSave != "" {
        if *mergeDest == "" || *mergeSrc == "" || *mergeSave == "" {
            log.Fatal("Merge param error, merge-src, merge-dest, merge-save must be not empty")
        }
        err := merge.Merge(*mergeSrc, *mergeDest, *mergeSave)
        if err != nil {
            log.Fatal(err.Error())
        }
    }

    checkResource()

    t, err := transport.Open(
        "file",
        config.GConfig.DestUri,
        config.GConfig.Incremental,
        config.GConfig.NewRepo)
    if err != nil {
        log.Fatal(err.Error())
    }
    defer t.Close()
    s, err := store.Open(
        "file",
        filepath.Join(filepath.Dir(config.GConfig.SourceDir), config.InfoDir),
        config.GConfig.SourceDir)
    if err != nil {
        log.Fatal(err.Error())
    }
    defer s.Close()

    errP := process.Process(config.GConfig.SourceDir, t, s)
    if errP != nil {
        log.Fatal(errP.Error())
    }
}

func checkResource() {
    if config.GConfig.SourceDir == "" {
        log.Fatal(errors.SourceDirNotExists.Error())
    }

    if config.GConfig.DestUri == "" {
        log.Fatal(errors.TargetUriEmpty.Error())
    }

    if config.GConfig.SourceDir == config.GConfig.DestUri {
        log.Fatal(errors.SourceAndTargetSame.Error())
    }

    if strings.Index(config.GConfig.DestUri, config.GConfig.SourceDir) != -1 {
        log.Fatal(errors.SourceOrTargetError.Error())
    }
}
