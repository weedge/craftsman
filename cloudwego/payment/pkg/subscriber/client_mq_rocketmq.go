package subscriber

import (
	"context"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/rlog"
	"github.com/cloudwego/kitex/pkg/klog"
)

type RmqPushConsumerOptions struct {
	Name            string   `mapstructure:"name"`
	NameSrvs        []string `mapstructure:"nameSrvs"`
	GroupName       string   `mapstructure:"groupName"`
	TopicName       string   `mapstructure:"topicName"`
	Tag             string   `mapstructure:"tag"`
	PullRetryCn     int      `mapstructure:"pullRetryCn"`
	LogicMaxRetryCn int      `mapstructure:"logicMaxRetryCn"`

	// MaxReconsumeTimes retry over to dlq topic, -1 is default 16, 0 or <-1 don't retry,over to dlq topic
	MaxReconsumeTimes int `mapstructure:"maxReconsumeTimes"`

	PullBatchSize              int `mapstructure:"pullBatchSize"`
	ConsumeMessageBatchMaxSize int `mapstructure:"consumeMessageBatchMaxSize"`

	// The DelayLevel specify the waiting time that before next reconsume,
	// and this range is from 1 to 18 now.
	// The time of each level is the value of indexing of {level-1} in
	// [1s, 5s, 10s, 30s, 1m, 2m, 3m, 4m, 5m, 6m, 7m, 8m, 9m, 10m, 20m, 30m, 1h, 2h]
	//delayLevel := 2
	// out delay level range, use default retry, retry cn: maxReconsumeTimes default 16
	// [10s, 30s, 1m, 2m, 3m, 4m, 5m, 6m, 7m, 8m, 9m, 10m, 20m, 30m, 1h, 2h]
	DelayLevel int `mapstructure:"delayLevel"`
}

func DefaultRmqPushConsumerOptions() *RmqPushConsumerOptions {
	return &RmqPushConsumerOptions{
		NameSrvs:                   []string{"127.0.0.1:9876"},
		PullRetryCn:                2,
		LogicMaxRetryCn:            0,
		MaxReconsumeTimes:          20,
		PullBatchSize:              32,
		ConsumeMessageBatchMaxSize: 1,
		DelayLevel:                 -1,
	}
}

type IRocketMQConsumerSubscribeHandler interface {

	// SubMsgsHandle sub some msg to consume
	// maybe u can concurency batch do to improve throughput rate,
	// with min(ConsumeMessageBatchMaxSize,PullBatchSize)>1;
	// but batch done to return commit status
	// so pelease Idempotent consume every msg
	// notice: PullBatchSize/min(ConsumeMessageBatchMaxSize,PullBatchSize) subscribes are concurency
	SubMsgsHandle(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error)
}

var (
	gMapPushConsumer = map[string]rocketmq.PushConsumer{}
)

// InitPushConsumerSubscribes support msg tag subscribe
func InitPushConsumerSubscribes(opts map[string]*RmqPushConsumerOptions, mapSubscribeHandler map[string]IRocketMQConsumerSubscribeHandler) {
	rlog.SetLogLevel("error")
	primitive.PanicHandler = func(i interface{}) { klog.Errorf("[panic] %v", i) }

	for _, opt := range opts {
		if handler, ok := mapSubscribeHandler[opt.Name]; ok {
			c, err := initPushConsumerSubscribe(opt, handler)
			if err != nil {
				continue
			}
			gMapPushConsumer[opt.Name] = c
		}
	}
}

func initPushConsumerSubscribe(opt *RmqPushConsumerOptions, handler IRocketMQConsumerSubscribeHandler) (c rocketmq.PushConsumer, err error) {
	//todo: add otel tracing

	klog.CtxDebugf(context.TODO(), "opt %+v subscribeHandler %+v ", opt, handler)
	traceCfg := &primitive.TraceConfig{
		Access:    primitive.Local,
		Resolver:  primitive.NewPassthroughResolver(opt.NameSrvs),
		GroupName: opt.GroupName,
	}
	c, err = rocketmq.NewPushConsumer(
		consumer.WithGroupName(opt.GroupName),
		consumer.WithConsumerModel(consumer.Clustering),
		consumer.WithNsResolver(primitive.NewPassthroughResolver(opt.NameSrvs)),
		consumer.WithRetry(opt.PullRetryCn),
		consumer.WithTrace(traceCfg),
		consumer.WithMaxReconsumeTimes(int32(opt.MaxReconsumeTimes)),
		consumer.WithPullBatchSize(int32(opt.PullBatchSize)),
		consumer.WithConsumeMessageBatchMaxSize(opt.ConsumeMessageBatchMaxSize),
	)
	if err != nil {
		klog.Errorf("new consumer error: %s", err.Error())
		return
	}

	err = c.Subscribe(opt.TopicName, consumer.MessageSelector{
		Type:       consumer.TAG,
		Expression: opt.Tag,
	}, handler.SubMsgsHandle)
	if err != nil {
		klog.Errorf("sub error: %s", err.Error())
		return
	}

	// Note: start after subscribe
	err = c.Start()
	if err != nil {
		klog.Errorf("consumer %s start error:%s", opt.Name, err.Error())
		return
	}
	klog.Infof("consumer %s start ok", opt.Name)

	return
}

func Close() {
	for name, consumer := range gMapPushConsumer {
		err := consumer.Shutdown()
		if err != nil {
			klog.Errorf("consumer %s shutdown error:%s", name, err.Error())
			return
		}
		klog.Infof("consumer %s shutdown ok", name)
	}
}
