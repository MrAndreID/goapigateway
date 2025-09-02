package repositories

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/MrAndreID/goapigateway/internal/types"

	"github.com/go-resty/resty/v2"
)

type IGatewayRepository interface {
	SendToBackend(timeout int, method string, url string, headers map[string][]string, body any) (int, any)
}

type GatewayRepository struct {
	DefaultTimeout int
	Debug          bool
}

func NewGatewayRepository(defaultTimeout int, debug bool) *GatewayRepository {
	return &GatewayRepository{
		DefaultTimeout: defaultTimeout,
		Debug:          debug,
	}
}

func (r *GatewayRepository) SendToBackend(timeout int, method string, url string, headers map[string][]string, body any) (int, any) {
	client := resty.New()

	if timeout == 0 {
		timeout = r.DefaultTimeout
	}

	client.SetTimeout(time.Second * time.Duration(timeout))

	request := client.R()

	request.Method = method

	request.URL = url

	request.SetHeaderMultiValues(headers)

	request.SetBody(body)

	if r.Debug {
		request.EnableTrace()

		request.SetDebug(true)

		request.EnableGenerateCurlOnDebug()
	}

	response, err := request.Send()

	if err != nil {
		return http.StatusInternalServerError, types.MainResponse{
			Code:        fmt.Sprintf("%04d", http.StatusInternalServerError),
			Description: strings.ToUpper(strings.ReplaceAll(http.StatusText(http.StatusInternalServerError), " ", "_")),
		}
	}

	var responseBody any

	err = json.Unmarshal(response.Body(), &responseBody)

	if err != nil {
		return http.StatusInternalServerError, types.MainResponse{
			Code:        fmt.Sprintf("%04d", http.StatusInternalServerError),
			Description: strings.ToUpper(strings.ReplaceAll(http.StatusText(http.StatusInternalServerError), " ", "_")),
		}
	}

	return response.StatusCode(), responseBody
}
