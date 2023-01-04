package injectors

import (
	"os"
	"strings"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/genericclient"
	"github.com/cloudwego/kitex/pkg/generic"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/xds"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	xdsmanager "github.com/kitex-contrib/xds"
	"github.com/kitex-contrib/xds/xdssuite"
	commonConstants "github.com/weedge/craftsman/cloudwego/common/pkg/constants"
	"github.com/weedge/craftsman/cloudwego/common/pkg/metadata"
)

type HttpThriftGenericClientOpts struct {
	IdlDirPath string                  `mapstructure:"idlDirPath"`
	Endpoints  []*GenericEndpointsOpts `mapstructure:"endpoints"`
}

type GenericEndpointsOpts struct {
	ProjectName   string                       `mapstructure:"projectName"`
	SvcName       string                       `mapstructure:"svcName"`
	Version       string                       `mapstructure:"version"`
	HeaderKeys    []string                     `mapstructure:"headerKeys"`
	HostPorts     []string                     `mapstructure:"hostPorts"`
	ClosedMethods []GenericEndpointMethodsOpts `mapstructure:"closedMethods"`

	EnableXDS bool   `mapstructure:"enableXDS"`
	XDSAddr   string `mapstructure:"xdsAddr"`
	Endpoint  string `mapstructure:"endpoint"`
}
type GenericEndpointMethodsOpts struct {
	SvcMethod  string `mapstructure:"svcMethod"`
	HttpMethod string `mapstructure:"httpMethod"`
}

func DefaultHttpThriftGenericClientOpts() *HttpThriftGenericClientOpts {
	return &HttpThriftGenericClientOpts{
		IdlDirPath: os.Getenv(commonConstants.IdlDirPathEnv),
		Endpoints:  []*GenericEndpointsOpts{},
	}
}

func GetApiSvcKey(projectName, svcName, version string) string {
	key := strings.Join([]string{projectName, svcName, version}, ".")
	return key
}
func InitHttpThriftGenericEndpointsOpts(opts *HttpThriftGenericClientOpts) map[string]*GenericEndpointsOpts {
	mapSvcConf := map[string]*GenericEndpointsOpts{}
	for _, item := range opts.Endpoints {
		svcName := GetApiSvcKey(item.ProjectName, item.SvcName, item.Version)
		mapSvcConf[svcName] = item
	}

	return mapSvcConf
}

// InitHttpThriftGenericClients client need thrift idl provider
// just for out open api defined in idl, idl file ci git pull from common, config cd to run
// u can use generic + xds
func InitHttpThriftGenericClients(opts *HttpThriftGenericClientOpts) map[string]genericclient.Client {
	mapGenericClient := map[string]genericclient.Client{}

	for _, item := range opts.Endpoints {
		svcName := GetApiSvcKey(item.ProjectName, item.SvcName, item.Version)
		fileName := strings.Join([]string{item.ProjectName, item.SvcName, commonConstants.IdlFileSuffixThrift}, ".")

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

		genericClient, err := genericclient.NewClient(
			svcName,
			httpThriftGeneric,
			client.WithHostPorts(item.HostPorts...),
			client.WithSuite(tracing.NewClientSuite()),
		)
		if item.EnableXDS {
			err = xdsmanager.Init(xdsmanager.WithXDSServerAddress(item.XDSAddr))
			if err != nil {
				klog.Errorf("%s xdsmanager.WithXDSServerAddress init err: %s", svcName, err.Error())
				continue
			}
			genericClient, err = genericclient.NewClient(
				item.Endpoint,
				httpThriftGeneric,
				client.WithSuite(tracing.NewClientSuite()),
				client.WithXDSSuite(xds.ClientSuite{
					RouterMiddleware: xdssuite.NewXDSRouterMiddleware(
						xdssuite.WithRouterMetaExtractor(metadata.ExtractFromPropagator),
					),
					Resolver: xdssuite.NewXDSResolver(),
				}),
			)
		}
		if err != nil {
			klog.Errorf("genericclient.NewClient %s err: %s", svcName, err.Error())
			continue
		}

		klog.Infof("svcName %s fileName %s InitHttpThriftGenericClient ok", svcName, fileName)
		mapGenericClient[svcName] = genericClient
	}

	return mapGenericClient
}

// InitHttpThriftGenericClients client,server need thrift idl provider
func InitMapThriftGenericClients() {}

// InitHttpThriftGenericClients client,server don't need thrift idl provider
func InitBinaryThriftGenericClients() {}
