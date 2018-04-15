package common

import "fmt"

var IsDebug1 = false
var IsDebug2 = false
var IsDebug3 = false // http response

var ErrNeedLogin = fmt.Errorf("need login")
var ErrLoginFailed = fmt.Errorf("login failed")
