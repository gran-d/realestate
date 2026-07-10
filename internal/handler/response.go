package handler

import "github.com/gofiber/fiber/v2"


type Response struct {
    Data interface{} `json:"data"`
}

type ErrorResponse struct {
    Error string `json:"error"`
}

type TokenResponse struct {
    Token string `json:"token" example:"eyJhbGciOiJIUzI1NiJ9..."`
}

type RegisterInput struct {
    Email    string `json:"email"    example:"ivan@example.com"`
    Password string `json:"password" example:"secret123"`
    Name     string `json:"name"     example:"Иван Петров"`
}

type LoginInput struct {
    Email    string `json:"email"    example:"ivan@example.com"`
    Password string `json:"password" example:"secret123"`
}


func ok(c *fiber.Ctx, data interface{}) error {
    return c.JSON(Response{Data: data})
}

// created — 201 Created (object succesfully created)
func created(c *fiber.Ctx, data interface{}) error {
    return c.Status(fiber.StatusCreated).JSON(Response{Data: data})
}

// noContent — 204 No Content (operation completed successfully)
func noContent(c *fiber.Ctx) error {
    return c.SendStatus(fiber.StatusNoContent)
}

// badRequest — 400 (client sent invalid data)
func badRequest(c *fiber.Ctx, msg string) error {
    return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: msg})
}

// unauthorized — 401 (no token)
func unauthorized(c *fiber.Ctx, msg string) error {
    return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{Error: msg})
}

// forbidden — 403 (token is valid but user isnt allowed)
func forbidden(c *fiber.Ctx, msg string) error {
    return c.Status(fiber.StatusForbidden).JSON(ErrorResponse{Error: msg})
}

// notFound — 404 (object not found)
func notFound(c *fiber.Ctx) error {
    return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{Error: "не найдено"})
}

// unprocessable — 422 (data is valid but violates business rules, e.g. email already exists)
func unprocessable(c *fiber.Ctx, msg string) error {
    return c.Status(fiber.StatusUnprocessableEntity).JSON(ErrorResponse{Error: msg})
}

// internalError — 500 (server is bad)
func internalError(c *fiber.Ctx) error {
    return c.Status(fiber.StatusInternalServerError).
        JSON(ErrorResponse{Error: "внутренняя ошибка сервера"})
}