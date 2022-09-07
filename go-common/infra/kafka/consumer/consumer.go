package consumer

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Calmantara/go-common/logger"
	"github.com/Calmantara/go-common/setup/config"
	"github.com/Calmantara/go-common/topic"
	"github.com/lovoo/goka"
	"github.com/lovoo/goka/codec"
)

type Consumer interface {
	Consume()
}

type KafkaConsumerImpl struct {
	sugar   logger.CustomLogger
	brokers []string
	workers []KafkaWorker
}

type KafkaWorker struct {
	Topic  topic.EmitterTopic
	Group  string
	Method goka.ProcessCallback
}

func NewKafkaConsumer(sugar logger.CustomLogger, configKafka config.ConfigSetup, wrokers ...KafkaWorker) Consumer {
	// config
	brokers := map[string][]string{}
	configKafka.GetConfig("kafka", &brokers)
	return &KafkaConsumerImpl{sugar: sugar, brokers: brokers["brokers"], workers: wrokers}
}

func (k *KafkaConsumerImpl) Consume() {
	for _, val := range k.workers {
		g := goka.DefineGroup(goka.Group(val.Group),
			goka.Input(val.Topic.GokaStream(), new(codec.String), val.Method),
			goka.Persist(new(codec.Int64)))

		p, err := goka.NewProcessor(k.brokers, g)
		if err != nil {
			log.Fatalf("error creating processor: %v", err)
		}
		ctx, cancel := context.WithCancel(context.Background())
		done := make(chan bool)
		go func() {
			defer close(done)
			if err = p.Run(ctx); err != nil {
				log.Fatalf("error running processor: %v", err)
			} else {

				log.Printf("Processor shutdown cleanly")
			}
		}()

		wait := make(chan os.Signal, 1)
		signal.Notify(wait, syscall.SIGINT, syscall.SIGTERM)
		<-wait   // wait for SIGINT/SIGTERM
		cancel() // gracefully stop processor
		<-done
	}
}
