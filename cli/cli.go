package cli

import (
	"flag"
	"fmt"

	"github.com/song940/qqwry.go/qqwry"
)

func Run() {
	var ip string
	flag.StringVar(&ip, "ip", "127.0.0.1", "IP address")
	flag.Parse()

	q, err := qqwry.NewQQwry("qqwry.dat")
	if err != nil {
		panic(err)
	}
	result, err := q.Find(ip)
	if err != nil {
		panic(err)
	}
	fmt.Println(result.IP)
	fmt.Println(result.Country, result.City)
}
