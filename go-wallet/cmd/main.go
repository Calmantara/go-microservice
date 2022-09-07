package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	// call dependencies injection
	conf, grpc, err := BuildInRuntime()
	if err != nil {
		panic(err)
	}
	var wg sync.WaitGroup
	wg.Add(1)

	// go routine for grpc server
	go func() {
		grpc.SERVE()
		wg.Done()
	}()
	fmt.Printf("runtime:%v\napplication-name:%s\napplication-port:%v",
		time.Now().Format(time.RFC850), conf["name"], conf["grpcport"])
	wg.Wait()
}
