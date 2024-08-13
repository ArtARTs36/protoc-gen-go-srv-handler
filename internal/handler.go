package internal

import (
	"strings"
)

type HandlerFileNaming string

const (
	HandlerFileNamingAsIs          HandlerFileNaming = "as_is"
	HandlerFileNamingWithoutDomain HandlerFileNaming = "without_domain"
)

type Handler struct {
	Filename            string
	MethodName          string
	InputMsgStructName  string
	OutputMsgStructName string

	Service *Service
}

func (h *Handler) TestFileName() string {
	return strings.Replace(h.Filename, ".go", "_test.go", 1)
}

func CreateHandlerFileNaming(val string) HandlerFileNaming {
	if val == string(HandlerFileNamingWithoutDomain) {
		return HandlerFileNamingWithoutDomain
	}

	return HandlerFileNamingAsIs
}
