package app

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"cookaholic/internal/infrastructure/cloudinary"
	"cookaholic/internal/infrastructure/db"
	"cookaholic/internal/infrastructure/http"
	"cookaholic/internal/interfaces"
)

// Application holds all services and dependencies
type Application struct {
	DB                       *gorm.DB
	UserService              *UserService
	EmailService             *EmailService
	EventBus                 interfaces.EventBus
	EmailVerificationHandler *EmailVerificationHandler
	RecipeService            *recipeService
	CategoryService          *categoryService
	CollectionService        *collectionService
	CloudinaryService        interfaces.CloudinaryService
	ImageService             *ImageService
	Server                   *http.Server
}

// GetUserService returns the user service
func (app *Application) GetUserService() interfaces.UserService {
	return app.UserService
}

func (app *Application) GetEmailService() interfaces.EmailService {
	return app.EmailService
}

func (app *Application) GetRecipeService() interfaces.RecipeService {
	return app.RecipeService
}

func (app *Application) GetCategoryService() interfaces.CategoryService {
	return app.CategoryService
}

func (app *Application) GetCollectionService() interfaces.CollectionService {
	return app.CollectionService
}

func (app *Application) GetImageService() interfaces.ImageService {
	return app.ImageService
}

// NewApplication creates a new Application instance
func NewApplication() (*Application, error) {
	// Initialize database
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPass, dbHost, dbPort, dbName)

	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto migrate schemas
	if err := database.AutoMigrate(&db.UserEntity{}, &db.CategoryEntity{}, &db.RecipeEntity{}, &db.CollectionEntity{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database schema: %w", err)
	}

	// Initialize repositories
	userRepo := db.NewUserRepository(database)
	recipeRepo := db.NewRecipeRepository(database)
	categoryRepo := db.NewCategoryRepository(database)
	collectionRepo := db.NewCollectionRepository(database)

	// Initialize Cloudinary service
	cloudinaryService, err := cloudinary.NewCloudinaryService()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Cloudinary service: %w", err)
	}

	// Initialize services
	emailService := NewEmailService()
	eventBus := NewEventBus()
	userService := NewUserService(userRepo, eventBus)
	emailVerificationHandler := NewEmailVerificationHandler(userRepo, emailService)

	recipeService := NewRecipeService(recipeRepo)
	categoryService := NewCategoryService(categoryRepo)
	collectionService := NewCollectionService(collectionRepo)
	imageService := NewImageService(cloudinaryService)

	// Subscribe to events
	eventBus.Subscribe("user.created", emailVerificationHandler)

	// Initialize application
	app := &Application{
		DB:                       database,
		UserService:              userService,
		EmailService:             emailService,
		EventBus:                 eventBus,
		EmailVerificationHandler: emailVerificationHandler,
		RecipeService:            recipeService,
		CategoryService:          categoryService,
		CollectionService:        collectionService,
		CloudinaryService:        cloudinaryService,
		ImageService:             imageService,
	}

	// Initialize HTTP server
	app.Server = http.NewServer(app)

	return app, nil
}

// Start starts the application
func (app *Application) Start() error {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s", port)
	return app.Server.Start(":" + port)
}
