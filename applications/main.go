package applications

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/MrAndreID/goapigateway/configs"
	"github.com/MrAndreID/goapigateway/internal/handlers"

	"github.com/MrAndreID/gomiddleware"
	"github.com/MrAndreID/gopackage"
	"github.com/fatih/color"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"github.com/unrolled/secure"
	"go.elastic.co/apm/module/apmechov4"
)

type Application struct {
	Config       *configs.Config
	TimeLocation *time.Location
}

func Start(toggle bool) any {
	var tag string = "Applications.Main.New."

	cfg, err := configs.New(toggle)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"tag":   tag + "01",
			"error": err.Error(),
		}).Error("failed to initiate configuration")

		return nil
	}

	timeLocation, err := time.LoadLocation(cfg.AppLocation)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"tag":   tag + "02",
			"error": err.Error(),
		}).Error("failed to load location for time")

		return nil
	}

	app := Application{
		Config:       cfg,
		TimeLocation: timeLocation,
	}

	echo.NotFoundHandler = func(c echo.Context) error {
		logrus.WithFields(logrus.Fields{
			"tag": tag + "03",
		}).Error("route not found")

		return c.JSON(http.StatusNotFound, map[string]string{
			"code":        fmt.Sprintf("%04d", http.StatusNotFound),
			"description": strings.ToUpper(strings.ReplaceAll(http.StatusText(http.StatusNotFound), " ", "_")),
		})
	}

	echo.MethodNotAllowedHandler = func(c echo.Context) error {
		logrus.WithFields(logrus.Fields{
			"tag": tag + "04",
		}).Error("method not allowed")

		return c.JSON(http.StatusMethodNotAllowed, map[string]string{
			"code":        fmt.Sprintf("%04d", http.StatusMethodNotAllowed),
			"description": strings.ToUpper(strings.ReplaceAll(http.StatusText(http.StatusMethodNotAllowed), " ", "_")),
		})
	}

	e := echo.New()

	e.Validator = gopackage.CustomValidator()

	e.HTTPErrorHandler = gopackage.EchoCustomHTTPErrorHandler

	e.JSONSerializer = gopackage.CustomJSON()

	e.Pre(middleware.RemoveTrailingSlash())

	e.Pre(gomiddleware.EchoSetRequestID)

	e.Use(apmechov4.Middleware())

	e.Use(middleware.Recover())

	e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:         "1; mode=block",
		ContentTypeNosniff:    "nosniff",
		XFrameOptions:         "SAMEORIGIN",
		HSTSMaxAge:            3600,
		ContentSecurityPolicy: "default-src 'self'",
	}))

	secureMiddleware := secure.Options{
		SSLProxyHeaders:      map[string]string{"X-Forwarded-Proto": "https"},
		STSSeconds:           63072000,
		STSIncludeSubdomains: true,
		STSPreload:           true,
		ForceSTSHeader:       true,
		IsDevelopment:        true,
	}

	e.Use(echo.WrapMiddleware(secure.New(secureMiddleware).Handler))

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: cfg.AllowedOrigins,
		AllowHeaders: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))

	e.Use(gomiddleware.EchoSetNoCache)

	e.Use(gomiddleware.EchoSetMaintenanceMode("storages/maintenance.flag"))

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now().In(timeLocation)

			err := next(c)

			stop := time.Now().In(timeLocation)

			dateTime := stop.Format(time.DateTime)

			redBackground := color.New(color.FgWhite).Add(color.BgRed).SprintFunc()
			greenBackground := color.New(color.FgWhite).Add(color.BgGreen).SprintFunc()

			redForeground := color.New(color.FgRed).SprintFunc()
			greenForeground := color.New(color.FgGreen).SprintFunc()

			var status, statusType string

			if c.Response().Status >= 400 {
				status = redBackground(" " + strconv.Itoa(c.Response().Status) + " ")

				statusType = redForeground("[ ERROR ]")
			} else {
				status = greenBackground(" " + strconv.Itoa(c.Response().Status) + " ")

				statusType = greenForeground("[ SUCCESS ]")
			}

			fmt.Printf(
				"%s %s %s %s %s | %s | %s\n",
				dateTime,
				status,
				statusType,
				c.Request().Method,
				c.Request().URL.String(),
				stop.Sub(start).String(),
				cast.ToString(c.Get("RequestID")),
			)

			return err
		}
	})

	if cfg.AppDebug {
		e.Logger.SetLevel(log.DEBUG)

		e.Debug = true
	}

	initService(&app)

	if toggle {
		handlers.NewGatewayHandler(e, GatewayService)

		return e.Start(":" + cfg.AppPort)
	}

	return e
}
