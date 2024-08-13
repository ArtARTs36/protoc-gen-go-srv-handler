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
	SrvNaming         SrvNaming
	PkgNaming         PkgNaming
	HandlerFileNaming HandlerFileNaming
}

type handlerFileNames struct {
	selected      string
	asIs          string
	withoutDomain string
}

func (c *SrvCollector) Collect(file *protogen.File, opts CollectOpts) (*Services, error) {
	services := &Services{
		Services: []*Service{},
	}

	apiImportPkg := APIImportPackage{
		FullName: string(file.GoImportPath),
		Alias:    string(file.GoPackageName),
	}

	if apiImportPkg.Alias == "" {
		importPathParts := strings.Split(apiImportPkg.FullName, "/")
		apiImportPkg.Alias = importPathParts[len(importPathParts)-1]
		apiImportPkg.AliasEqualsLastPackage = true
	} else {
		importPathParts := strings.Split(apiImportPkg.FullName, "/")
		apiImportPkg.AliasEqualsLastPackage = apiImportPkg.Alias == importPathParts[len(importPathParts)-1]
	}

	for _, service := range file.Services {
		srvPkg := c.generatePackageName(service, opts)

		srv := &Service{
			PackageName:      srvPkg,
			Name:             c.generateServiceName(service, opts),
			Domain:           strings.TrimSuffix(service.GoName, "Service"),
			RPCName:          service.GoName,
			APIImportPackage: apiImportPkg,
			GoFileName:       srvPkg + "/service.go",
			TestFileName:     srvPkg + "/service_test.go",
		}

		handlersByFiles := map[string]handlerFileNames{}

		srv.Handlers = map[string]*Handler{}

		for _, method := range service.Methods {
			names := c.generateHandlerFilename(method, srv, srvPkg, opts)
			if otherHandlerNames, exists := handlersByFiles[names.selected]; exists {
				otherHandler, alreadyReplaced := srv.Handlers[names.selected]
				if !alreadyReplaced {
					otherHandler.Filename = otherHandlerNames.asIs
					delete(srv.Handlers, names.selected)
					srv.Handlers[otherHandlerNames.asIs] = otherHandler
				}

				names.selected = names.asIs
			}

			handlersByFiles[names.selected] = names

			handler := &Handler{
				Filename:            names.selected,
				MethodName:          method.GoName,
				InputMsgStructName:  string(method.Input.Desc.Name()),
				OutputMsgStructName: string(method.Output.Desc.Name()),
				Service:             srv,
			}

			srv.Handlers[handler.Filename] = handler
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

func (*SrvCollector) generateHandlerFilename(
	method *protogen.Method,
	srv *Service,
	pkg string,
	opts CollectOpts,
) handlerFileNames {
	names := handlerFileNames{
		asIs: fmt.Sprintf(
			"%s/%s.go",
			pkg,
			strings.ToLower(method.GoName),
		),
		withoutDomain: fmt.Sprintf(
			"%s/%s.go",
			pkg,
			strings.ReplaceAll(strings.ToLower(method.GoName), strings.ToLower(srv.Domain), ""),
		),
	}

	if opts.HandlerFileNaming == HandlerFileNamingAsIs {
		names.selected = names.asIs
	} else {
		names.selected = names.withoutDomain
	}

	return names
}
