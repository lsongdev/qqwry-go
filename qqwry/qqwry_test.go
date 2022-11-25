package qqwry

import (
	"log"
	"testing"
)

func TestQQwry(t *testing.T) {
	q := NewQQwry("qqwry.dat")
	result, err := q.Find("8.8.8.8")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("ip:%v, country:%v, city:%v", result.IP, result.Country, result.City)
}
