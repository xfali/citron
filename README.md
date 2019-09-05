# fbt

## 介绍

fbt是一个简单实现的文件备份服务

## 使用

```$xslt
./fbt -s c:/tmp/src -d c:/tmp/dest
```

-s 源目录

-d 目标备份目录

默认使用增量备份，但首次备份是会自动全量备份

## 注意事项

不要删除与源目录同级目录下的.fbtmeta目录(默认隐藏)。删除会造成增量备份不可用。

## TODO

1. 目前只支持本地目录 -> 本地目录的备份，后续支持其他类型的备份方式，如FTP等；
2. 文件校验功能较弱，目前仅以修改时间作为校验依据，不够严谨；
3. 增量备份的merge功能未实现（由多个增量备份目录自动、高效、安全的合并为一个文件目录）；
4. 备份进度监视器待开发。