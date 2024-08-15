package generator

import (
	"fmt"
	"os"
	"path/filepath"

	"google.golang.org/protobuf/compiler/protogen"

	"github.com/artarts36/protoc-gen-go-srv-handler/internal/collector"
	"github.com/artarts36/protoc-gen-go-srv-handler/internal/entity"
	"github.com/artarts36/protoc-gen-go-srv-handler/internal/options"
	"github.com/artarts36/protoc-gen-go-srv-handler/internal/renderer"
)

type Generator struct {
	srvCollector *collector.SrvCollector
	renderer     *renderer.Renderer
}

type GenerateParams struct {
	OutputDir         string
	FileOverwrite     bool
	HandlerFileNaming options.HandlerFileNaming
	PkgNaming         options.PkgNaming
	SrvNaming         options.SrvNaming
	ReqValidator      options.RequestValidator
	GenTests          bool
}

func NewGenerator(srvCollector *collector.SrvCollector, renderer *renderer.Renderer) *Generator {
	return &Generator{
		srvCollector: srvCollector,
		renderer:     renderer,
	}
}

func (c *Generator) Generate(gen *protogen.Plugin, file *protogen.File, params GenerateParams) error {
	services, err := c.srvCollector.Collect(file, collector.CollectOpts{
		PkgNaming:         params.PkgNaming,
		SrvNaming:         params.SrvNaming,
		HandlerFileNaming: params.HandlerFileNaming,
		RequestValidator:  params.ReqValidator,
	})
	if err != nil {
		return fmt.Errorf("failed to collect services: %w", err)
	}

	for _, srv := range services.Services {
		if c.skipFile(srv.GoFileName, params) {
			continue
		}

		srvGenFile := gen.NewGeneratedFile(srv.GoFileName, file.GoImportPath)

		err = c.renderer.RenderService(srvGenFile, srv)
		if err != nil {
			return fmt.Errorf("failed rendering service: %w", err)
		}

		if params.GenTests {
			srvTestGenFile := gen.NewGeneratedFile(srv.TestFileName, file.GoImportPath)
			if !c.skipFile(srv.TestFileName, params) {
				err = c.renderer.RenderServiceTest(srvTestGenFile, srv)
				if err != nil {
					return fmt.Errorf("failed rendering service test: %w", err)
				}
			}
		}

		err = c.genHandlers(gen, file, srv, params)
		if err != nil {
			return fmt.Errorf("failed generating handlers: %w", err)
		}
	}

	return nil
}

func (c *Generator) genHandlers(
	gen *protogen.Plugin,
	file *protogen.File,
	srv *entity.Service,
	params GenerateParams,
) error {
	for _, handler := range srv.Handlers {
		handlerGenFile := gen.NewGeneratedFile(handler.Filename, file.GoImportPath)

		err := c.renderer.RenderHandler(handlerGenFile, handler, renderer.RenderHandlerParams{
			RequestValidator: params.ReqValidator,
		})
		if err != nil {
			return fmt.Errorf("failed rendering handler: %w", err)
		}

		if params.GenTests {
			err = c.genHandlerTest(gen, file, handler, params)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *Generator) genHandlerTest(
	gen *protogen.Plugin,
	file *protogen.File,
	handler *entity.Handler,
	params GenerateParams,
) error {
	handlerTestGenFile := gen.NewGeneratedFile(handler.TestFileName(), file.GoImportPath)
	if !c.skipFile(handler.TestFileName(), params) {
		err := c.renderer.RenderHandlerTest(handlerTestGenFile, handler, renderer.RenderHandlerParams{
			RequestValidator: params.ReqValidator,
		})
		if err != nil {
			return fmt.Errorf("failed rendering handler test: %w", err)
		}
	}

	return nil
}

func (c *Generator) skipFile(path string, params GenerateParams) bool {
	if params.FileOverwrite {
		return false
	}

	path = filepath.Join(params.OutputDir, path)
	_, err := os.Stat(path)
	return err == nil
}
