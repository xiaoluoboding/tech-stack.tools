package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sha/backend/api"
	"strconv"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func (s *Server) registerUserRoutes(g *echo.Group) {
	g.POST("/user", func(c echo.Context) error {
		userCreate := &api.UserCreate{}
		if err := json.NewDecoder(c.Request().Body).Decode(userCreate); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Malformatted post user request").SetInternal(err)
		}

		passwordHash, err := bcrypt.GenerateFromPassword([]byte(userCreate.Password), bcrypt.DefaultCost)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate password hash").SetInternal(err)
		}

		userCreate.PasswordHash = string(passwordHash)
		user, err := s.Store.CreateUser(userCreate)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create user").SetInternal(err)
		}

		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		if err := json.NewEncoder(c.Response().Writer).Encode(composeResponse(user)); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to encode user response").SetInternal(err)
		}

		return nil
	})

	g.GET("/user", func(c echo.Context) error {
		userList, err := s.Store.FindUserList(&api.UserFind{})
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch user list").SetInternal(err)
		}

		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		if err := json.NewEncoder(c.Response().Writer).Encode(composeResponse(userList)); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to encode user list response").SetInternal(err)
		}

		return nil
	})

	g.GET("/user/me", func(c echo.Context) error {
		userSessionID := c.Get(getUserIDContextKey())
		if userSessionID == nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Missing auth session")
		}

		userID := userSessionID.(int)
		userFind := &api.UserFind{
			ID: &userID,
		}
		user, err := s.Store.FindUser(userFind)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch user").SetInternal(err)
		}

		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		if err := json.NewEncoder(c.Response().Writer).Encode(composeResponse(user)); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to encode user response").SetInternal(err)
		}

		return nil
	})

	g.PATCH("/user/me", func(c echo.Context) error {
		userID := c.Get(getUserIDContextKey()).(int)
		userPatch := &api.UserPatch{
			ID: userID,
		}
		if err := json.NewDecoder(c.Request().Body).Decode(userPatch); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Malformatted patch user request").SetInternal(err)
		}

		if userPatch.Email != nil {
			userFind := api.UserFind{
				Email: userPatch.Email,
			}
			user, err := s.Store.FindUser(&userFind)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to find user by email %s", *userPatch.Email)).SetInternal(err)
			}
			if user != nil {
				return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("User with email %s existed", *userPatch.Email)).SetInternal(err)
			}
		}

		if userPatch.Password != nil && *userPatch.Password != "" {
			passwordHash, err := bcrypt.GenerateFromPassword([]byte(*userPatch.Password), bcrypt.DefaultCost)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate password hash").SetInternal(err)
			}

			passwordHashStr := string(passwordHash)
			userPatch.PasswordHash = &passwordHashStr
		}

		user, err := s.Store.PatchUser(userPatch)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to patch user").SetInternal(err)
		}

		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		if err := json.NewEncoder(c.Response().Writer).Encode(composeResponse(user)); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to encode user response").SetInternal(err)
		}

		return nil
	})

	g.PATCH("/user/:userId", func(c echo.Context) error {
		currentUserID := c.Get(getUserIDContextKey()).(int)
		currentUser, err := s.Store.FindUser(&api.UserFind{
			ID: &currentUserID,
		})
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to find user").SetInternal(err)
		}
		if currentUser == nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Current session user not found with ID: %d", currentUserID)).SetInternal(err)
		} else if currentUser.Role != api.Owner {
			return echo.NewHTTPError(http.StatusForbidden, "Access forbidden for current session user").SetInternal(err)
		}

		userID, err := strconv.Atoi(c.Param("userId"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("ID is not a number: %s", c.Param("userId"))).SetInternal(err)
		}

		userPatch := &api.UserPatch{
			ID: userID,
		}
		if err := json.NewDecoder(c.Request().Body).Decode(userPatch); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Malformatted patch user request").SetInternal(err)
		}

		if userPatch.Password != nil && *userPatch.Password != "" {
			passwordHash, err := bcrypt.GenerateFromPassword([]byte(*userPatch.Password), bcrypt.DefaultCost)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate password hash").SetInternal(err)
			}

			passwordHashStr := string(passwordHash)
			userPatch.PasswordHash = &passwordHashStr
		}

		user, err := s.Store.PatchUser(userPatch)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to patch user").SetInternal(err)
		}

		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		if err := json.NewEncoder(c.Response().Writer).Encode(composeResponse(user)); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to encode user response").SetInternal(err)
		}

		return nil
	})

	g.DELETE("/user/:userId", func(c echo.Context) error {
		userID, err := strconv.Atoi(c.Param("userId"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("ID is not a number: %s", c.Param("userId"))).SetInternal(err)
		}

		userDelete := &api.UserDelete{
			ID: userID,
		}

		err = s.Store.DeleteUser(userDelete)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to delete user ID: %v", userID)).SetInternal(err)
		}

		c.JSON(http.StatusOK, true)

		return nil
	})
}
