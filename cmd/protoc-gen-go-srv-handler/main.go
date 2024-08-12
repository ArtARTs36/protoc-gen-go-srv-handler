package main

import (
	_ "embed"
	"flag"
	"fmt"
	"github.com/artarts36/protoc-gen-go-srv-handler/internal"
	"google.golang.org/protobuf/compiler/protogen"
	"os"
	"path/filepath"
)

func main() {
	var flags flag.FlagSet
	outDir := flags.String("out_dir", "", "Output directory for generated files")
	overwrite := flags.Bool("overwrite", false, "Overwrite existing files")
	pkgNamingVal := flags.String("pkg_naming", string(internal.PkgNamingAsIs), "Package naming: `as_is`, `without_service_suffix`")
	srvNamingVal := flags.String("srv_naming", string(internal.SrvNamingAsIs), "Service naming: `as_is`, `just_service`")

	protogen.Options{
		ParamFunc: flags.Set,
	}.Run(func(gen *protogen.Plugin) error {
		if *outDir == "" {
			return fmt.Errorf("out_dir is required, set --go-srv-handler_opt=out_dir=./dir")
		}

		pkgNaming := internal.CreatePkgNaming(*pkgNamingVal)
		srvNaming := internal.CreateSrvNaming(*srvNamingVal)

		renderer, err := internal.NewRenderer()
		if err != nil {
			return err
		}

		currDir, err := os.Getwd()
		if err != nil {
			return err
		}

		cmd := &command{
			outputDir:    filepath.Join(currDir, *outDir),
			overwrite:    *overwrite,
			pkgNaming:    pkgNaming,
			srvCollector: internal.NewSrvCollector(),
			renderer:     renderer,
			srvNaming:    srvNaming,
		}

		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}

			_, genErr := cmd.gen(gen, f)
			if genErr != nil {
				return fmt.Errorf("failed generating proto files: %v", genErr)
			}
		}
		return nil
	})
}

type command struct {
	outputDir    string
	overwrite    bool
	pkgNaming    internal.PkgNaming
	srvNaming    internal.SrvNaming
	srvCollector *internal.SrvCollector
	renderer     *internal.Renderer
}

func (c *command) gen(gen *protogen.Plugin, file *protogen.File) ([]*protogen.GeneratedFile, error) {
	services, err := c.srvCollector.Collect(file, internal.CollectOpts{
		PkgNaming: c.pkgNaming,
		SrvNaming: c.srvNaming,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to collect services: %w", err)
	}

	var generatedFiles []*protogen.GeneratedFile

	for _, srv := range services.Services {
		filename := srv.PbFileName
		if c.skipFile(filename) {
			continue
		}

		g := gen.NewGeneratedFile(filename, file.GoImportPath)

		err = c.renderer.RenderService(g, srv)
		if err != nil {
			return nil, fmt.Errorf("failed rendering service: %w", err)
		}

		generatedFiles = append(generatedFiles, g)

		for _, handler := range srv.Handlers {
			hf := gen.NewGeneratedFile(handler.Filename, file.GoImportPath)

			err = c.renderer.RenderHandler(hf, srv, handler)
			if err != nil {
				return nil, fmt.Errorf("failed rendering handler: %w", err)
			}

			generatedFiles = append(generatedFiles, hf)
		}
	}

	return generatedFiles, nil
}

func (c *command) skipFile(path string) bool {
	if c.overwrite {
		return false
	}

	path = filepath.Join(c.outputDir, path)
	_, err := os.Stat(path)
	return err == nil
}
