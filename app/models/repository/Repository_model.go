package repositoryModels

import "gorm.io/gorm"

type RepositoryModel[T any] struct {
	Statue gorm.DB
	Result T
}
