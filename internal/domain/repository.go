package domain

type Repository interface {
	Save(model *Model) error
}
