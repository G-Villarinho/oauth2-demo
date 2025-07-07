package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/aetheris-lab/aetheris-id/api/internal/domain"
	"github.com/aetheris-lab/aetheris-id/api/internal/models"
	"github.com/aetheris-lab/aetheris-id/api/internal/services"
	"github.com/labstack/echo/v4"
)

type ClientHandler interface {
	CreateClient(ectx echo.Context) error
}

type clientHandler struct {
	clientService services.ClientService
}

func NewClientHandler(clientService services.ClientService) ClientHandler {
	return &clientHandler{
		clientService: clientService,
	}
}

func (h *clientHandler) CreateClient(ectx echo.Context) error {
	logger := slog.With(
		slog.String("handler", "client"),
		slog.String("method", ectx.Request().Method),
		slog.String("path", ectx.Request().URL.Path),
	)

	var payload models.CreateClientPayload
	if err := ectx.Bind(&payload); err != nil {
		logger.Error("bind payload", "error", err)
		return echo.ErrBadRequest
	}

	if err := ectx.Validate(&payload); err != nil {
		logger.Error("validate payload", "error", err)
		return err
	}

	response, err := h.clientService.CreateClient(ectx.Request().Context(), payload.Name, payload.Description, payload.RedirectURIs, payload.GrantTypes)
	if err != nil {
		if errors.Is(err, domain.ErrClientAlreadyExists) {
			logger.Error(err.Error())
			return echo.ErrConflict
		}

		logger.Error("create client", "error", err)
		return echo.ErrInternalServerError
	}

	return ectx.JSON(http.StatusCreated, response)
}
