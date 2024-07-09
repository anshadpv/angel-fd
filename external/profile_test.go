package external

import (
	"fmt"
	"testing"

	"github.com/angel-one/fd-core/commons/config"
	"github.com/angel-one/fd-core/commons/context"
	"github.com/angel-one/fd-core/commons/flags"
	"github.com/angel-one/fd-core/constants"
	httpclient "github.com/angel-one/go-http-client"
	"github.com/stretchr/testify/assert"
)

var ctx = context.Background("test")
var profileService ProfileService
var testClientCode1 = "S1614297"
var testClientCode2 = "B53586"
var testClientCode3 = "XXXXXX"

func init() {
	ctx = context.Background("test")

	configMap := map[string]interface{}{
		"url":             "http://internal-bbe-profile-uat-alb2-474368331.ap-south-1.elb.amazonaws.com/v1/profile/get",
		"method":          "GET",
		"timeoutinmillis": 60000,
		"retrycount":      3,
		"headers": map[string]interface{}{
			"s2sEnabled":    "true",
			"Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyRGF0YSI6eyJjb3VudHJ5X2NvZGUiOiIiLCJtb2Jfbm8iOiIiLCJ1c2VyX2lkIjoiU0dCMi4wLVVBVCIsInNvdXJjZSI6IjM2MWIyODBjLWU0NWQtNDJiNi05ZjE2LWQwNjcyOWI3ZWY0YyIsImFwcF9pZCI6IiIsImNyZWF0ZWRfYXQiOiIyMDIzLTA2LTAxVDA3OjU5OjM5LjE0OTU1ODU5MVoiLCJkYXRhQ2VudGVyIjoiIn0sIm9tbmVtYW5hZ2VyaWQiOjAsInRva2VuIjoiIiwic291cmNlaWQiOiIiLCJ1c2VyX3R5cGUiOiJhcHBsaWNhdGlvbiIsInJlcXVlc3RlZF90b2tlbl90eXBlIjoiYWNjZXNzX3Rva2VuIiwic291cmNlIjoiMzYxYjI4MGMtZTQ1ZC00MmI2LTlmMTYtZDA2NzI5YjdlZjRjIiwiYWN0Ijp7InN1YiI6IjM2MWIyODBjLWU0NWQtNDJiNi05ZjE2LWQwNjcyOWI3ZWY0YyJ9LCJpc3MiOiJBbmdlbEJyb2tpbmciLCJzdWIiOiIzNjFiMjgwYy1lNDVkLTQyYjYtOWYxNi1kMDY3MjliN2VmNGMiLCJhdWQiOlsiMzZiMDc3NGMtZDcyMC00ODQ1LTg1OTAtYjVhMGYyNzk1YjY5Il0sImV4cCI6MTc0OTUzNzMyMiwiaWF0IjoxNjg1NjA2Mzc5LCJqdGkiOiIxMzQxMTZiOC1hMWU0LTQ3ZjMtYTljNi00NmUwOWMwYTQ2YmQifQ.1oc0eOtGgqlRIgFd4bkQjEZjkJliOadzPOlhDOM1nE0",
		},
		"backoffpolicy": map[string]interface{}{
			"constantbackoff": map[string]interface{}{
				"intervalinmillis":          5,
				"maxjitterintervalinmillis": 10,
			},
		},
	}

	requestConfigs := []*httpclient.RequestConfig{httpclient.NewRequestConfig(constants.ProfileServerConfig, configMap)}
	client := httpclient.ConfigureHTTPClient(requestConfigs...)
	config.InitTestMode(fmt.Sprintf("%s/%s", flags.BaseConfigPath(), flags.Env()), constants.HTTPClientConfig)
	profileService = DefaultProfileService(client)
}

func TestGetProfileDetails(t *testing.T) {
	if profileResponse, err := profileService.GetUserProfileDetails(ctx, testClientCode1); err != nil {
		t.Error("get profile server failed", err)
	} else {
		t.Logf("%+v", profileResponse)
		if assert.NotNil(t, profileResponse, "ProfileResponse is NULL") {
			assert.Equal(t, "success", profileResponse.Status, fmt.Sprintf("ProfileServer API call was not success, Reason: %s", profileResponse.Message))
			assert.Equal(t, testClientCode1, profileResponse.Data.ClientID, "Invalid client-code when fetching")
		}
	}
}

func TestGetProfileDetails2(t *testing.T) {
	if profileResponse, err := profileService.GetUserProfileDetails(ctx, testClientCode2); err != nil {
		t.Error("get profile server failed", err)
	} else {
		t.Logf("%+v", profileResponse)
		if assert.NotNil(t, profileResponse, "ProfileResponse is NULL") {
			assert.Equal(t, "success", profileResponse.Status, fmt.Sprintf("ProfileServer API call was not success, Reason: %s", profileResponse.Message))
			assert.Equal(t, testClientCode2, profileResponse.Data.ClientID, "Invalid client-code when fetching")
		}
	}
}

func TestGetProfileDetails3(t *testing.T) {
	if profileResponse, err := profileService.GetUserProfileDetails(ctx, testClientCode3); err != nil {
		t.Error("get profile server failed", err)
	} else {
		t.Logf("%+v", profileResponse)
		if assert.NotNil(t, profileResponse, "ProfileResponse is NULL") {
			assert.Equal(t, "error", profileResponse.Status, fmt.Sprintf("ProfileServer API call was not error, Reason: %s", profileResponse.Message))
			assert.Equal(t, "data not found", profileResponse.Message, fmt.Sprintf("ProfileServer API call was not right, Reason: %s", profileResponse.Message))
			assert.Equal(t, "", profileResponse.Data.ClientID, "Invalid client-code when fetching")
		}
	}
}
