package router

import (
	"github.com/gofiber/fiber/v2"

	"FinalTask/internal/handler"
	"FinalTask/internal/repository"
	"FinalTask/internal/service"
)

func SetupRoutes(app *fiber.App) {
	// ===== Repository Layer =====
	userRepo := repository.NewUserRepository()
	storeRepo := repository.NewStoreRepository()
	addressRepo := repository.NewAddressRepository()
	categoryRepo := repository.NewCategoryRepository()
	productRepo := repository.NewProductRepository()
	trxRepo := repository.NewTransactionRepository()

	// ===== Service Layer =====
	authService := service.NewAuthService(userRepo, storeRepo) // contoh: auth butuh user & store
	userService := service.NewUserService(userRepo)
	storeService := service.NewStoreService(storeRepo)
	addressService := service.NewAddressService(addressRepo)
	categoryService := service.NewCategoryService(categoryRepo)
	productService := service.NewProductService(productRepo, storeRepo, categoryRepo)
	trxService := service.NewTransactionService(trxRepo, productRepo)

	// ===== Handler Layer =====
	api := app.Group("/api/v1")

	handler.NewAuthHandler(api, authService)
	handler.NewUserHandler(api, userService)
	handler.NewStoreHandler(api, storeService)
	handler.NewAddressHandler(api, addressService)
	handler.NewCategoryHandler(api, categoryService)
	handler.NewProductHandler(api, productService)
	handler.NewTransactionHandler(api, trxService)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("API FinalTask Rakamin is running 🚀")
	})
}
