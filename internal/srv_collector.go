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

func (c *SrvCollector) Collect(file *protogen.File) (*Services, error) {
	services := &Services{
		Services: []*Service{},
	}

	apiImportPkg := ApiImportPackage{
		FullName: string(file.GoImportPath),
	}

	importPathParts := strings.Split(apiImportPkg.FullName, "/")
	apiImportPkg.Alias = importPathParts[len(importPathParts)-1]

	for _, service := range file.Services {
		srv := &Service{
			PackageName:      strings.ToLower(service.GoName),
			Name:             service.GoName,
			ApiImportPackage: apiImportPkg,
			PbFileName:       strings.ToLower(service.GoName) + "/service.go",
		}

		for _, method := range service.Methods {
			handler := &Handler{
				Filename: fmt.Sprintf(
					"%s/%s.go",
					strings.ToLower(service.GoName),
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
