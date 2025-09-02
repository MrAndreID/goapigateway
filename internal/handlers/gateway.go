package handlers

import (
	"strings"

	"github.com/MrAndreID/goapigateway/internal/services"

	"github.com/labstack/echo/v4"
	"github.com/spf13/cast"
)

type gatewayHandler struct {
	GatewayService services.IGatewayService
}

func NewGatewayHandler(e *echo.Echo, gatewayService services.IGatewayService) *gatewayHandler {
	handler := &gatewayHandler{
		GatewayService: gatewayService,
	}

	e.POST(":path", handler.Gateway)
	e.GET(":path", handler.Gateway)
	e.PATCH(":path", handler.Gateway)
	e.PUT(":path", handler.Gateway)
	e.DELETE(":path", handler.Gateway)

	return handler
}

func (h *gatewayHandler) Gateway(c echo.Context) error {
	return c.JSON(h.GatewayService.Gateway(services.GatewayRequest{
		Method:    c.Request().Method,
		URL:       c.Request().URL.String(),
		Headers:   c.Request().Header,
		Body:      c.Request().Body,
		IPAddress: strings.Split(c.Request().RemoteAddr, ":")[0],
		Host:      c.Request().Host,
		RequestID: cast.ToString(c.Get("RequestID")),
	}))
}
