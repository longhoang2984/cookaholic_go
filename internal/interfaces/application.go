package interfaces

// Application defines the interface for the application
type Application interface {
	GetUserService() UserService
	GetEmailService() EmailService
	GetRecipeService() RecipeService
	GetCategoryService() CategoryService
	GetCollectionService() CollectionService
}
