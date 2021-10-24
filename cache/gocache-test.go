package cache

import (
	"testing"
)

type TestStruct struct {
	Num      int
	Children []*TestStruct
}

var cacheTest *Cache

func Init() {
	cacheTest = New(DefaultExpiration, 0)
}

func TestCacheFill(t *testing.T) {
	testFill(t, cacheTest)
}

func testFill(t *testing.T, tc *Cache) {
	cacheTest.Set("data1", "data1", DefaultExpiration)
	cacheTest.Set("data2", "data2", DefaultExpiration)
	cacheTest.Set("data3", "data3", DefaultExpiration)

	oc := New(DefaultExpiration, 0)
	first, found := oc.Get("data1")
	if !found {
		t.Error("a not found")
	}
	if first.(string) != "data1" {
		t.Error("data1 not equal data1")
	}

	second, found := oc.Get("data2")
	if !found {
		t.Error("data2 not found")
	}
	if second.(string) != "b" {
		t.Error("data2 not equal data2")
	}

	third, found := oc.Get("data3")
	if !found {
		t.Error("data3 not found")
	}
	if third.(string) != "data3" {
		t.Error("data3 not equal data3")
	}
}

func TestGet(t *testing.T) {

	a, found := cacheTest.Get("a")
	if found || a != nil {
		t.Error("Get-Error: a shouldn't exist but found", a)
	}

	b, found := cacheTest.Get("b")
	if found || b != nil {
		t.Error("Get-Error: b shouldn't exist but found", b)
	}

	c, found := cacheTest.Get("c")
	if found || c != nil {
		t.Error("Get-Error: c shouldn't exist but found", c)
	}

	cacheTest.Set("a", 1, DefaultExpiration)
	_, found = cacheTest.Get("a")
	if !found {
		t.Error("Get-Error: not found added item")
	}

	cacheTest.Set("b", "b", DefaultExpiration)
	_, found = cacheTest.Get("b")
	if !found {
		t.Error("Get-Error: not found added item")
	}
}

func TestSet(t *testing.T) {

	err := cacheTest.Set("data1", "data2", DefaultExpiration)
	if err != nil {
		t.Error("Set-Error: not set")
	}
}

func TestDelete(t *testing.T) {

	cacheTest.Set("data1", "data2", DefaultExpiration)
	cacheTest.Delete("data1")
	_, found := cacheTest.Get("data1")
	if found {
		t.Error("Delete-Error: item should be deleted")
	}
}

func TestFlush(t *testing.T) {

	cacheTest.Set("data1", "data2", DefaultExpiration)
	cacheTest.Flush()
	_, found := cacheTest.Get("data1")
	if found {
		t.Error("Flush-Error: data1 found, but shouldn't")
	}
}
