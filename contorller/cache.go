package contorller

import (
	"fmt"
	"go-cache-restapi/cache"
	"go-cache-restapi/helper"
	"go-cache-restapi/model"
)

func Get(c *cache.Cache, key string, method string) ([]byte, *helper.Resp) {
	data, exist := c.Get(key)
	res := helper.NewResponse()
	if !exist {
		fmt.Println(key, " data not exist")
		res.Success = false
		res.Action = "Get"
		res.Data = data
		res.Errors = nil
		res.Method = method
		return nil, nil
	}

	resByte, err := data.([]byte)
	if err {
		res.Success = false
		res.Action = "Get"
		res.Data = data
		res.Errors = nil
		res.Method = method
		return nil, nil
	}

	res.Success = true
	res.Action = "Get"
	res.Data = data
	res.Errors = nil
	res.Method = method

	fmt.Println("get method: ", data)
	return resByte, res
}

//If exist item key, this item will update
func Set(data model.KeyValue, c *cache.Cache, method string) *helper.Resp {
	err := c.Set(data.Key, data.Value, cache.NoExpiration)
	res := helper.NewResponse()
	if err == nil {
		res.Success = true
		res.Action = "Set"
		res.Data = data
		res.Errors = nil
		res.Method = method
	} else {
		res.Success = false
		res.Action = "Set"
		res.Data = nil
		res.Errors = map[string]string{"Request Body": "Provided request body malformed"}
		res.Method = method
	}

	fmt.Println("set method: ", c)

	return res
}

func GetAll(c *cache.Cache, method string) (map[string]cache.Item, *helper.Resp) {
	items := c.Items()
	res := helper.NewResponse()
	fmt.Println("Items: ", c.Items())
	res.Success = true
	res.Action = "Get"
	res.Data = items
	res.Errors = nil
	res.Method = method

	return items, res
}

func DeleteItem(c *cache.Cache, key string, method string) *helper.Resp {
	err := c.Delete(key)
	res := helper.NewResponse()
	if err == nil {
		res.Success = true
		res.Action = "DELETE"
		res.Data = nil
		res.Errors = nil
		res.Method = method
	}
	return res
}

func Flush(c *cache.Cache, m string) *helper.Resp {
	c.Flush()
	res := helper.NewResponse()
	res.Success = true
	res.Action = "Get"
	res.Data = ""
	res.Errors = nil
	res.Method = m
	fmt.Println("flush")

	return res
}

func SaveCacheFile(c *cache.Cache, fileName string) {
	c.SaveFile(fileName)
	fmt.Println("/controller/savecachefile çalıştı")
}

func LoadCacheFile(c *cache.Cache, fileName string) {
	c.LoadFile(fileName)
	fmt.Println("/controller/LoadCacheFile çalıştı")
}
