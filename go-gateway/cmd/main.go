package main

import (
	"fmt"
	"sync"
	"time"

	ginrouter "github.com/Calmantara/go-common/infra/gin/router"
)

func main() {
	// call dependencies injection
	conf, http, err := BuildInRuntime()
	if err != nil {
		panic(err)
	}
	var wg sync.WaitGroup
	wg.Add(1)

	// go routine for http server
	go func() {
		http.SERVE(ginrouter.WithPort(fmt.Sprintf("%v", conf["httpport"])))
		wg.Done()
	}()
	fmt.Printf("runtime:%v\napplication-name:%s\napplication-port:%v",
		time.Now().Format(time.RFC850), conf["name"], conf["httpport"])
	wg.Wait()
}
