package external

import (
	"context"
	"strings"
	"time"

	"github.com/angel-one/fd-core/business/model"
	"github.com/angel-one/fd-core/commons/config"
	"github.com/angel-one/fd-core/commons/httpclient"
	"github.com/angel-one/fd-core/commons/log"
	"github.com/angel-one/fd-core/constants"
	"github.com/angel-one/fd-core/utils"
	"github.com/angel-one/goerr"
)

const (
	clientId     = "client_id"
	grantType    = "grant_type"
	clientSecret = "client_secret"
	scope        = "scope"
)

var token string
var tokenExpiresIn time.Time

type UpSwing interface {
	ValidateToken(ctx context.Context) error
	DoPCIRegistration(ctx context.Context, clientCode string) (*model.PCIRegistrationResponse, error)
	GetNetWorthData(ctx context.Context, clientCode string) (*model.NetWorthResponse, error)
	GetPendingJourneyData(ctx context.Context, clientCode string) (*model.PendingJourneyResponse, error)
	generateAccessToken(ctx context.Context) error
}

type upSwingImpl struct {
	httpClient     httpclient.Client
	profileService ProfileService
	tokenPayload   map[string]string
}

func DefaultUpSwing(ctx context.Context) UpSwing {
	secrets := config.Default().Secrets

	tokenPayload := make(map[string]string)
	tokenPayload[clientId] = secrets[constants.UpswingClientId]
	tokenPayload[grantType] = secrets[constants.UpswingGrantType]
	tokenPayload[clientSecret] = secrets[constants.UpswingClientSecret]
	tokenPayload[scope] = secrets[constants.UpswingScope]

	u := &upSwingImpl{tokenPayload: tokenPayload, httpClient: httpclient.Default(), profileService: DefaultProfileService(httpclient.Default())}
	u.generateAccessToken(ctx)
	return u
}

func (i *upSwingImpl) getHeaders(configs map[string]interface{}, appendToken bool) map[string]string {
	headers := utils.GetHeaders(configs)
	if appendToken {
		headers[constants.HeaderAuthorization] = constants.HeaderAuthorizationBearer + " " + token
	}
	return headers
}

func (u *upSwingImpl) generateAccessToken(ctx context.Context) error {

	// todo: token reuse ; validate its expiry and invoke only if its expired
	response := model.GenerateTokenResponse{}
	configs, _ := config.Default().GetMap(constants.HTTPClientConfig, constants.UpSwingGenerateToken)
	err := utils.DoEncodeRequest(ctx, constants.UpSwingGenerateToken, u.getHeaders(configs, false), utils.GetBaseUrl(configs), u.tokenPayload, &response)
	if err != nil {
		return goerr.New(err, "external failed : failed to generate upswing token")
	}
	log.Info(ctx).Msgf("upswing token generation is complete")

	token = response.AccessToken
	tokenExpiresIn = time.Now().Add(time.Second * (time.Duration(response.ExpiresIn) - 60))
	return nil
}

func (u *upSwingImpl) DoPCIRegistration(ctx context.Context, clientCode string) (*model.PCIRegistrationResponse, error) {

	pciResponse := model.PCIRegistrationResponse{}
	pciRequest := model.PCIRegistrationRequest{PartnerCustomerId: clientCode}
	pciConfigs, _ := config.Default().GetMap(constants.HTTPClientConfig, constants.UpSwingPCIRegistration)
	err := utils.DoRequest(ctx, constants.UpSwingPCIRegistration, u.httpClient, u.getHeaders(pciConfigs, true), utils.GetBaseUrl(pciConfigs), nil, pciRequest, &pciResponse)
	if err != nil {
		return nil, goerr.New(err, "external failed : failed to do upswing PCI registration")
	}
	log.Info(ctx).Msgf("upswing user registration compelte for client-code: %s", clientCode)

	// do data-ingestion via goroutine
	profileData, err := u.profileService.GetUserProfileDetails(ctx, clientCode)
	if err != nil {
		return nil, err
	}
	dataIngestionRequest := model.DataIngestionRequest{Pan: profileData.Data.Pan}
	err = u.postDataIngestion(ctx, clientCode, dataIngestionRequest)
	if err != nil {
		return nil, err
	}
	log.Info(ctx).Msgf("upswing data ingestion compelte for client-code: %s", clientCode)

	return &pciResponse, err
}

func (u *upSwingImpl) postDataIngestion(ctx context.Context, clientCode string, request model.DataIngestionRequest) error {
	configs, _ := config.Default().GetMap(constants.HTTPClientConfig, constants.UpswingDataIngestion)
	url := strings.Replace(utils.GetBaseUrl(configs), constants.UpPCIField, clientCode, -1)
	err := utils.DoRequest(ctx, constants.UpswingDataIngestion, u.httpClient, u.getHeaders(configs, true), url, nil, request, nil)
	if err != nil {
		return goerr.New(err, "external failed : data ingestion call failed with upswing")
	}
	return nil
}

func (u *upSwingImpl) GetNetWorthData(ctx context.Context, clientCode string) (*model.NetWorthResponse, error) {
	response := model.NetWorthResponse{}
	configs, _ := config.Default().GetMap(constants.HTTPClientConfig, constants.UpSwingNetWorth)
	url := strings.Replace(utils.GetBaseUrl(configs), constants.UpPCIField, clientCode, -1)
	err := utils.DoRequest(ctx, constants.UpSwingNetWorth, u.httpClient, u.getHeaders(configs, true), url, nil, nil, &response)
	if err != nil {
		return nil, goerr.New(err, "external failed : networth api call failed with upswing")
	}
	return &response, nil
}

// refresh the access token if the token expired
func (u *upSwingImpl) ValidateToken(ctx context.Context) error {
	if time.Until(tokenExpiresIn).Seconds() < 1 {
		log.Info(ctx).Msg("upswing token expires in 1 minute - renewing now")
		u.generateAccessToken(ctx)
	} else {
		log.Info(ctx).Msg("upswing token renewal not required")
	}
	return nil
}

func (u *upSwingImpl) GetPendingJourneyData(ctx context.Context, clientCode string) (*model.PendingJourneyResponse, error) {
	response := model.PendingJourneyResponse{}
	configs, _ := config.Default().GetMap(constants.HTTPClientConfig, constants.UpswingPendingJourney)
	url := utils.GetBaseUrl(configs)

	queryParams := make(map[string]string)
	queryParams["pci"] = clientCode
	err := utils.DoRequest(ctx, constants.UpswingPendingJourney, u.httpClient, u.getHeaders(configs, true), url, queryParams, nil, &response)

	if err != nil {
		return nil, goerr.New(err, "external failed : pending journey api call failed with upswing")
	}
	return &response, nil
}
