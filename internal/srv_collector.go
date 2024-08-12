package internal

import (
	"fmt"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
)

type SrvCollector struct {
}

func NewSrvCollector() *SrvCollector {
	return &SrvCollector{}
}

type CollectOpts struct {
	SrvNaming SrvNaming
	PkgNaming PkgNaming
}

func (c *SrvCollector) Collect(file *protogen.File, opts CollectOpts) (*Services, error) {
	services := &Services{
		Services: []*Service{},
	}

	apiImportPkg := ApiImportPackage{
		FullName: string(file.GoImportPath),
	}

	importPathParts := strings.Split(apiImportPkg.FullName, "/")
	apiImportPkg.Alias = importPathParts[len(importPathParts)-1]

	for _, service := range file.Services {
		srvPkg := c.generatePackageName(service, opts)

		srv := &Service{
			PackageName:      srvPkg,
			Name:             c.generateServiceName(service, opts),
			RpcName:          service.GoName,
			ApiImportPackage: apiImportPkg,
			PbFileName:       srvPkg + "/service.go",
		}

		for _, method := range service.Methods {
			handler := &Handler{
				Filename: fmt.Sprintf(
					"%s/%s.go",
					srvPkg,
					strings.ToLower(method.GoName),
				),
				MethodName:          method.GoName,
				InputMsgStructName:  string(method.Input.Desc.Name()),
				OutputMsgStructName: string(method.Output.Desc.Name()),
			}

			srv.Handlers = append(srv.Handlers, handler)
		}

		services.Services = append(services.Services, srv)
	}

	return services, nil
}

func (*SrvCollector) generatePackageName(srv *protogen.Service, opts CollectOpts) string {
	if opts.PkgNaming == PkgNamingAsIs {
		return strings.ToLower(srv.GoName)
	}

	return strings.TrimSuffix(strings.ToLower(srv.GoName), "service")
}

func (*SrvCollector) generateServiceName(srv *protogen.Service, opts CollectOpts) string {
	if opts.SrvNaming == SrvNamingAsIs {
		return srv.GoName
	}

	return "Service"
}
