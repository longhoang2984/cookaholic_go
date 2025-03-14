package interfaces

// Application defines the interface for the application
type Application interface {
	GetUserService() UserService
	GetRecipeService() RecipeService
	GetCategoryService() CategoryService
}
