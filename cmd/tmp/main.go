package main

import (
	"context"
	"fmt"
	"time"

	"github.com/daysleep666/short_chain/pkg"
)

func main() {
	s := pkg.NewUniqueIDSnowflakeService(1)
	do := func() {
		for i := 0; i < 10000; i++ {
			id, _ := s.Generate(context.Background())
			fmt.Println(id)
		}
	}
	go do()
	go do()
	time.Sleep(time.Second)
}
