package outbound

import (
	"encoding/json"
	"go-httpix-cli/entity"
	"go-httpix-cli/utils"
	"net/http"
	"time"

	"github.com/imroc/req/v3"
)

type HttpClient struct {
	clientInstance *req.Client
}

func httpClientSetup() *HttpClient {
	reqClient := req.C().
		SetUserAgent("httpix/1.0 TUI-HTTP-Client").
		SetJsonMarshal(json.Marshal).
		SetJsonUnmarshal(json.Unmarshal)

	httpClient := reqClient.GetClient()
	httpClient.Timeout = 1 * time.Minute
	httpClient.Transport = &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
	}

	return &HttpClient{clientInstance: reqClient}
}

func (httpClient *HttpClient) Execute(requestEntity *entity.Request) (responseEntity *entity.Response, err error) {
	buildedUrl, err := utils.BuildURL(requestEntity.URL, requestEntity.Params)
	if err != nil {
		return nil, err
	}
	httpResponse, err := httpClient.clientInstance.R().
		SetBody(requestEntity.Body).
		SetHeaders(utils.ApplyHeaders(requestEntity.Headers)).
		SetContentType(requestEntity.ContentType).
		Send(requestEntity.Method, buildedUrl)
	if err != nil {
		return nil, err
	}
	bodyBytes := httpResponse.Bytes()
	return &entity.Response{
		StatusCode: httpResponse.StatusCode,
		Status:     httpResponse.Status,
		Proto:      httpResponse.Proto,
		Headers:    httpResponse.Header,
		Body:       string(bodyBytes),
		Duration:   httpResponse.TotalTime(),
		Size:       len(httpResponse.Bytes()),
	}, nil
}
