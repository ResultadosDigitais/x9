package handler

import (
	"fmt"
	"net/http"

	"github.com/ResultadosDigitais/x9/cmd/app/auth"
	"github.com/ResultadosDigitais/x9/log"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := session.Get("x9-session", c)
		if err != nil {
			log.Error("Error on getting session", map[string]interface{}{"error": err.Error()})
			return c.JSON(http.StatusInternalServerError, nil)
		}
		token := sess.Values["token"]
		if token != nil {
			oidcToken, err := auth.VerifyToken(token.(string))
			if err == nil {
				userInfo := struct {
					Email string `json:"email"`
					Name  string `json:"name"`
				}{"", ""}
				oidcToken.Claims(&userInfo)
				log.Info(fmt.Sprintf("User %s requested %s", userInfo.Email, c.Path()), nil)
				c.Set("email", userInfo.Email)
				c.Set("name", userInfo.Name)
				if err := next(c); err != nil {
					log.Error("Error handling function", map[string]interface{}{"error": err.Error()})
					return c.JSON(http.StatusInternalServerError, nil)
				}
				return nil
			}
			log.Error("Error on verifying user token", map[string]interface{}{"error": err.Error()})

		} else {
			log.Error("Missing token in user session", nil)
		}

		return c.Redirect(http.StatusFound, "/login")
	}
}
