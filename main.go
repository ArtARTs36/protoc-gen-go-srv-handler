package main

import (
	_ "embed"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/artarts36/protoc-gen-go-srv-handler/internal/collector"
	"github.com/artarts36/protoc-gen-go-srv-handler/internal/entity"
	"github.com/artarts36/protoc-gen-go-srv-handler/internal/options"
	"github.com/artarts36/protoc-gen-go-srv-handler/internal/renderer"

	"google.golang.org/protobuf/types/pluginpb"

	"google.golang.org/protobuf/compiler/protogen"
)

func main() {
	var flags flag.FlagSet
	outDir := flags.String("out_dir", "", "Output directory for generated files")
	overwrite := flags.Bool("overwrite", false, "Overwrite existing files")
	genTests := flags.Bool("gen_tests", false, "Generate test files")
	pkgNamingVal := flags.String(
		"pkg_naming", string(options.PkgNamingAsIs), "Package naming: `as_is`, `without_service_suffix`")
	srvNamingVal := flags.String(
		"srv_naming", string(options.SrvNamingAsIs), "Service naming: `as_is`, `just_service`")
	handlerFileNamingVal := flags.String(
		"handler_file_naming",
		string(options.HandlerFileNamingAsIs),
		"Handler file naming: `as_is`, `without domain`",
	)
	requestValidatorVal := flags.String(
		"request_validator",
		string(entity.RequestValidatorTypeNo),
		"Request validator: `no`, `ozzo`",
	)
	requestValidatorFieldsVal := flags.String(
		"request_validator_fields",
		string(entity.RequestValidatorFieldsNonOptional),
		"Request validator fields: `non_optional`",
	)

	protogen.Options{
		ParamFunc: flags.Set,
	}.Run(func(gen *protogen.Plugin) error {
		gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)

		if *outDir == "" {
			return fmt.Errorf("out_dir is required, set --go-srv-handler_opt=out_dir=./dir")
		}

		pkgNaming := options.CreatePkgNaming(*pkgNamingVal)
		srvNaming := options.CreateSrvNaming(*srvNamingVal)
		handlerFileNaming := entity.CreateHandlerFileNaming(*handlerFileNamingVal)
		reqValidator := entity.CreateRequestValidator(*requestValidatorVal)
		reqValidatorFields := entity.CreateRequestValidatorFields(*requestValidatorFieldsVal)

		rend, err := renderer.NewRenderer()
		if err != nil {
			return err
		}

		currDir, err := os.Getwd()
		if err != nil {
			return err
		}

		cmd := &command{
			outputDir:         filepath.Join(currDir, *outDir),
			overwrite:         *overwrite,
			pkgNaming:         pkgNaming,
			handlerFileNaming: handlerFileNaming,
			genTests:          *genTests,
			srvCollector:      collector.NewSrvCollector(),
			renderer:          rend,
			srvNaming:         srvNaming,
			reqValidator: entity.RequestValidator{
				Type:   reqValidator,
				Fields: reqValidatorFields,
			},
		}

		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}

			genErr := cmd.gen(gen, f)
			if genErr != nil {
				return fmt.Errorf("failed generating proto files: %v", genErr)
			}
		}
		return nil
	})
}

type command struct {
	outputDir         string
	overwrite         bool
	handlerFileNaming options.HandlerFileNaming
	pkgNaming         options.PkgNaming
	srvNaming         options.SrvNaming
	reqValidator      entity.RequestValidator
	srvCollector      *collector.SrvCollector
	renderer          *renderer.Renderer
	genTests          bool
}

func (c *command) gen(gen *protogen.Plugin, file *protogen.File) error {
	services, err := c.srvCollector.Collect(file, collector.CollectOpts{
		PkgNaming:         c.pkgNaming,
		SrvNaming:         c.srvNaming,
		HandlerFileNaming: c.handlerFileNaming,
		RequestValidator:  c.reqValidator,
	})
	if err != nil {
		return fmt.Errorf("failed to collect services: %w", err)
	}

	for _, srv := range services.Services {
		if c.skipFile(srv.GoFileName) {
			continue
		}

		srvGenFile := gen.NewGeneratedFile(srv.GoFileName, file.GoImportPath)

		err = c.renderer.RenderService(srvGenFile, srv)
		if err != nil {
			return fmt.Errorf("failed rendering service: %w", err)
		}

		if c.genTests {
			srvTestGenFile := gen.NewGeneratedFile(srv.TestFileName, file.GoImportPath)
			if !c.skipFile(srv.TestFileName) {
				err = c.renderer.RenderServiceTest(srvTestGenFile, srv)
				if err != nil {
					return fmt.Errorf("failed rendering service test: %w", err)
				}
			}
		}

		err = c.genHandlers(gen, file, srv)
		if err != nil {
			return fmt.Errorf("failed generating handlers: %w", err)
		}
	}

	return nil
}

func (c *command) genHandlers(gen *protogen.Plugin, file *protogen.File, srv *entity.Service) error {
	for _, handler := range srv.Handlers {
		handlerGenFile := gen.NewGeneratedFile(handler.Filename, file.GoImportPath)

		err := c.renderer.RenderHandler(handlerGenFile, handler, renderer.RenderHandlerParams{
			RequestValidator: c.reqValidator,
		})
		if err != nil {
			return fmt.Errorf("failed rendering handler: %w", err)
		}

		if c.genTests {
			err = c.genHandlerTest(gen, file, handler)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *command) genHandlerTest(gen *protogen.Plugin, file *protogen.File, handler *entity.Handler) error {
	handlerTestGenFile := gen.NewGeneratedFile(handler.TestFileName(), file.GoImportPath)
	if !c.skipFile(handler.TestFileName()) {
		err := c.renderer.RenderHandlerTest(handlerTestGenFile, handler, renderer.RenderHandlerParams{
			RequestValidator: c.reqValidator,
		})
		if err != nil {
			return fmt.Errorf("failed rendering handler test: %w", err)
		}
	}

	return nil
}

func (c *command) skipFile(path string) bool {
	if c.overwrite {
		return false
	}

	path = filepath.Join(c.outputDir, path)
	_, err := os.Stat(path)
	return err == nil
}
