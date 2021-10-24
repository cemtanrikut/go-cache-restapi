package service

import (
	"go-cache-restapi/cache"
	"go-cache-restapi/controller"
	"go-cache-restapi/helper"
	"go-cache-restapi/model"
)

func Get(c *cache.Cache, key string, method string) ([]byte, *helper.Resp) {
	return controller.Get(c, key, method)

}

func GetAll(c *cache.Cache, method string) (map[string]cache.Item, *helper.Resp) {
	return controller.GetAll(c, method)
}

func Set(keyVal model.KeyValue, c *cache.Cache, method string) *helper.Resp {
	return controller.Set(keyVal, c, method)
}

func Delete(c *cache.Cache, key string, method string) *helper.Resp {
	return controller.DeleteItem(c, key, method)
}

func Flush(c *cache.Cache, method string) *helper.Resp {
	return controller.Flush(c, method)
}

func SaveCacheFile(c *cache.Cache, fileName string) {
	controller.SaveCacheFile(c, fileName)
}

func LoadCacheFile(c *cache.Cache, fileName string) {
	controller.LoadCacheFile(c, fileName)
}
