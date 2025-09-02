package applications

import (
	"github.com/MrAndreID/goapigateway/internal/repositories"
	"github.com/MrAndreID/goapigateway/internal/services"
)

var (
	GatewayService *services.GatewayService
)

func initService(app *Application) {
	GatewayService = services.NewGatewayService(app.Config.HeaderPrefix, repositories.NewGatewayRepository(app.Config.DefaultTimeout, app.Config.AppDebug))
}
