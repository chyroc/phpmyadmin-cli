# phpmyadmin-cli
access phpmyadmin from cli / 通过shell操作phpmyadmin

## install
```
go get github.com/Chyroc/phpmyadmin-cli
```

## use

### help
```
➜  ~ phpmyadmin-cli -h
NAME:
   phpmyadmin-cli - access phpmyadmin from shell cli

USAGE:
   phpmyadmin-cli [global options] [arguments...]

GLOBAL OPTIONS:
   --url value      phpmyadmin url
   --prune          清理命令记录
   --history value  command history file (default: "~/.phpmyadmin_cli_history")
   --help, -h       show help
```

### connect to phpadmin
```
➜  ~ phpmyadmin-cli -url ip:port
>>>
```
