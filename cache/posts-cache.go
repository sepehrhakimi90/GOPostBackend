package cache

import (
	"crashGo/entity"
)

type PostCache interface {
	Set(key string, value *entity.Post)
	Get(key string) *entity.Post
}
