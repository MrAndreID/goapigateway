package services

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/MrAndreID/goapigateway/internal/repositories"
	"github.com/MrAndreID/goapigateway/internal/types"

	"github.com/spf13/cast"
)

type IGatewayService interface {
	Gateway(gatewayRequest GatewayRequest) (int, any)
}

type GatewayRequest struct {
	Method    string
	URL       string
	Headers   map[string][]string
	Body      any
	IPAddress string
	Host      string
	RequestID string
}

type GatewayService struct {
	HeaderPrefix      string
	GatewayRepository repositories.IGatewayRepository
}

func NewGatewayService(headerPrefix string, gatewayRepository repositories.IGatewayRepository) *GatewayService {
	return &GatewayService{
		HeaderPrefix:      headerPrefix,
		GatewayRepository: gatewayRepository,
	}
}

func (s *GatewayService) Gateway(gatewayRequest GatewayRequest) (int, any) {
	var (
		statusCode   int
		responseBody any
		gatewayData  = map[string][]map[string]any{
			"POST":   types.Post,
			"GET":    types.Get,
			"PATCH":  types.Patch,
			"PUT":    types.Put,
			"DELETE": types.Delete,
		}
	)

	value, exist := gatewayData[gatewayRequest.Method]

	if exist {
		for _, gatewayValue := range value {
			endpointData := strings.Split(cast.ToString(gatewayValue["endpoint"]), "/")

			urlData := strings.Split(gatewayRequest.URL, "?")

			pathData := strings.Split(urlData[0], "/")

			if len(endpointData) != len(pathData) {
				continue
			}

			var (
				nextData bool
				param    []map[string]string
			)

			for i, v := range endpointData {
				if strings.Index(v, "{") == 0 && strings.Index(v, "}") == len(v)-1 {
					param = append(param, map[string]string{
						"id":    v,
						"value": pathData[i],
					})

					continue
				}

				if v != pathData[i] {
					nextData = true

					break
				}
			}

			if nextData {
				continue
			}

			var (
				methodValue, urlValue string
				headerValue           map[string][]string
			)

			backendData := cast.ToStringMapString(gatewayValue["backend"])

			if backendData["method"] != "" {
				methodValue = backendData["method"]
			} else {
				methodValue = gatewayRequest.Method
			}

			if backendData["path"] == "" {
				backendData["path"] = cast.ToString(gatewayValue["endpoint"])
			}

			if len(urlData) > 1 {
				urlValue = backendData["host"] + backendData["path"] + "?" + urlData[1]
			} else {
				urlValue = backendData["host"] + backendData["path"]
			}

			if len(param) > 0 {
				for _, v := range param {
					urlValue = strings.Replace(urlValue, v["id"], v["value"], -1)
				}
			}

			headerValue = gatewayRequest.Headers

			headerValue[s.HeaderPrefix+"Ip-Address"] = []string{gatewayRequest.IPAddress}

			headerValue[s.HeaderPrefix+"Host"] = []string{gatewayRequest.Host}

			headerValue[s.HeaderPrefix+"Request-ID"] = []string{gatewayRequest.RequestID}

			statusCode, responseBody = s.GatewayRepository.SendToBackend(0, methodValue, urlValue, headerValue, gatewayRequest.Body)
		}
	}

	if !exist || responseBody == nil {
		return http.StatusNotFound, types.MainResponse{
			Code:        fmt.Sprintf("%04d", http.StatusNotFound),
			Description: strings.ToUpper(strings.ReplaceAll(http.StatusText(http.StatusNotFound), " ", "_")),
		}
	}

	return statusCode, responseBody
}
