package main

import (
	"encoding/json"
	"fmt"
	"go-cache-restapi/cache"
	"go-cache-restapi/helper"
	"go-cache-restapi/middleware"
	"go-cache-restapi/model"
	"go-cache-restapi/service"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

var c *cache.Cache
var r *mux.Router

func init() {
	fmt.Println("** init çalıştı")
	initCache()
	initLoadCacheFromFile()
	initSaveFile()
	initMuxRouter()
	//I use Goroutine. Bcs goroutine works better performance than thread
	go SaveDataInterval()
}

func main() {
	r.HandleFunc("/get/{key}", Get).Methods(http.MethodGet)
	r.HandleFunc("/set", Set).Methods(http.MethodPost)
	r.HandleFunc("/flush", Flush).Methods(http.MethodGet)
	r.HandleFunc("/get", GetAll).Methods(http.MethodGet)
	r.HandleFunc("/delete/{key}", Delete).Methods(http.MethodDelete)
	r.HandleFunc("/", Index).Methods(http.MethodGet)

	httpLogServer()

	http.ListenAndServe(":8000", r)

}

func Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("Index Page"))
}

func Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	param := mux.Vars(r)
	key := param["key"]
	_, resp := service.Get(c, key, r.Method)
	response := helper.JsonMarshall(resp)
	w.Write(response)
	initSaveFile()
}
func GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, resp := service.GetAll(c, r.Method)
	response := helper.JsonMarshall(resp)
	w.Write(response)
	initSaveFile()
}

//If exist item key, this item will update
func Set(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var keyVal model.KeyValue
	err := json.NewDecoder(r.Body).Decode(&keyVal)
	if err == nil {
		fmt.Println("/main/Set çalıştı", keyVal)
		resp := service.Set(keyVal, c, r.Method)
		response := helper.JsonMarshall(resp)
		w.Write(response)
		initSaveFile()
	}
}

func Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	param := mux.Vars(r)
	key := param["key"]
	resp := service.Delete(c, key, r.Method)
	response := helper.JsonMarshall(resp)
	w.Write(response)
	initSaveFile()
}

func Flush(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	resp := service.Flush(c, r.Method)
	response := helper.JsonMarshall(resp)
	w.Write(response)
	initSaveFile()
}

type Writer interface {
	Write(p []byte) (n int, err error)
}
type Reader interface {
	Read(p []byte) (n int, err error)
}

func initCache() {
	c = cache.New(5*time.Minute, 10*time.Minute)
}

//Load the cache from saved file when you can run the project
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

//Save the cache and cache file periodically
//I prefer it to run every 2 minutes
func SaveDataInterval() {
	for range time.Tick(time.Minute * 2) {
		initSaveFile()
		service.GetAll(c, "GET")
		fmt.Println("File automatic saved.")
	}
	time.Sleep(time.Second * 5)
}

func httpLogServer() {
	logger := log.New(os.Stdout, "", log.LstdFlags)
	logMiddleware := middleware.NewLogMiddleware(logger)
	r.Use(logMiddleware.Func())
}
