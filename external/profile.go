package external

import (
	"context"

	"github.com/angel-one/fd-core/business/model"
	"github.com/angel-one/fd-core/commons/config"
	c "github.com/angel-one/fd-core/commons/context"
	"github.com/angel-one/fd-core/commons/httpclient"
	"github.com/angel-one/fd-core/commons/log"
	"github.com/angel-one/fd-core/constants"
	"github.com/angel-one/fd-core/utils"
	"github.com/angel-one/goerr"
)

const (
	ClientCodeKey = "clientid"
)

type ProfileService interface {
	GetUserProfileDetails(ctx context.Context, clientCode string) (*model.ProfileResponse, error)
}

type profileServerImpl struct {
	httpClient httpclient.Client
	jwtToken   string
}

func DefaultProfileService(httpclient httpclient.Client) ProfileService {
	jwtToken := config.Default().Secrets[constants.ProfileServiceToken]
	if jwtToken == "" {
		log.Fatal(c.Background("init")).Msgf("Profile service token is missing from configs")
	}
	return &profileServerImpl{httpClient: httpclient, jwtToken: jwtToken}
}

func (p *profileServerImpl) getHeaders(configs map[string]interface{}, appendToken bool) map[string]string {
	headers := utils.GetHeaders(configs)
	if appendToken {
		headers[constants.HeaderAuthorization] = constants.HeaderAuthorizationBearer + " " + p.jwtToken
	}
	return headers
}

func (p *profileServerImpl) GetUserProfileDetails(ctx context.Context, clientCode string) (*model.ProfileResponse, error) {
	response := model.ProfileResponse{}
	configs, _ := config.Default().GetMap(constants.HTTPClientConfig, constants.ProfileServerConfig)
	queryParams := make(map[string]string)
	queryParams[ClientCodeKey] = clientCode
	err := utils.DoRequest(ctx, constants.ProfileServerConfig, p.httpClient, p.getHeaders(configs, true), utils.GetBaseUrl(configs), queryParams, nil, &response)
	if err != nil {
		return nil, goerr.New(err, "external call failed : profile service api invocation failed.")
	}
	return &response, nil
}
