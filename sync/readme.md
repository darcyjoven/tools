# sync 同步工具

> windows 端运行上传和下载linux服务器的工具

## CLI

+ sync up --file D:\app\apbak\u1\topprod\topcust\cim\4gl\cimi999.4gl
+ sync down -f D:\app\apbak\u1\topprod\topcust\cim\4gl\cimi999.4gl

## 配置文件

```yaml
author: darcy
diffdownload: info
diffupload: info
format: "060102"
localdir: D:/app/apbak/u1/
gitDir: D:/app/apbak/u1
logdir: ./temp
loginterval: one
logname: main
remotedir: /u1/
remotestr: root@192.168.1.1
tempdir: D:/app/apbak/u1/tmp
gitcomment: "💾sync($filename):$fullfilename"
```