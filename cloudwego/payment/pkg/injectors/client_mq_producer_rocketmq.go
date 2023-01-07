package injectors

import (
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/apache/rocketmq-client-go/v2/rlog"
	"github.com/cloudwego/kitex/pkg/klog"
)

type RmqProducerClientOptions struct {
	Name      string   `mapstructure:"name"`
	NameSrvs  []string `mapstructure:"nameSrvs"`
	GroupName string   `mapstructure:"groupName"`
	RetryCn   int      `mapstructure:"retryCn"`
}

func DefaultRmqProducerClientOptions() *RmqProducerClientOptions {
	return &RmqProducerClientOptions{
		Name:      "",
		NameSrvs:  []string{"127.0.0.1:9876"},
		GroupName: "",
		RetryCn:   2,
	}
}

func InitRmqTransactionProducerClient(opts *RmqProducerClientOptions, listener primitive.TransactionListener) (p rocketmq.TransactionProducer, err error) {
	rlog.SetLogLevel("error")
	primitive.PanicHandler = func(i interface{}) { klog.Errorf("[panic] %v", i) }

	namesrvs := opts.NameSrvs
	groupName := opts.GroupName
	traceCfg := &primitive.TraceConfig{
		Access:    primitive.Local,
		Resolver:  primitive.NewPassthroughResolver(namesrvs),
		GroupName: groupName,
	}

	p, err = rocketmq.NewTransactionProducer(
		listener,
		producer.WithNsResolver(primitive.NewPassthroughResolver(namesrvs)),
		producer.WithGroupName(groupName),
		producer.WithRetry(opts.RetryCn),
		producer.WithTrace(traceCfg),
	)
	if err != nil {
		klog.Fatalf("NewTransactionProducer error: %s", err.Error())
	}

	err = p.Start()
	if err != nil {
		klog.Fatalf("start producer error: %s", err.Error())
	}

	klog.Infof("%s start produce ok", opts.Name)

	return
}
