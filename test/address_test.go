package test

import (
	address2 "booktastic-server-go/address"
	json2 "encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestAddress(t *testing.T) {
	// Get logged out.
	resp, _ := getApp().Test(httptest.NewRequest("GET", "/api/address", nil))
	assert.Equal(t, 401, resp.StatusCode)

	user, token := GetUserWithToken(t)

	resp, _ = getApp().Test(httptest.NewRequest("GET", "/api/address?jwt="+token, nil))
	assert.Equal(t, 200, resp.StatusCode)

	var addresses []address2.Address
	json2.Unmarshal(rsp(resp), &addresses)
	assert.Greater(t, len(addresses), 0)
	assert.Equal(t, addresses[0].Userid, user.ID)

	// Get by id
	idstr := strconv.FormatUint(addresses[0].ID, 10)
	resp, _ = getApp().Test(httptest.NewRequest("GET", "/api/address/"+idstr+"?jwt="+token, nil))
	assert.Equal(t, 200, resp.StatusCode)
	var address address2.Address
	json2.Unmarshal(rsp(resp), &address)
	assert.Equal(t, address.ID, addresses[0].ID)
	assert.Equal(t, address.Userid, user.ID)

	// Invalid id.
	resp, _ = getApp().Test(httptest.NewRequest("GET", "/api/address/0?jwt="+token, nil))
	assert.Equal(t, 404, resp.StatusCode)

	// Without token
	resp, _ = getApp().Test(httptest.NewRequest("GET", "/api/address/"+idstr, nil))
	assert.Equal(t, 404, resp.StatusCode)
}
