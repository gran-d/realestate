// @title           Real Estate API
// @version         1.0
// @description     Мини-Циан: объявления о недвижимости на Go
// @host            localhost:3000
// @BasePath        /api/v1

// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
// @description                 Вставь: Bearer {токен из /auth/login}

package main

import (
	"database/sql"
	"log"
	"os"
	"time"

	
	_ "realestate/docs"

	"realestate/internal/handler"
	"realestate/internal/middleware"
	"realestate/internal/repository"
	"realestate/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" 
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

func main() {
    
    if err := godotenv.Load(); err != nil {
        log.Println("Файл .env не найден — используем переменные окружения")
    }

    // ── подключение к постгрес ──────────────────────────────────────────────
    db := mustConnectDB()
    defer db.Close()
    
    userRepo     := repository.NewUserRepository(db)
    listingRepo  := repository.NewListingRepository(db)
    favoriteRepo := repository.NewFavoriteRepository(db)
    tokenRepo    := repository.NewTokenRepository(db)

    userSvc := service.NewUserService(
    userRepo,
    tokenRepo,
    os.Getenv("JWT_SECRET"),
)
    listingSvc  := service.NewListingService(listingRepo)
    favoriteSvc := service.NewFavoriteService(favoriteRepo)

    userH     := handler.NewUserHandler(userSvc)
    listingH  := handler.NewListingHandler(listingSvc)
    favoriteH := handler.NewFavoriteHandler(favoriteSvc)

    // ──────────────────────────────────────────────────────
    app := fiber.New(fiber.Config{
        
        ErrorHandler: func(c *fiber.Ctx, err error) error {
            return c.Status(fiber.StatusInternalServerError).
                JSON(fiber.Map{"error": err.Error()})
        },
    })

    // http://localhost:3000/swagger/index.html
    app.Get("/swagger/*", fiberSwagger.WrapHandler)

    api := app.Group("/api/v1")
    

    // public routes  
    api.Post("/auth/register", userH.Register)
    api.Post("/auth/login",    userH.Login)
    api.Get("/listings",       listingH.Search)
    api.Get("/listings/:id",   listingH.GetByID)
    api.Post("/auth/register-agent", userH.RegisterAsAgent)
    api.Post("/auth/refresh", userH.Refresh)
    api.Post("/auth/logout", userH.Logout)
    // private routes
    auth := api.Group("/", middleware.JWTProtected(os.Getenv("JWT_SECRET")))
    auth.Post("/listings",        listingH.Create)
    auth.Put("/listings/:id",     listingH.Update)
    auth.Delete("/listings/:id",  listingH.Delete)
    auth.Get("/favorites",        favoriteH.List)
    auth.Post("/favorites/:id",   favoriteH.Add)
    auth.Delete("/favorites/:id", favoriteH.Remove)

    //────────────────────────────────────────────────────────────────
    port := os.Getenv("APP_PORT")
    if port == "" {
        port = "3000" //default port
    }
    log.Printf("Сервер:   http://localhost:%s", port)
    log.Printf("Swagger:  http://localhost:%s/swagger/index.html", port)
    log.Fatal(app.Listen(":" + port))
    
}

func mustConnectDB() *sql.DB {
    dsn := os.Getenv("DATABASE_URL")
    if dsn == "" {
        log.Fatal("DATABASE_URL не задан в .env")
    }

    for attempt := 1; attempt <= 10; attempt++ {
        db, err := sql.Open("postgres", dsn)
        if err == nil {
            if pingErr := db.Ping(); pingErr == nil {
                log.Printf("PostgreSQL подключён (попытка %d/10)", attempt)
                db.SetMaxOpenConns(25)  
                db.SetMaxIdleConns(5)
                return db
            }
        }
        log.Printf("Postgres не готов, жду... (попытка %d/10)", attempt)
        time.Sleep(2 * time.Second)
    }

    log.Fatal("Не удалось подключиться к PostgreSQL после 10 попыток")
    return nil
}