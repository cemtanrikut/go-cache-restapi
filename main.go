package main

import (
	"fmt"
	"go-cache-restapi/cache"
	"go-cache-restapi/service"
	"os"
	"time"

	"github.com/gorilla/mux"
)

var c *cache.Cache
var r *mux.Router

func main() {

}

func initCache() {
	c = cache.New(5*time.Minute, 10*time.Minute)
}
func initLoadCacheFromFile() {
	fileName := "TIMESTAMP-data.gob"
	_, err := os.Open("tmp/" + fileName)
	if err == nil {
		//Load cache same time
		service.LoadCacheFile(c, "tmp/"+fileName)
		fmt.Println("--LoadCacheFile çalıştı")
	} else {
		fmt.Println("init load hatası ", err)
	}
}
func initSaveFile() {
	fileName := "TIMESTAMP-data.gob"
	r, err := os.Create("tmp/" + fileName)
	fmt.Println("init save data: ", r)
	fmt.Println("init save err: ", err)
	if err == nil {
		//Save cash same time
		service.SaveCacheFile(c, "tmp/"+fileName)
		fmt.Println("-SaveCacheFile çalıştı")
	} else {
		fmt.Println("init save hatası", err)
	}
}
func initMuxRouter() {
	r = mux.NewRouter()
}
