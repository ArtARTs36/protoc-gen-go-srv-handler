package collector

import (
	"github.com/artarts36/protoc-gen-go-srv-handler/internal/options"
)

type CollectOpts struct {
	SrvNaming         options.SrvNaming
	PkgNaming         options.PkgNaming
	HandlerFileNaming options.HandlerFileNaming
	RequestValidator  options.RequestValidator
}
