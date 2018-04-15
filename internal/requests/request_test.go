package requests

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequest(t *testing.T) {
	as := assert.New(t)

	session, err := NewSession()
	as.Nil(err)
	as.NotNil(session)

	{
		resp, err := request(session, http.MethodGet, "localhost:8000", "", nil, nil, nil)
		as.Nil(err)
		as.NotNil(resp)
		defer resp.Body.Close()

		b, err := ioutil.ReadAll(resp.Body)
		as.Nil(err)
		as.NotNil(b)
	}

	{
		resp, err := request(session, http.MethodGet, "http://localhost:8000", "", nil, nil, nil)
		as.Nil(err)
		as.NotNil(resp)
		defer resp.Body.Close()

		b, err := ioutil.ReadAll(resp.Body)
		as.Nil(err)
		as.NotNil(b)
	}
}
