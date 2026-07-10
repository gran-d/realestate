package handler

import (
	"os"
	"realestate/internal/service"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{service: s}
}

// Register godoc
// @Summary      Регистрация
// @Description  Создаёт нового пользователя. Email должен быть уникальным.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input  body      RegisterInput              true  "Данные"
// @Success      201    {object}  Response{data=domain.User}
// @Failure      400    {object}  ErrorResponse
// @Failure      422    {object}  ErrorResponse
// @Router       /auth/register [post]
func (h *UserHandler) Register(c *fiber.Ctx) error {
	var input RegisterInput

	if err := c.BodyParser(&input); err != nil {
		return badRequest(c, "неверный формат запроса")
	}

	if input.Email == "" || input.Password == "" || input.Name == "" {
		return badRequest(c, "email, password и name обязательны")
	}

	user, err := h.service.Register(input.Email, input.Password, input.Name)
	if err != nil {
		return unprocessable(c, err.Error())
	}

	return created(c, user)
}

// Login godoc
// @Summary      Вход
// @Description  Авторизация пользователя
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input  body      LoginInput  true  "Credentials"
// @Success      200    {object}  Response
// @Failure      401    {object}  ErrorResponse
// @Router       /auth/login [post]
func (h *UserHandler) Login(c *fiber.Ctx) error {
	var input LoginInput

	if err := c.BodyParser(&input); err != nil {
		return badRequest(c, "неверный формат запроса")
	}

	if input.Email == "" || input.Password == "" {
		return badRequest(c, "email и password обязательны")
	}

	resp, err := h.service.Login(input.Email, input.Password)
	if err != nil {
		return unauthorized(c, err.Error())
	}

	return ok(c, resp)
}

// Refresh godoc
// @Summary      Обновить access token
// @Description  Возвращает новый access token по refresh token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input  body      object  true  "Refresh Token"
// @Success      200    {object}  Response
// @Failure      401    {object}  ErrorResponse
// @Router       /auth/refresh [post]
func (h *UserHandler) Refresh(c *fiber.Ctx) error {
	var input struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.BodyParser(&input); err != nil {
		return badRequest(c, "неверный формат запроса")
	}

	if input.RefreshToken == "" {
		return badRequest(c, "refresh_token обязателен")
	}

	accessToken, err := h.service.Refresh(input.RefreshToken)
	if err != nil {
		return unauthorized(c, err.Error())
	}

	return ok(c, fiber.Map{
		"access_token": accessToken,
	})
}

// Logout godoc
// @Summary      Выход
// @Description  Удаляет refresh token
// @Tags         auth
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        input  body      object  true  "Refresh Token"
// @Success      204
// @Router       /auth/logout [post]
func (h *UserHandler) Logout(c *fiber.Ctx) error {
	var input struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.BodyParser(&input); err != nil {
		return badRequest(c, "неверный формат запроса")
	}

	if input.RefreshToken == "" {
		return badRequest(c, "refresh_token обязателен")
	}

	if err := h.service.Logout(input.RefreshToken); err != nil {
		return badRequest(c, err.Error())
	}

	return noContent(c)
}

// RegisterAsAgent godoc
// @Summary      Регистрация как агент
// @Description  Создаёт аккаунт агента
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input  body      AgentRegisterInput  true  "Данные агента"
// @Success      201    {object}  Response
// @Router       /auth/register-agent [post]
func (h *UserHandler) RegisterAsAgent(c *fiber.Ctx) error {
	var input struct {
		Email      string `json:"email"`
		Password   string `json:"password"`
		Name       string `json:"name"`
		Phone      string `json:"phone"`
		InviteCode string `json:"invite_code"`
	}

	if err := c.BodyParser(&input); err != nil {
		return badRequest(c, "неверный формат")
	}

	if input.InviteCode != os.Getenv("AGENT_INVITE_CODE") {
		return forbidden(c, "неверный код приглашения")
	}

	user, err := h.service.RegisterWithRole(
		input.Email,
		input.Password,
		input.Name,
		input.Phone,
		"agent",
	)
	if err != nil {
		return unprocessable(c, err.Error())
	}

	return created(c, user)
}