package phpmyadmin

import "testing"

func TestLogin(t *testing.T) {
	DefaultPHPMyAdmin.SetURI("localhost:8000")
	err := DefaultPHPMyAdmin.Login("root", "pass")
	if err != nil {
		t.Errorf(err.Error())
	}
}
