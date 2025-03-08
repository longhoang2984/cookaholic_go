# Cookaholic

A recipe sharing and management application built with Go, following Hexagonal Architecture principles. Users can create accounts, share recipes, and discover new dishes from other users.

## Features

- User authentication and authorization
- User profile management (create, update, delete)
- Recipe management (create, read, update, delete recipes)
- Recipe filtering and search
- More features coming soon!

## Project Structure
```
internal/
  ├── domain/              # Core business models and rules
  │   └── user.go         # User domain model
  ├── interfaces/          # Interface definitions
  │   ├── user_service.go # User service interface
  │   ├── user_repository.go # User repository interface
  │   └── application.go  # Application interface
  ├── app/                # Application implementation
  │   ├── user_service.go # User service implementation
  │   └── application.go  # Application implementation
  └── infrastructure/     # External implementations
      ├── http/          # HTTP handlers and routing
      ├── middleware/    # HTTP middleware
      └── db/           # Database implementations
```

## How to Create a New Service

To create a new service in this project, follow these steps:

1. Create a new domain model in the `internal/domain` directory
2. Create a new repository interface in `internal/interfaces`
3. Create a new service interface in `internal/interfaces`
4. Implement the service in `internal/app`
5. Create a new handler in `internal/infrastructure/http`
6. Add the new routes in `internal/infrastructure/http/router.go`

Example structure for a new feature:

```
internal/
  ├── domain/
  │   └── your_model.go
  ├── interfaces/
  │   ├── your_repository.go
  │   └── your_service.go
  ├── app/
  │   └── your_service.go
  └── infrastructure/
      └── http/
          └── your_handler.go
```

Follow the existing patterns in the codebase to maintain consistency with the hexagonal architecture principles.

### Example: Creating a Recipe Service

Here's a practical example of creating a Recipe service:

1. First, create the domain model in `internal/domain/model/recipe.go`:
```go
package model

type Recipe struct {
    ID          string
    Title       string
    Description string
    Ingredients []string
    Steps       []string
    UserID      string
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
```

2. Define the repository interface in `internal/interfaces/recipe.go`:
```go
package interfaces

type RecipeRepository interface {
    Create(ctx context.Context, recipe *model.Recipe) error
    GetByID(ctx context.Context, id string) (*model.Recipe, error)
    ListByUserID(ctx context.Context, userID string) ([]*model.Recipe, error)
    Update(ctx context.Context, recipe *model.Recipe) error
    Delete(ctx context.Context, id string) error
}
```

3. Define the service interface in `internal/interfaces/recipe.go`:
```go
package interfaces

type RecipeService interface {
    CreateRecipe(ctx context.Context, recipe *model.Recipe) error
    GetRecipe(ctx context.Context, id string) (*model.Recipe, error)
    ListUserRecipes(ctx context.Context, userID string) ([]*model.Recipe, error)
    UpdateRecipe(ctx context.Context, recipe *model.Recipe) error
    DeleteRecipe(ctx context.Context, id string) error
}
```

4. Implement the service in `internal/app/recipe.go`:
```go
package app

type recipeService struct {
    recipeRepo interfaces.RecipeRepository
}

func NewRecipeService(recipeRepo interfaces.RecipeRepository) RecipeService {
    return &recipeService{
        recipeRepo: recipeRepo,
    }
}

func (s *recipeService) CreateRecipe(ctx context.Context, input interfaces.CreateRecipeInput) error {
    // Add business logic here
    return s.recipeRepo.Create(ctx, recipe)
}

// Implement other methods...
```

5. Create the handler in `internal/infrastructure/http/recipe.go`:
```go
package infrastructure

type RecipeHandler struct {
    recipeService RecipeService
}

func NewRecipeHandler(recipeService RecipeService) *RecipeHandler {
    return &RecipeHandler{
        recipeService: recipeService,
    }
}

func (h *RecipeHandler) CreateRecipe(c *gin.Context) {
    var recipe model.Recipe
    if err := c.ShouldBindJSON(&recipe); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := h.recipeService.CreateRecipe(c.Request.Context(), &recipe); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, recipe)
}

// Implement other handlers...
```

6. Add routes in `internal/infrastructure/http/router.go`:
```go
func SetupRouter(handlers ...interface{}) *gin.Engine {
    router := gin.Default()
    
    // Add recipe routes
    recipeHandler := handlers[0].(*RecipeHandler)
    recipes := router.Group("/api/recipes")
    {
        recipes.POST("/", recipeHandler.CreateRecipe)
        recipes.GET("/:id", recipeHandler.GetRecipe)
        recipes.GET("/user/:userId", recipeHandler.ListUserRecipes)
        recipes.PUT("/:id", recipeHandler.UpdateRecipe)
        recipes.DELETE("/:id", recipeHandler.DeleteRecipe)
    }
    
    return router
}
```

This example demonstrates how to implement a complete feature following the hexagonal architecture pattern. Each layer has its specific responsibility:
- Domain layer: Contains business logic and interfaces
- Repository layer: Handles data persistence
- Service layer: Implements business logic
- Handler layer: Manages HTTP requests and responses
- Router layer: Defines API endpoints

