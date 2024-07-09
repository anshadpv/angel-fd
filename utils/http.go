package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/angel-one/fd-core/commons/config"
	fdcontext "github.com/angel-one/fd-core/commons/context"
	"github.com/angel-one/fd-core/commons/httpclient"
	"github.com/angel-one/fd-core/commons/log"
	"github.com/angel-one/fd-core/constants"
	utils2 "github.com/angel-one/go-utils"
	"github.com/angel-one/goerr"
)

func GetBaseUrl(configs map[string]interface{}) string {
	return config.Default().GetConfigEnv(configs, constants.UrlKey)
}

func GetHeaders(configs map[string]interface{}) map[string]string {
	return config.Default().GetConfigOptionMap(configs, "headers")
}

func DoEncodeRequest(ctx context.Context, configKey string, headers map[string]string, httpUrl string, request map[string]string, response interface{}) error {
	httpClient := &http.Client{}

	formData := url.Values{}
	for k, v := range request {
		formData.Add(k, v)
	}

	req, err := http.NewRequest(http.MethodPost, httpUrl, bytes.NewBufferString(formData.Encode()))
	if err != nil {
		return goerr.New(err, "http new request creation failed")
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := httpClient.Do(req)
	if err != nil {
		return goerr.New(err, "http client Do call failed")
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		if res == nil {
			return goerr.New(err, "response is null from client")
		}
		err = utils2.GetJSONData(res.Body, &response)
		if err != nil {
			return goerr.New(err, "http status ok but body content parsing failed")
		}
		log.Debug(ctx).Msgf("http-client request complete")
		return nil
	}

	data, err := utils2.GetDataAsString(res.Body)
	if err != nil {
		return goerr.New(err, "http response body data as string failed")
	}

	log.Debug(ctx).Msgf("http-client request failed. Status: %s; Response: %s", res.Status, data)
	return goerr.New(errors.New(data), "http-client request failed")
}

func DoRequest(ctx context.Context, configKey string, httpClient httpclient.Client, headers map[string]string, url string, queryParams map[string]string, request interface{}, response interface{}) error {
	httpRequest := httpclient.NewRequest(configKey).SetHeaderParams(headers).SetURL(url)

	var bodyBytes *bytes.Buffer
	if request != nil {
		requestJson, err := json.Marshal(request)
		if err != nil {
			return goerr.New(err, "json marshal failed")
		}
		bodyBytes = bytes.NewBuffer(requestJson)
		log.Debug(ctx).Msgf("http-client Payload: %s", string(requestJson))
	}

	if bodyBytes != nil {
		httpRequest.SetBody(bodyBytes)
	}

	if queryParams != nil {
		httpRequest.SetQueryParams(queryParams)
	}

	resp, err := httpClient.Request(httpRequest)
	if err != nil {
		msg := fmt.Sprintf("http-client request failed for clientId %s", GetClientId(ctx))
		return goerr.New(err, msg)
	}
	defer CloseHttpRequest(ctx, resp.Body)

	if resp.StatusCode == http.StatusOK {
		if response == nil {
			return nil
		}
		err = utils2.GetJSONData(resp.Body, &response)
		if err != nil {
			return goerr.New(err, "http status ok but body content parsing failed")
		}
		log.Debug(ctx).Msgf("http-client request complete")
		return nil
	}

	data, err := utils2.GetDataAsString(resp.Body)
	if err != nil {
		return goerr.New(err, "http response body data as string failed")
	}

	log.Debug(ctx).Msgf("http-client request failed. Status: %s; Response: %s", resp.Status, data)
	return goerr.New(errors.New(data), "http-client request failed")
}

func CloseHttpRequest(ctx context.Context, body io.ReadCloser) {
	err := body.Close()
	if err != nil {
		log.Warn(ctx).Stack().Err(err).Msg("error closing request body")
	}
}

func GetClientId(ctx context.Context) string {
	return fdcontext.Get(ctx).UserID
}
