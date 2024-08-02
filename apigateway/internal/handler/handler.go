package handler

import (
	"github.com/alextotalk/atanika/internal/domain"
	"github.com/alextotalk/atanika/internal/service"
	"github.com/labstack/echo/v4"
	"html/template"
	"net/http"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

var tpl = template.Must(template.ParseFiles("a"))

func indexHandler(c echo.Context) error {
	// Pass any necessary data to the template (optional)
	users := []domain.User{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
	}

	data := map[string]interface{}{
		"Users": users}
	return tpl.Execute(c.Response().Writer, data) // Write to context's response writer

}

func auth(c echo.Context) error {

	return c.String(http.StatusOK, "Authorized")
}

func (h *Handler) NewRouter() *echo.Echo {
	e := echo.New()

	e.GET("/", indexHandler)
	e.GET("/auth", auth)
	return e
}
