package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"github.com/sepehrhakimi90/GOPostBackend/cache"
	"github.com/sepehrhakimi90/GOPostBackend/controller"
	"github.com/sepehrhakimi90/GOPostBackend/repository"
	"github.com/sepehrhakimi90/GOPostBackend/service"
)

const (
	port  = ":8000"
	)

func main() {

	postRepo, err := repository.NewMongoPostRepo()
	if err != nil {
		log.Fatal(err)
	}



	cachePost, err := getRedisCache()
	if err != nil {
		panic(err)
	}
	postController := controller.NewPostController(service.NewPostService(postRepo), cachePost)

	router := mux.NewRouter()

	router.HandleFunc("/posts", postController.GetPosts).Methods(http.MethodGet)
	router.HandleFunc("/add", postController.AddPost).Methods(http.MethodPost)
	router.HandleFunc("/posts/{id}", postController.GetPostByID).Methods(http.MethodGet)
	router.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintln(writer, "Up and Running...")
	})

	repository.NewMongoPostRepo()

	log.Println("Server is listening on port", port)

	log.Fatalln(http.ListenAndServe(port, router))
}

func getRedisCache() (cache.PostCache, error){
	redisHost := os.Getenv("REDIS_IP")
	redisPORT := os.Getenv("REDIS_PORT")
	redisEx := os.Getenv("REDIS_EXPIRY")
	if redisPORT == "" {
		redisPORT = "6379"
	}
	if redisEx == "" {
		redisEx = "10000"
	}
	redisExpiry, err := strconv.Atoi(redisEx)
	if err != nil {
		return nil, err
	}

	return cache.NewRedisCache(fmt.Sprintf("%s:%s", redisHost, redisPORT), 10, time.Duration(redisExpiry) * time.Millisecond), nil
}
