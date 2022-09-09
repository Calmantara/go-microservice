package producer

import (
	"context"

	"github.com/Calmantara/go-common/pb"
	serviceutil "github.com/Calmantara/go-common/service/util"

	"github.com/Calmantara/go-common/logger"
	"github.com/Calmantara/go-common/setup/config"
	"github.com/Calmantara/go-common/topic"
	"github.com/google/uuid"
	"github.com/lovoo/goka"
	"github.com/lovoo/goka/codec"
	"google.golang.org/protobuf/encoding/protojson"
)

type Producer interface {
	Publish(ctx context.Context, topic topic.EmitterTopic, message any) (err error)
}

type KafkaProducerImpl struct {
	sugar   logger.CustomLogger
	util    serviceutil.UtilService
	brokers []string
}

func NewKafkaProducer(sugar logger.CustomLogger, configKafka config.ConfigSetup, util serviceutil.UtilService) Producer {
	// config
	brokers := map[string][]string{}
	configKafka.GetConfig("kafka", &brokers)
	return &KafkaProducerImpl{sugar: sugar, util: util, brokers: brokers["brokers"]}
}

func (k *KafkaProducerImpl) Publish(ctx context.Context, topic topic.EmitterTopic, message any) (err error) {
	emitter, err := goka.NewEmitter(k.brokers, topic.GokaStream(), new(codec.Bytes))
	if err != nil {
		k.sugar.WithContext(ctx).Errorf("error creating emitter:%v", err)
		return err
	}
	defer emitter.Finish()

	// generate correlation for header
	corrId := k.util.GetCorrelationIdFromContext(ctx)
	header := goka.Headers{logger.CorrelationKey.String(): []byte(corrId)}

	// auto generate unique key
	key := uuid.New()
	k.sugar.WithContext(ctx).Infof("emitting message with key:%v", key)

	// converting to proto
	msg := &pb.Emitter{}
	k.util.ObjectMapper(message, msg)
	buff, _ := protojson.Marshal(msg)

	if err = emitter.EmitSyncWithHeaders(key.String(), buff, header); err != nil {
		k.sugar.WithContext(ctx).Errorf("error emmit message:%v", err)
	}
	return err
}
