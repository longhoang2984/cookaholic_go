package http

import (
	"cookaholic/internal/infrastructure/middleware"
	"cookaholic/internal/interfaces"

	"github.com/gin-gonic/gin"
)

// Server holds all HTTP handlers and router configuration
type Server struct {
	router        *gin.Engine
	app           interfaces.Application
	userHandler   *UserHandler
	recipeHandler *RecipeHandler
	// Add other handlers here
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

	// Public routes
	s.router.POST("/api/users/login", s.userHandler.Login)
	s.router.POST("/api/users/register", s.userHandler.Create)

	// Protected routes
	protected := s.router.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		users := protected.Group("/users")
		{
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
		}
	}
}

// Start starts the HTTP server
func (s *Server) Start(addr string) error {
	return s.router.Run(addr)
}
