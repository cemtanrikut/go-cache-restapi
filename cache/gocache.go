package cache

import (
	"encoding/gob"
	"fmt"
	"go-cache-restapi/model"
	"io"
	"os"
	"sync"
	"time"
)

func init() {
	gob.Register(model.Value{})
}

//Structs
type Item struct {
	Object     interface{}
	Expiration int64
}
type Cache struct {
	*cache
}
type cache struct {
	defaultExpiration time.Duration
	items             map[string]Item
	mu                sync.RWMutex
	onEvicted         func(string, interface{})
}

func (item Item) Expired() bool {
	if item.Expiration == 0 {
		return false
	}
	//If item has expired return true
	return time.Now().UnixNano() > item.Expiration
}

//consts
const (
	NoExpiration      time.Duration = -1
	DefaultExpiration time.Duration = 0
)

//Adds item to cache
//If exist item key, this item will update
func (c *cache) Set(k string, i interface{}, td time.Duration) error {
	var e int64
	if td == DefaultExpiration {
		td = c.defaultExpiration
	}
	if td > 0 {
		e = time.Now().Add(td).UnixNano()
	}
	c.mu.Lock()
	c.items[k] = Item{
		Object:     i,
		Expiration: e,
	}
	c.mu.Unlock()
	return nil
}
func (c *cache) set(k string, i interface{}, td time.Duration) {
	var e int64
	if td == DefaultExpiration {
		td = c.defaultExpiration
	}
	if td > 0 {
		e = time.Now().Add(td).UnixNano()
	}
	c.items[k] = Item{
		Object:     i,
		Expiration: e,
	}
}

//Gets the item where has been to cache
func (c *cache) Get(k string) (interface{}, bool) {
	c.mu.RLock()
	// "Inlining" of get and Expired
	item, found := c.items[k]
	if !found {
		c.mu.RUnlock()
		return nil, false
	}
	if item.Expiration > 0 {
		if time.Now().UnixNano() > item.Expiration {
			c.mu.RUnlock()
			return nil, false
		}
	}
	c.mu.RUnlock()
	return item.Object, true
}

//Deletes the item where has been to cache
func (c *cache) Delete(k string) error {
	c.mu.Lock()
	v, evicted := c.delete(k)
	c.mu.Unlock()
	if evicted {
		c.onEvicted(k, v)
	}
	return nil
}
func (c *cache) delete(k string) (interface{}, bool) {
	if c.onEvicted != nil {
		if v, found := c.items[k]; found {
			delete(c.items, k)
			return v.Object, true
		}
	}
	delete(c.items, k)
	return nil, false
}

//Saving cache with using Gob
func (c *cache) Save(w io.Writer) (err error) {
	enc := gob.NewEncoder(w)
	defer func() {
		if x := recover(); x != nil {
			err = fmt.Errorf("Error register to Gob Lib")
			fmt.Println(err)
		}
	}()
	c.mu.RLock()
	defer c.mu.RUnlock()
	for _, v := range c.items {
		gob.Register(v.Object)
	}
	err = enc.Encode(&c.items)
	fmt.Println("/cache/save : ", err)
	return
}

//Save file and run Save() func
func (c *cache) SaveFile(fileName string) error {
	fp, err := os.Create(fileName)
	if err != nil {
		fmt.Println("/cache/SaveFile err:", err)
		return err
	}
	err = c.Save(fp)
	if err != nil {
		fp.Close()
		fmt.Println("/cache/Doesn' save file")
		return err
	}
	return fp.Close()
}

func (c *cache) Load(r io.Reader) error {
	dec := gob.NewDecoder(r)
	items := map[string]Item{}
	err := dec.Decode(&items)
	if err == nil {
		c.mu.Lock()
		defer c.mu.Unlock()
		for k, v := range items {
			ov, found := c.items[k]
			if !found || ov.Expired() {
				c.items[k] = v
			}
		}
	}
	fmt.Println("/cache/Load error: ", err)
	return err
}

//LoadFile and run Load() func
func (c *cache) LoadFile(fileName string) error {
	fp, err := os.Open(fileName)
	if err != nil {
		fmt.Println("/cache/LoadFile open err")
		return err
	}
	err = c.Load(fp)
	if err != nil {
		fp.Close()
		fmt.Println("/cache/LoadFile load err", err)
		return err
	}
	return fp.Close()
}

//Gets all items in cache
func (c *cache) Items() map[string]Item {
	c.mu.RLock()
	defer c.mu.RUnlock()
	m := make(map[string]Item, len(c.items))
	currentTime := time.Now().UnixNano()
	for c, t := range c.items {
		if t.Expiration > 0 {
			if currentTime > t.Expiration {
				continue
			}
		}
		m[c] = t
	}
	return m
}

// Flush the cache
func (c *cache) Flush() {
	c.mu.Lock()
	c.items = map[string]Item{}
	c.mu.Unlock()
}

//Creates new cache
func New(defaultExpiration, cleanupInterval time.Duration) *Cache {
	items := make(map[string]Item)
	return newCache(defaultExpiration, items)
}
func newCache(de time.Duration, m map[string]Item) *Cache {
	if de == 0 {
		de = -1
	}
	c := &cache{
		defaultExpiration: de,
		items:             m,
	}
	Cch := Cache{c}

	return &Cch
}
