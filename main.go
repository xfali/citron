// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package main

import (
    "citron/config"
    "citron/ctx"
    "citron/errors"
    "citron/filter"
    "citron/merge"
    "citron/process"
    "citron/statistic"
    "citron/store"
    "citron/transport"
    "flag"
    "fmt"
    "github.com/xfali/goutils/log"
    "os"
    "path/filepath"
    "strings"
)

var logLv = map[string]int{
    "DEBUG": log.DEBUG,
    "INFO":  log.INFO,
    "WARN":  log.WARN,
    "ERROR": log.ERROR,
}

func main() {
    sourceDir := flag.String("s", "", "source dir")
    destUri := flag.String("d", "", "dest uri")
    checksumType := flag.String("checksum", "", "checksum type: MD5 | SHA1")
    incremental := flag.String("incr", "true", "incremental backup")
    newRepo := flag.String("n", "true", "creating a new backup repository every time")
    mergeSrc := flag.String("merge-src", "", "path of src merge dir")
    mergeDest := flag.String("merge-dest", "", "path of dest merge dir")
    mergeSave := flag.String("merge-save", "", "dir save merge result")
    logPath := flag.String("log-path", "./citron.log", "log file path")
    logLevel := flag.String("log-lv", "INFO", "log level: DEBUG | INFO | WARN | ERROR")
    multiTask := flag.Int("multi-task", 1, "backup multi task number")
    rmSrc := flag.Bool("remove-source", false, "remove source file")
    rmDelFile := flag.Bool("remove-del", false, "remove deleted source file")
    limit := flag.String("limit", "", "backup rate limit, for example: 20M/S or 256K/S")

    regexpBackup := flag.String("regexp-backup", "", "backup file select regexp")
    regexpHidden := flag.String("regexp-hidden", "", "hidden file regexp")
    regexpHelp := flag.Bool("regexp-help", false, "show regexp example")

    flag.Parse()

    if *regexpHelp {
        filter.PrintRegexp()
        os.Exit(0)
    }

    config.GConfig.SourceDir = *sourceDir
    config.GConfig.DestUri = *destUri
    config.GConfig.ChecksumType = *checksumType
    config.GConfig.Incremental = *incremental == "true"
    config.GConfig.NewRepo = *newRepo == "true"
    config.GConfig.MultiTaskNum = *multiTask
    config.GConfig.RmDel = *rmDelFile
    config.GConfig.RmSrc = *rmSrc
    config.GConfig.Limit = *limit

    config.GConfig.RegexpHidden = *regexpHidden
    config.GConfig.RegexpBackup = *regexpBackup

    fmt.Printf("config: %s\n", config.GConfig.String())
    logWriter := log.NewFileLogWriter(*logPath)
    if logWriter == nil {
        log.Fatal("log writer init failed")
    }
    defer logWriter.Close()
    log.Level = logLv[*logLevel]
    log.Writer = logWriter
    log.Log(log.Level, "config: %s\n", config.GConfig.String())

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

    st := statistic.New()
    rate, interval := config.GConfig.ParseLimit()
    limiter := statistic.NewLimiter(st, rate, interval)
    t, err := transport.Open(
        "file",
        config.GConfig.DestUri,
        config.GConfig.Incremental,
        config.GConfig.NewRepo,
        limiter)
    if err != nil {
        log.Fatal(err.Error())
    }
    defer t.Close()
    s, err := store.Open(
        "file",
        filepath.Join(filepath.Dir(config.GConfig.SourceDir), config.InfoDir, "root.meta"))
    if err != nil {
        log.Fatal(err.Error())
    }
    defer s.Close()

    ctx := ctx.Context{
        Transport: t,
        Store:     s,
        Statistic: st,
        Limiter:   limiter,
    }
    ctx.ConfigFilter(config.GConfig)
    errP := process.Process(config.GConfig.SourceDir, &ctx)

    fmt.Printf(st.String())
    log.Log(log.Level, st.String())

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
