package handler

import (
	"net/http"

	"github.com/ResultadosDigitais/x9/cmd/app/auth"
	"github.com/ResultadosDigitais/x9/config"
	"github.com/ResultadosDigitais/x9/log"
	"github.com/ResultadosDigitais/x9/management"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func (h *Handler) Login(c echo.Context) error {
	return c.Render(http.StatusOK, "login.html", nil)
}

func (h *Handler) OIDCAuth(c echo.Context) error {
	sess, err := session.Get("x9-session", c)
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
	}
	state := auth.GetState()
	sess.Values["state"] = state
	sess.Save(c.Request(), c.Response())
	authURL, err := auth.GetAuthCodeURL(state, nil)
	if err != nil {
		log.Error("Cannot get OIDC auth code url", map[string]interface{}{"error": err.Error()})
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.Redirect(http.StatusFound, authURL)
}

func (h *Handler) OIDCCallback(c echo.Context) error {
	sess, _ := session.Get("x9-session", c)
	state := sess.Values["state"]
	if c.QueryParam("state") != state {
		msg := "Missing state param"
		log.Info(msg, nil)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": msg})
	}
	rawIDToken, err := auth.GetRawIDToken(c.QueryParam("code"))
	if err != nil {
		log.Error("Cannot get raw id token from given code", map[string]interface{}{"error": err.Error()})
		return c.JSON(http.StatusBadRequest, nil)
	}
	_, err = auth.VerifyToken(rawIDToken)
	if err != nil {
		log.Info("Invalid token", map[string]interface{}{"error": err.Error()})
		return c.JSON(http.StatusBadRequest, nil)
	}

	sess.Values["token"] = rawIDToken
	sess.Save(c.Request(), c.Response())

	return c.Redirect(http.StatusFound, "/dashboard")
}

func (h *Handler) Dashboard(c echo.Context) error {
	data := struct {
		AppURL string
	}{config.AppOpts.AppURL}
	return c.Render(http.StatusOK, "dashboard.html", data)
}

func (h *Handler) GetVulnerabilities(c echo.Context) error {

	search := c.QueryParam("search")

	var err error
	var vulns []management.Vulnerability

	if search != "" {
		vulns, err = management.GetVulnerabilitiesByNameOrRepo(search)
	} else {
		vulns, err = management.GetVulnerabilities()
	}
	if err != nil {
		log.Error("Error on getting vulnerabilities", map[string]interface{}{"error": err.Error()})
		return c.JSON(http.StatusBadRequest, nil)
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"rows": vulns, "total": len(vulns)})
}
