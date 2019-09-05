// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package history

import (
    "encoding/json"
    "fbt/io"
    utilio "github.com/xfali/goutils/io"
    "io/ioutil"
    "os"
    "time"
)

type History struct {
    Timestamp   time.Time `json:"timestamp"`
    Path        string    `json:"path"`
    Version     string    `json:"version"`
    Incremental bool      `json:"incremental"`
}

type Recorder struct {
    file *os.File
    data []History
}

func New() *Recorder {
    return &Recorder{}
}

func (r *Recorder) Open(uri string) error {
    if utilio.IsPathExists(uri) {
        file, err := os.OpenFile(uri, os.O_RDONLY, 0644)
        if err != nil {
            return err
        }
        defer file.Close()
        b, err := ioutil.ReadAll(file)
        if err == nil && len(b) > 0 {
            var ret []History
            err := json.Unmarshal(b, &ret)
            if err != nil {
                return err
            }
            r.data = ret
        }
    }
    file, err := os.OpenFile(uri, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
    if err != nil {
        return err
    }
    r.file = file

    return nil
}

func (r *Recorder) Append(record History) error {
    r.data = append(r.data, record)
    return nil
}

func (r *Recorder) Save() error {
    if r.file != nil && len(r.data) > 0 {
        b, err := json.Marshal(r.data)
        if err != nil {
            return err
        }

        err2 := io.SaveWrite(r.file, b)
        if err2 != nil {
            return err
        }
    }
    return nil
}

func (r *Recorder) Close() error {
    r.Save()
    r.data = []History{}
    if r.file != nil {
        return r.file.Close()
    }
    return nil
}
