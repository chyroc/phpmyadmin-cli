package phpmyadmin

import (
	"testing"
	"github.com/Chyroc/phpmyadmin-cli/internal/requests"
)

func initClient() *phpMyAdmin {
	s, err := requests.NewSession()
	if err != nil {
		panic(err)
	}
	return &phpMyAdmin{
		Session: s,
	}
}

func TestLogin(t *testing.T) {
	p := initClient()
	p.SetURI("10.102.3.114:8000")

	if err := p.Login("root", "pass"); err != nil {
		t.Errorf(err.Error())
	}

	if err := p.Login("root", "error"); err == nil {
		t.Errorf("should return err")
	}
}
