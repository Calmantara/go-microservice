package main

import (
	"fmt"
	"os"
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
	fmt.Printf(`-------------------------------------\n
	runtime:%v\napplication-name:%s\napplication-port:%v\napplication-env:%v\n
	--------------------------------------------------`,
		time.Now().Format(time.RFC850), conf["name"], conf["grpcport"], os.Getenv("ENV"))
	wg.Wait()
}
