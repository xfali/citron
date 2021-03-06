# citron

## 介绍

![citron Logo](images/citron.png)

citron是一个简单实现的文件备份服务

## 安装

```
go get github.com/xfali/citron
```

## 使用

```$xslt
./citron -s $SRC_DIR -d $DEST_DIR
```

-s 源目录

-d 目标备份目录

默认使用增量备份，首次备份会自动全量备份

```$xslt
  -checksum string
        文件校验方式: MD5 | SHA1
  -s string
        源备份路径
  -d string
        备份目标路径
  -incr string
        是否增量备份 (default "true")
  -n string
        是否每次都创建一个新的备份仓库，不覆盖旧仓库 (default "true")
  -sync
        是否同步备份（默认异步）
  -log-lv string
        日志级别: DEBUG | INFO | WARN | ERROR (default "INFO")
  -log-path string
        日志文件 (default "./citron.log")
  -multi-task int
        备份任务数，1为同步备份 (default 1)
  -remove-del
        当源文件被删除时，同步删除备份仓库的文件（默认为false不删除）
  -remove-source
        文件备份完成后删除源文件（默认为false不删除）
  -regexp-backup string
        设置备份过滤正则表达式
  -regexp-hidden string
        设置隐藏过滤正则表达式
  -regexp-help
        打印常用正则表达式

```

### 建议

**默认的配置为：**

开启增量备份，且不覆盖旧备份仓库，这种方式是最安全的备份方式。

优势：
1. 性能好
2. 数据全
3. 保留修改的历史数据

劣势：
1. 目前还未完成仓库合并功能，会出现多个仓库，用户查找文件比较麻烦
2. 占用较大存储空间：文件删除、修改都会保留历史文件

**其他配置参数：**

```$xslt
./citron -s $SRC_DIR -d $DEST_DIR -n=false
```

增加-n参数，表示不使用多仓库，只使用同一个仓库（目录）进行增量备份

优势：
1. 性能好
2. 存储空间小
3. 由于只有一个备份目录，用户查找文件方便

劣势：
1. 文件每次修改都会被覆盖，历史副本会丢失
2. 如果配置了-remove-del，当源目录的文件被删除时，则备份目录会同步删除，且无法找回

### 高级 - 正则过滤

* 备份过滤

参数：-regexp-backup

```
 ./citron -s $SRC_DIR -d $DEST_DIR -n=false -regexp-backup="^\S+\.go$"
```
备份.go后缀的所有文件


* 隐藏过滤

参数：-regexp-hidden

```
 ./citron -s $SRC_DIR -d $DEST_DIR -n=false -regexp-hidden="^\S+\.go$"
```
隐藏备份仓库中.go后缀的所有文件


## 注意事项

不要删除与源目录同级目录下的.citronmeta目录(默认隐藏)。删除会造成增量备份异常或不可用。

## TODO

1. 目前只支持本地目录 -> 本地目录的备份，后续支持其他类型的备份方式，如NAS等；
2. 增量备份多仓库的merge功能未实现（由多个增量备份目录自动、高效、安全的合并为一个文件目录）；
3. 文件监视功能待开发；
4. 避免重复文件备份功能？属于用户行为，暂不支持；
5. ~~备份文件过滤功能；~~ (已实现)
6. ~~显示备份进度条；~~ (已实现)