package iredis

import (
	"log"
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	if v, err := Get("test"); err == nil && v != "" {
		t.Fail()
	}
}

func TestSet(t *testing.T) {
	if err := Set("test2", "123456"); err != nil {
		log.Println(err)
		t.Fail()
		return
	}
	if v, err := Get("test2"); err == nil && v != "123456" {
		t.Fail()
	}
}

func TestDel(t *testing.T) {
	TestSet(t)
	if err := Del("test2"); err != nil {
		log.Println(err)
		t.Fail()
		return
	}
	if v, err := Get("test2"); err == nil && v != "" {
		t.Fail()
	}
}

func TestSetEX(t *testing.T) {
	SetEx("aaa", "222", 3*time.Second)
	log.Println(client.Get("aaa"))
	time.Sleep(4 * time.Second)
	log.Println(client.Get("aaa"))
}
