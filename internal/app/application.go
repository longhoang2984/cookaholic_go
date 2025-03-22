package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"cookaholic/internal/infrastructure/cloudinary"
	"cookaholic/internal/infrastructure/db"
	"cookaholic/internal/infrastructure/http"
	"cookaholic/internal/interfaces"

	"github.com/google/uuid"
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
	RecipeCollectionService  interfaces.RecipeCollectionService
	RecipeRatingService      interfaces.RecipeRatingService
	CloudinaryService        interfaces.CloudinaryService
	ImageService             *ImageService
	Server                   *http.Server
	stopRatingCron           chan bool
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

func (app *Application) GetRecipeCollectionService() interfaces.RecipeCollectionService {
	return app.RecipeCollectionService
}

func (app *Application) GetRecipeRatingService() interfaces.RecipeRatingService {
	return app.RecipeRatingService
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
	if err := database.AutoMigrate(&db.UserEntity{}, &db.CategoryEntity{}, &db.RecipeEntity{}, &db.CollectionEntity{}, &db.RecipeCollectionEntity{}, &db.RecipeRatingEntity{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database schema: %w", err)
	}

	// Initialize repositories
	userRepo := db.NewUserRepository(database)
	recipeRepo := db.NewRecipeRepository(database)
	categoryRepo := db.NewCategoryRepository(database)
	collectionRepo := db.NewCollectionRepository(database)
	recipeCollectionRepo := db.NewRecipeCollectionRepository(database)
	recipeRatingRepo := db.NewRecipeRatingRepository(database)

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
	recipeCollectionService := NewRecipeCollectionService(recipeCollectionRepo, recipeRepo, collectionRepo)
	recipeRatingService := NewRecipeRatingService(recipeRatingRepo, recipeRepo)
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
		RecipeCollectionService:  recipeCollectionService,
		RecipeRatingService:      recipeRatingService,
		CloudinaryService:        cloudinaryService,
		ImageService:             imageService,
		stopRatingCron:           make(chan bool),
	}

	// Initialize HTTP server
	app.Server = http.NewServer(app)

	// Start the rating update cron job
	go app.startRatingUpdateCron()

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

// Stop stops the application
func (app *Application) Stop() {
	log.Println("Stopping application...")
	// Stop the rating update cron job
	app.stopRatingCron <- true
}

// startRatingUpdateCron starts a goroutine that periodically updates recipe ratings
func (app *Application) startRatingUpdateCron() {
	ticker := time.NewTicker(6 * time.Hour)
	defer ticker.Stop()

	log.Println("Starting recipe rating update cron job...")

	for {
		select {
		case <-ticker.C:
			log.Println("Running recipe rating update job...")
			app.updateAllRecipeRatings()
		case <-app.stopRatingCron:
			log.Println("Stopping recipe rating update cron job...")
			return
		}
	}
}

// updateAllRecipeRatings updates all recipe ratings
func (app *Application) updateAllRecipeRatings() {
	// Get all recipes
	var recipeIDs []string

	if err := app.DB.Table("recipes").Select("id").Where("status = ?", 1).Pluck("id", &recipeIDs).Error; err != nil {
		log.Printf("Error fetching recipe IDs: %v", err)
		return
	}

	// Update ratings for each recipe
	for _, idStr := range recipeIDs {
		// Get ID as UUID
		id, err := uuid.Parse(idStr)
		if err != nil {
			log.Printf("Error parsing UUID %s: %v", idStr, err)
			continue
		}

		// Create a context for the operation
		ctx := context.Background()

		// Get the recipe rating repository
		recipeRatingRepo := db.NewRecipeRatingRepository(app.DB)

		// Update recipe rating summary
		if err := recipeRatingRepo.UpdateRecipeRatingSummary(ctx, id); err != nil {
			log.Printf("Error updating rating for recipe %s: %v", idStr, err)
		}
	}
}
