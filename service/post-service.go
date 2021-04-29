package service

import (
	"errors"

	"crashGo/entity"
	"crashGo/repository"
)

type PostService interface {
	Validate(post *entity.Post) error
	Create(post *entity.Post) (*entity.Post, error)
	FindAll()([]entity.Post, error)
	FindById(id string)(*entity.Post, error)
}

type service struct {
	repo repository.PostRepository
}

func NewPostService(repo repository.PostRepository) PostService {
	return &service{
		repo: repo,
	}
}

func (s *service) Validate(post *entity.Post) error {
	if post == nil {
		return errors.New("the post is empty")
	}
	if post.Title == "" {
		return errors.New("the title of the post is empty")
	}
	if post.Id != "" {
		post.Id = ""
	}
	return nil
}

func (s *service) Create(post *entity.Post) (*entity.Post, error) {
	return s.repo.Save(post)
}

func (s *service) FindAll() ([]entity.Post, error) {
	return s.repo.FindAll()
}

func (s *service) FindById(id string) (*entity.Post, error) {
	return s.repo.FindById(id)
}
