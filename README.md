# basher (猛击者)

# usage

默认每次请求都生成不同的userAgent.

```shell
./basher --help
Usage:
  -c int
        concurrent (default 8) 并发线程数
  -debug
        debug
  -fake
        Random X-Forwarded-For and X-Real-IP (default true)
  -s string
        target url (default "https://proof.ovh.net/files/1Mb.dat")
```

实例:
运行1个线程下载，打开debug日志:

`./basher -c 1 -s https://proof.ovh.net/files/1Mb.dat -debug`

一些测试下载地址:

`https://proof.ovh.net/files/1Mb.dat`

`https://proof.ovh.net/files/10Mb.dat`

`https://proof.ovh.net/files/1Gb.dat`





