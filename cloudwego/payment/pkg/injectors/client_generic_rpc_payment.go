package injectors

import (
	"os"
	"strings"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/genericclient"
	"github.com/cloudwego/kitex/pkg/generic"
	"github.com/cloudwego/kitex/pkg/klog"
	commonConstants "github.com/weedge/craftsman/cloudwego/common/pkg/constants"
	"github.com/weedge/craftsman/cloudwego/payment/pkg/constants"
)

type HttpThriftGenericClientOpts struct {
	IdlDirPath string `mapstructure:"idlDirPath"`
	Endpoint   string `mapstructure:"endpoint"`
}

// InitHttpThriftGenericClients client need thrift idl provider
func InitHttpThriftGenericClients(opts HttpThriftGenericClientOpts) (mapGenericClient map[string]genericclient.Client) {
	mapGenericClient = map[string]genericclient.Client{}

	// idl dir , don't have sub dir
	files, err := os.ReadDir(opts.IdlDirPath)
	if err != nil {
		klog.Fatalf("generic routers add idl err: %s", err.Error())
	}

	for _, file := range files {
		fileName := file.Name()
		if file.IsDir() {
			klog.Infof("dir %s can't include", fileName)
			continue
		}

		idlStationName := strings.Join([]string{constants.ProjectName, commonConstants.StationServiceName, commonConstants.IdlFileSuffixThrift}, ".")
		idlDaName := strings.Join([]string{constants.ProjectName, commonConstants.StationServiceName, commonConstants.IdlFileSuffixThrift}, ".")
		if fileName != idlStationName &&
			fileName != idlDaName {
			klog.Infof("file %s can't include", fileName)
			continue
		}

		provider, err := generic.NewThriftFileProvider(fileName, opts.IdlDirPath)
		if err != nil {
			klog.Errorf("generic.NewThriftFileProvider %s %s err: %s", fileName, opts.IdlDirPath, err.Error())
			continue
		}

		httpThriftGeneric, err := generic.HTTPThriftGeneric(provider)
		if err != nil {
			klog.Errorf("generic.HTTPThriftGeneric %s %s err: %s", fileName, opts.IdlDirPath, err.Error())
			continue
		}

		svcName := strings.ReplaceAll(fileName, ".thrift", "")
		genericClient, err := genericclient.NewClient(
			svcName,
			httpThriftGeneric,
			client.WithHostPorts(opts.Endpoint),
		)
		if err != nil {
			klog.Errorf("genericclient.NewClient %s err: %s", svcName, err.Error())
			continue
		}

		mapGenericClient[svcName] = genericClient
	}

	return
}

// InitHttpThriftGenericClients client,server need thrift idl provider
func InitMapThriftGenericClients() {}

// InitHttpThriftGenericClients client,server don't need thrift idl provider
func InitBinaryThriftGenericClients() {}
