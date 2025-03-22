package http

import (
	"cookaholic/internal/infrastructure/middleware"
	"cookaholic/internal/interfaces"

	"github.com/gin-gonic/gin"
)

// Server holds all HTTP handlers and router configuration
type Server struct {
	router                  *gin.Engine
	app                     interfaces.Application
	userHandler             *UserHandler
	recipeHandler           *RecipeHandler
	categoryHandler         *CategoryHandler
	collectionHandler       *CollectionHandler
	recipeCollectionHandler *RecipeCollectionHandler
	recipeRatingHandler     *RecipeRatingHandler
	imageHandler            *ImageHandler
}

// NewServer creates a new Server instance
func NewServer(application interfaces.Application) *Server {
	router := gin.Default()

	server := &Server{
		router: router,
		app:    application,
	}

	// Initialize all handlers
	server.setupHandlers()

	return server
}

// setupHandlers initializes all HTTP handlers
func (s *Server) setupHandlers() {
	s.userHandler = NewUserHandler(s.router, s.app.GetUserService())
	s.recipeHandler = NewRecipeHandler(s.router, s.app.GetRecipeService())
	s.categoryHandler = NewCategoryHandler(s.router, s.app.GetCategoryService())
	s.collectionHandler = NewCollectionHandler(s.router, s.app.GetCollectionService())
	s.recipeCollectionHandler = NewRecipeCollectionHandler(s.app.GetRecipeCollectionService())
	s.recipeRatingHandler = NewRecipeRatingHandler(s.app.GetRecipeRatingService())
	s.imageHandler = NewImageHandler(s.router, s.app.GetImageService())

	// Public routes
	s.router.POST("/api/users/login", s.userHandler.Login)
	s.router.POST("/api/users/register", s.userHandler.Create)

	// Protected routes
	protected := s.router.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		users := protected.Group("/users")
		{
			users.POST("/email-verify", s.userHandler.VerifyOTP)
			users.POST("/resend-otp", s.userHandler.ResendOTP)
			users.GET("/:id", s.userHandler.GetByID)
			users.PUT("/:id", s.userHandler.Update)
			users.DELETE("/:id", s.userHandler.Delete)
			users.GET("", s.userHandler.List)
		}

		recipes := protected.Group("/recipes")
		{
			recipes.POST("", s.recipeHandler.CreateRecipe)
			recipes.GET("/:id", s.recipeHandler.GetRecipe)
			recipes.PUT("/:id", s.recipeHandler.UpdateRecipe)
			recipes.DELETE("/:id", s.recipeHandler.DeleteRecipe)
			recipes.GET("", s.recipeHandler.FilterRecipes)
			recipes.GET("/:id/collections", s.recipeCollectionHandler.GetCollectionsByRecipeID)
			recipes.GET("/:id/collections/:collectionId/check", s.recipeCollectionHandler.IsRecipeInCollection)
			recipes.GET("/:id/ratings", s.recipeRatingHandler.GetRatingsByRecipeID)
			recipes.GET("/:id/ratings/me", s.recipeRatingHandler.GetUserRatingForRecipe)
			recipes.POST("/:id/ratings", s.recipeRatingHandler.RateRecipe)
		}

		categories := protected.Group("/categories")
		{
			categories.POST("", s.categoryHandler.CreateCategory)
			categories.GET("/:id", s.categoryHandler.GetCategory)
			categories.PUT("/:id", s.categoryHandler.UpdateCategory)
			categories.DELETE("/:id", s.categoryHandler.DeleteCategory)
			categories.GET("", s.categoryHandler.ListCategories)
		}

		collections := protected.Group("/collections")
		{
			collections.POST("", s.collectionHandler.CreateCollection)
			collections.GET("/:id", s.collectionHandler.GetCollection)
			collections.PUT("/:id", s.collectionHandler.UpdateCollection)
			collections.DELETE("/:id", s.collectionHandler.DeleteCollection)
			collections.GET("", s.collectionHandler.GetUserCollections)
			collections.POST("/:id/recipes/:recipeId", s.recipeCollectionHandler.SaveRecipeToCollection)
			collections.DELETE("/:id/recipes/:recipeId", s.recipeCollectionHandler.RemoveRecipeFromCollection)
			collections.GET("/:id/recipes", s.recipeCollectionHandler.GetRecipesByCollectionID)
		}

		images := protected.Group("/images")
		{
			images.POST("/upload", s.imageHandler.UploadImage)
			images.POST("/upload-multiple", s.imageHandler.UploadMultipleImages)
		}

		ratings := protected.Group("/ratings")
		{
			ratings.PUT("/:id", s.recipeRatingHandler.UpdateRating)
			ratings.DELETE("/:id", s.recipeRatingHandler.DeleteRating)
		}
	}
}

// Start starts the HTTP server
func (s *Server) Start(addr string) error {
	return s.router.Run(addr)
}
