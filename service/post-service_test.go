package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/sepehrhakimi90/GOPostBackend/entity"
)

type mockRepository struct {
	mock.Mock
}

func (m mockRepository) Save(post *entity.Post) (*entity.Post, error) {
	args := m.Called()
	result := args.Get(0)
	return result.(*entity.Post), args.Error(1)
}

func (m mockRepository) FindAll() ([]entity.Post, error) {
	args := m.Called()
	return args.Get(0).([]entity.Post), args.Error(1)
}

func TestCreat(t *testing.T) {
	mockRepo := new(mockRepository)
	post := entity.Post{
		Id:    "1",
		Title: "A",
		Text:  "B",
	}

	mockRepo.On("Save").Return(&post, nil)

	testService := NewPostService(mockRepo)

	result, err := testService.Create(&post)

	mockRepo.AssertExpectations(t)

	assert.Nil(t, err)
	assert.Equal(t, post.Title, result.Title)
	assert.Equal(t, post.Text, result.Text)
}

func TestFindAll(t *testing.T) {
	mockRepo := new(mockRepository)

	post := entity.Post{
		Id:    "1",
		Title: "A",
		Text:  "B",
	}

	mockRepo.On("FindAll").Return([]entity.Post{post}, nil)

	testService := NewPostService(mockRepo)

	result, err := testService.FindAll()

	mockRepo.AssertExpectations(t)

	assert.Nil(t, err)
	assert.Equal(t, post.Title, result[0].Title)
	assert.Equal(t, post.Text, result[0].Text)
}

func TestValidateEmptyPost(t *testing.T) {
	testService := NewPostService(nil)

	err := testService.Validate(nil)

	assert.NotNil(t, err)
	assert.Equal(t, "the post is empty", err.Error())
}

func TestValidateEmptyPostTitle(t *testing.T) {
	post := entity.Post{
		Id:    "1",
		Title: "",
		Text:  "B",
	}

	testService := NewPostService(nil)

	err := testService.Validate(&post)

	assert.NotNil(t, err)
	assert.Equal(t, "the title of the post is empty", err.Error())
}