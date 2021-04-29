package cache

import (
	"github.com/sepehrhakimi90/GOPostBackend/entity"
)

type PostCache interface {
	Set(key string, value *entity.Post)
	Get(key string) *entity.Post
}
