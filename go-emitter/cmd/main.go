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

// package main

// import (
// 	"context"
// 	"log"
// 	"os"
// 	"os/signal"
// 	"syscall"
// 	"time"

// 	"github.com/lovoo/goka"
// 	"github.com/lovoo/goka/codec"
// )

// var (
// 	brokers             = []string{"localhost:9092"}
// 	topic   goka.Stream = "example-transaction"
// 	group   goka.Group  = "balance-group"
// )

// // process messages until ctrl-c is pressed
// func runProcessor() {
// 	// process callback is invoked for each message delivered from
// 	// "example-stream" topic.
// 	cb := func(ctx goka.Context, msg interface{}) {
// 		var counter int64
// 		// ctx.Value() gets from the group table the value that is stored for
// 		// the message's key.
// 		if val := ctx.Value(); val != nil {
// 			counter = val.(int64)
// 		}
// 		counter++
// 		// SetValue stores the incremented counter in the group table for in
// 		// the message's key.
// 		ctx.SetValue(counter)
// 		log.Printf("key = %s, counter = %v, msg = %v", ctx.Key(), counter, msg)
// 	}

// 	// Define a new processor group. The group defines all inputs, outputs, and
// 	// serialization formats. The group-table topic is "example-group-table".
// 	g := goka.DefineGroup(group,
// 		goka.Input(topic, new(codec.String), cb),
// 		goka.Persist(new(codec.Int64)))

// 	p, err := goka.NewProcessor(brokers, g)
// 	if err != nil {
// 		log.Fatalf("error creating processor: %v", err)
// 	}
// 	ctx, cancel := context.WithCancel(context.Background())
// 	done := make(chan bool)
// 	go func() {
// 		defer close(done)
// 		if err = p.Run(ctx); err != nil {
// 			log.Fatalf("error running processor: %v", err)
// 		} else {

// 			log.Printf("Processor shutdown cleanly")
// 		}
// 	}()

// 	wait := make(chan os.Signal, 1)
// 	signal.Notify(wait, syscall.SIGINT, syscall.SIGTERM)
// 	<-wait   // wait for SIGINT/SIGTERM
// 	cancel() // gracefully stop processor
// 	<-done
// }

// func main() {
// 	go runEmitter() // emits one message and stops
// 	runProcessor()  // press ctrl-c to stop
// }
