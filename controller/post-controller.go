package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/sepehrhakimi90/GOPostBackend/cache"
	"github.com/sepehrhakimi90/GOPostBackend/entity"
	"github.com/sepehrhakimi90/GOPostBackend/errors"
	"github.com/sepehrhakimi90/GOPostBackend/service"
)

type PostController interface {
	GetPosts(writer http.ResponseWriter, request *http.Request)
	AddPost(writer http.ResponseWriter, request *http.Request)
	GetPostByID(writer http.ResponseWriter, request *http.Request)
}

type controller struct {
	postService service.PostService
	postCache   cache.PostCache
}

func NewPostController(service service.PostService, postCache cache.PostCache) PostController {
	return &controller{
		postService: service,
		postCache: postCache,
	}
}

func (c *controller) GetPostByID(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-type", "application/json")
	id := strings.Split(request.URL.Path, "/")[2]
	var err error
	post := c.postCache.Get(id)
	if post == nil {
		post, err = c.postService.FindById(id)
		if err != nil {
			if err.Error() == "mongo: no documents in result" {
				writer.WriteHeader(http.StatusNotFound)
				json.NewEncoder(writer).Encode(errors.ServiceError{Message: "post does not exist"})
				return
			}
			writer.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(writer).Encode(errors.ServiceError{Message: "error in retrieving post"})
			return
		}
		c.postCache.Set(id, post)
		log.Println("Controller: data was not in cache")
	} else {
		log.Println("Controller: data was in cache")
	}

	encoder := json.NewEncoder(writer)
	err = encoder.Encode(post)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(writer).Encode(errors.ServiceError{Message: "error in marshalling data"})
		return
	}
	writer.WriteHeader(http.StatusOK)
}

func (c *controller) GetPosts(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-type", "application/json")
	encoder := json.NewEncoder(writer)
	posts, err := c.postService.FindAll()
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(writer).Encode(errors.ServiceError{Message: "error in retrieving posts"})
		return
	}
	err = encoder.Encode(posts)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(writer).Encode(errors.ServiceError{Message: "error in marshalling data"})
		return
	}
	writer.WriteHeader(http.StatusOK)
}

func (c *controller) AddPost(writer http.ResponseWriter, request *http.Request) {
	post := &entity.Post{}
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(post)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode(errors.ServiceError{Message: "error in post data"})
		return
	}

	err = c.postService.Validate(post)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode(errors.ServiceError{Message: err.Error()})
		return
	}

	newPost, err := c.postService.Create(post)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(writer).Encode(errors.ServiceError{Message: "error in saving post"})
		return
	}
	encoder := json.NewEncoder(writer)
	err = encoder.Encode(newPost)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(writer).Encode(errors.ServiceError{Message: "error in marshalling data"})
		return
	}

	writer.WriteHeader(http.StatusOK)
}
