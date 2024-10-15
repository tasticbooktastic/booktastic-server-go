package test

import (
	json2 "encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/tasticbooktastic/booktastic-server-go/misc"
	"net/http/httptest"
	"testing"
)

func TestMisc(t *testing.T) {
	resp, _ := getApp().Test(httptest.NewRequest("GET", "/api/online", nil))
	assert.Equal(t, 200, resp.StatusCode)

	var result misc.OnlineResult

	json2.Unmarshal(rsp(resp), &result)
	assert.True(t, result.Online)
}
