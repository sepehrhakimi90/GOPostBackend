package repository

import (
	"github.com/sepehrhakimi90/GOPostBackend/entity"
)

type PostRepository interface {
	Save(post *entity.Post) (*entity.Post, error)
	FindAll() ([]entity.Post, error)
	FindById(id string) (*entity.Post, error)
}
