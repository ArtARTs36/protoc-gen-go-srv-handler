package main

import (
	_ "embed"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/artarts36/protoc-gen-go-srv-handler/internal/generator"

	"github.com/artarts36/protoc-gen-go-srv-handler/internal/collector"
	"github.com/artarts36/protoc-gen-go-srv-handler/internal/entity"
	"github.com/artarts36/protoc-gen-go-srv-handler/internal/options"
	"github.com/artarts36/protoc-gen-go-srv-handler/internal/renderer"

	"google.golang.org/protobuf/types/pluginpb"

	"google.golang.org/protobuf/compiler/protogen"
)

var (
	outDir                    = flag.String("out_dir", "", "Output directory for generated files")
	overwrite                 = flag.Bool("overwrite", false, "Overwrite existing files")
	genTests                  = flag.Bool("gen_tests", false, "Generate test files")
	pkgNamingVal              = flag.String("pkg_naming", string(options.PkgNamingAsIs), "Package naming: `as_is`, `without_service_suffix`")                          //nolint: lll // not need
	srvNamingVal              = flag.String("srv_naming", string(options.SrvNamingAsIs), "Service naming: `as_is`, `just_service`")                                    //nolint: lll // not need
	handlerFileNamingVal      = flag.String("handler_file_naming", string(options.HandlerFileNamingAsIs), "Handler file naming: `as_is`, `without domain`")            //nolint: lll // not need
	requestValidatorVal       = flag.String("request_validator", string(options.RequestValidatorTypeNo), "Request validator: `no`, `ozzo`")                            //nolint: lll // not need
	requestValidatorFieldsVal = flag.String("request_validator_fields", string(options.RequestValidatorFieldsNonOptional), "Request validator fields: `non_optional`") //nolint: lll // not need
)

func main() {
	flag.Parse()

	protogen.Options{
		ParamFunc: flag.CommandLine.Set,
	}.Run(func(gen *protogen.Plugin) error {
		gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)

		if *outDir == "" {
			return errors.New("out_dir is required, set --go-srv-handler_opt=out_dir=./dir")
		}

		rend, err := renderer.NewRenderer()
		if err != nil {
			return err
		}

		currDir, err := os.Getwd()
		if err != nil {
			return err
		}

		params := generator.GenerateParams{
			OutputDir:         filepath.Join(currDir, *outDir),
			FileOverwrite:     *overwrite,
			PkgNaming:         options.CreatePkgNaming(*pkgNamingVal),
			SrvNaming:         options.CreateSrvNaming(*srvNamingVal),
			HandlerFileNaming: entity.CreateHandlerFileNaming(*handlerFileNamingVal),
			GenTests:          *genTests,
			ReqValidator: options.RequestValidator{
				Type:   options.CreateRequestValidator(*requestValidatorVal),
				Fields: options.CreateRequestValidatorFields(*requestValidatorFieldsVal),
			},
		}

		genr := generator.NewGenerator(collector.NewSrvCollector(), rend)

		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}

			genErr := genr.Generate(gen, f, params)
			if genErr != nil {
				return fmt.Errorf("failed generating proto files: %v", genErr)
			}
		}
		return nil
	})
}
