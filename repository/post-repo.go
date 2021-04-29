package repository

import (
	"crashGo/entity"
)

type PostRepository interface {
	Save(post *entity.Post) (*entity.Post, error)
	FindAll() ([]entity.Post, error)
	FindById(id string) (*entity.Post, error)
}
