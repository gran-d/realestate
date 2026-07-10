package handler

import (
	"realestate/internal/domain"
	"realestate/internal/service"

	"github.com/gofiber/fiber/v2"
)

type FavoriteHandler struct {
    service *service.FavoriteService
}

func NewFavoriteHandler(s *service.FavoriteService) *FavoriteHandler {
    return &FavoriteHandler{service: s}
}

// List godoc
// @Summary      Мои избранные
// @Description  Возвращает объявления которые добавил текущий пользователь.
// @Tags         favorites
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  Response{data=[]domain.Listing}
// @Router       /favorites [get]
func (h *FavoriteHandler) List(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(int)
    listings, err := h.service.List(userID)
    if err != nil {
        return internalError(c)
    }
    if listings == nil {
        listings = []domain.Listing{}
    }
    return ok(c, listings)
}

// Add godoc
// @Summary      Добавить в избранное
// @Tags         favorites
// @Security     BearerAuth
// @Param        id   path  int  true  "ID объявления"
// @Success      204  "Добавлено"
// @Failure      422  {object}  ErrorResponse  "Уже в избранном"
// @Router       /favorites/{id} [post]
func (h *FavoriteHandler) Add(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(int)
    listingID, err := c.ParamsInt("id")
    if err != nil {
        return badRequest(c, "id должен быть числом")
    }
    if err := h.service.Add(userID, listingID); err != nil {
        return unprocessable(c, "уже в избранном")
    }
    return noContent(c)
}

// Remove godoc
// @Summary      Убрать из избранного
// @Tags         favorites
// @Security     BearerAuth
// @Param        id   path  int  true  "ID объявления"
// @Success      204  "Убрано"
// @Failure      404  {object}  ErrorResponse
// @Router       /favorites/{id} [delete]
func (h *FavoriteHandler) Remove(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(int)
    listingID, err := c.ParamsInt("id")
    if err != nil {
        return badRequest(c, "id должен быть числом")
    }
    if err := h.service.Remove(userID, listingID); err != nil {
        return notFound(c)
    }
    return noContent(c)
}