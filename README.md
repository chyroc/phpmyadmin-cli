# phpmyadmin-cli

[![Build Status](https://travis-ci.org/Chyroc/phpmyadmin-cli.svg?branch=master)](https://travis-ci.org/Chyroc/phpmyadmin-cli)
[![Go Report Card](https://goreportcard.com/badge/github.com/Chyroc/phpmyadmin-cli)](https://goreportcard.com/report/github.com/Chyroc/phpmyadmin-cli)
[![codecov](https://codecov.io/gh/Chyroc/phpmyadmin-cli/branch/master/graph/badge.svg)](https://codecov.io/gh/Chyroc/phpmyadmin-cli)

access phpmyadmin from cli / 通过shell操作phpmyadmin

## features
* access phpmyadmin from cli
* grammar tip

## install
```
go get github.com/Chyroc/phpmyadmin-cli
```

## use

```
➜  ~ phpmyadmin-cli -h
NAME:
   phpmyadmin-cli - access phpmyadmin from shell cli

USAGE:
   phpmyadmin-cli [global options] [arguments...]

GLOBAL OPTIONS:
   -host           phpMyAdmin host
   -port           phpMyAdmin port
   -server         选择server
   -username       phpMyAdmin用户名（为空则跳过验证）
   -password       phpMyAdmin密码
   -history        command history file (default: "~/.phpmyadmin_cli_history")
   -log            command log file (default: "~/.phpmyadmin_cli.log")
   -v              开启调试信息 v
   -vv             开启调试信息 vv
   -vvv            开启调试信息 vvv

   -list           获取server列表
   -prune          清理命令记录
   -h              show help
```

### 多server，命令

* 显示servers: `show servers;`
* 选择一个server: `set server_id;`
