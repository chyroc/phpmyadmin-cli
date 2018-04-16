package testdata

import (
	"os/exec"
	"bytes"

	"github.com/stretchr/testify/suite"
	"strings"
	"testing"
)

var helpInfo = `NAME:
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
   -h              show help`

type Cli struct {
	suite.Suite
	t *testing.T

	stdout  *bytes.Buffer
	stderr  *bytes.Buffer
	command []string
	bin     string

	expectStdout interface{}
	expectStderr interface{}
}

func (t *Cli) addCommand(s ...string) {
	t.command = append(t.command, s...)
}

func (t *Cli) SetupTest() {
	t.stdout = new(bytes.Buffer)
	t.stderr = new(bytes.Buffer)
	t.command = nil
	t.expectStderr = nil
	t.expectStdout = nil
}

func (t *Cli) TearDownTest() {
	if len(t.command) > 0 {
		c := exec.Command(t.command[0], t.command[1:]...)
		c.Stdout = t.stdout
		c.Stderr = t.stderr
		t.Nil(c.Run())
		if t.expectStdout != nil {
			t.Equal(t.expectStdout.(string), t.stdout.String())
		}
		if t.expectStderr != nil {
			t.Equal(t.expectStderr.(string), t.stderr.String())
		}
	}
}

func (t *Cli) TestFindBin() {
	t.addCommand()

	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	c := exec.Command("which", "phpmyadmin-cli")
	c.Stdout = stdout
	c.Stderr = stderr
	t.Nil(c.Run())
	t.bin = strings.Replace(stdout.String(), "\n", "", -1)
	t.t.Logf("bin %s\n", t.bin)
}

func (t *Cli) TestHelp() {
	t.addCommand(t.bin)
	t.expectStdout = helpInfo
}

func (t *Cli) TestHelp2() {
	t.addCommand(t.bin, "-h")
	t.expectStdout = helpInfo
}

func TestCli(t *testing.T) {
	suite.Run(t, &Cli{
		t: t,
	})
}
