package server

import (
	"encoding/json"
	"net/http"
	"sha/backend/api"

	"github.com/labstack/echo/v4"
)

func (s *Server) registerSystemRoutes(g *echo.Group) {
	g.GET("/ping", func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		if err := json.NewEncoder(c.Response().Writer).Encode(composeResponse(s.Profile)); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to compose system profile").SetInternal(err)
		}

		return nil
	})

	g.GET("/status", func(c echo.Context) error {
		ownerUserType := api.Owner
		ownerUserFind := api.UserFind{
			Role: &ownerUserType,
		}
		ownerUser, err := s.Store.FindUser(&ownerUserFind)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to find owner user").SetInternal(err)
		}

		systemStatus := api.SystemStatus{
			Owner:   ownerUser,
			Profile: s.Profile,
		}

		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		if err := json.NewEncoder(c.Response().Writer).Encode(composeResponse(systemStatus)); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to encode system status response").SetInternal(err)
		}

		return nil
	})
}
