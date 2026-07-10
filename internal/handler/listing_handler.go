package handler

import (
	"fmt"
	"realestate/internal/domain"
	"realestate/internal/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type ListingHandler struct {
    service *service.ListingService
}

func NewListingHandler(s *service.ListingService) *ListingHandler {
    return &ListingHandler{service: s}
}

// Search godoc
// @Summary      Поиск объявлений
// @Description  Все параметры необязательны. Без параметров — вернёт последние 50.
// @Tags         listings
// @Produce      json
// @Param        city       query  string  false  "Город (поиск по вхождению)"    example(Москва)
// @Param        type       query  string  false  "Тип сделки: sale или rent"
// @Param        property   query  string  false  "Тип жилья: apartment,house,commercial"
// @Param        min_price  query  int     false  "Цена от (манат)"
// @Param        max_price  query  int     false  "Цена до (манат)"
// @Param        rooms      query  int     false  "Количество комнат"
// @Success      200        {object}  Response{data=[]domain.Listing}
// @Failure      500        {object}  ErrorResponse
// @Router       /listings [get]
func (h *ListingHandler) Search(c *fiber.Ctx) error {
    filter := domain.ListingFilter{
        City:     c.Query("city"),
        DealType:     c.Query("deal_type"),
        PropertyType: c.Query("property_type"),
        Rooms:    c.QueryInt("rooms"), 
    }
    if p := c.Query("min_price"); p != "" {
        filter.MinPrice, _ = strconv.ParseInt(p, 10, 64)
    }
    if p := c.Query("max_price"); p != "" {
        filter.MaxPrice, _ = strconv.ParseInt(p, 10, 64)
    }

    listings, err := h.service.Search(filter)
    if err != nil {
        return internalError(c)
    }
    if listings == nil {
        listings = []domain.Listing{}
    }
    return ok(c, listings)
}

// GetByID godoc
// @Summary      Получить объявление по ID
// @Tags         listings
// @Produce      json
// @Param        id   path      int  true  "ID объявления"  example(1)
// @Success      200  {object}  Response{data=domain.Listing}
// @Failure      404  {object}  ErrorResponse
// @Router       /listings/{id} [get]
func (h *ListingHandler) GetByID(c *fiber.Ctx) error {
    id, err := c.ParamsInt("id")
    if err != nil {
        return badRequest(c, "id должен быть числом")
    }

    listing, err := h.service.GetByID(id)
    if err != nil {
        fmt.Println("ERROR:", err)
        return c.Status(500).JSON(fiber.Map{
            "error": err.Error(),
        })
    }

    return ok(c, listing)
}

// Create godoc
// @Summary      Создать объявление
// @Description  Требует JWT. user_id берётся из токена — не из тела запроса.
// @Tags         listings
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        input  body      domain.Listing  true  "Данные объявления"
// @Success      201    {object}  Response{data=domain.Listing}
// @Failure      400    {object}  ErrorResponse
// @Failure      401    {object}  ErrorResponse
// @Failure      422    {object}  ErrorResponse
// @Router       /listings [post]
func (h *ListingHandler) Create(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(int)

    var input domain.Listing
    if err := c.BodyParser(&input); err != nil {
        return badRequest(c, "неверный формат запроса")
    }
    input.UserID = userID

    listing, err := h.service.Create(&input)
    if err != nil {
        return unprocessable(c, err.Error())
    }
    return created(c, listing)
}

// Update godoc
// @Summary      Изменить объявление
// @Description  Только автор может редактировать своё объявление.
// @Tags         listings
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id     path      int             true  "ID объявления"
// @Param        input  body      domain.Listing  true  "Новые данные"
// @Success      200    {object}  Response{data=domain.Listing}
// @Failure      403    {object}  ErrorResponse
// @Router       /listings/{id} [put]
func (h *ListingHandler) Update(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(int)
    id, err := c.ParamsInt("id")
    if err != nil {
        return badRequest(c, "id должен быть числом")
    }
    var input domain.Listing
    if err := c.BodyParser(&input); err != nil {
        return badRequest(c, "неверный формат запроса")
    }
    input.ID = id
    if err := h.service.Update(&input, userID); err != nil {
        return forbidden(c, err.Error())
    }
    return ok(c, input)
}

// Delete godoc
// @Summary      Удалить объявление
// @Description  Мягкое удаление — скрывает объявление. Данные остаются в БД.
// @Tags         listings
// @Security     BearerAuth
// @Param        id   path  int  true  "ID объявления"
// @Success      204  "Удалено"
// @Failure      403  {object}  ErrorResponse
// @Router       /listings/{id} [delete]
func (h *ListingHandler) Delete(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(int)
    id, err := c.ParamsInt("id")
    if err != nil {
        return badRequest(c, "id должен быть числом")
    }
    if err := h.service.Delete(id, userID); err != nil {
        return forbidden(c, err.Error())
    }
    return noContent(c)
}