package entity

import (
	"strings"

	"github.com/artarts36/protoc-gen-go-srv-handler/internal/options"
)

type Message struct {
	Name       string
	Properties MessageProperties
}

type MessageProperties struct {
	All          []*MessageProperty
	Required     []*MessageProperty
	Validateable []*MessageProperty
}

type MessageProperty struct {
	GoName   string
	Type     ValType
	Required bool
	Optional bool
}

func (p *MessageProperty) ExampleValue() string {
	switch p.Type {
	case ValTypeUndefined:
		return "undefined"
	case ValTypeString:
		str := "test"

		if strings.EqualFold(p.GoName, "name") {
			str = "John"
		} else if strings.EqualFold(p.GoName, "email") {
			str = "john@gmail.com"
		}

		return str
	case ValTypeInt:
		return "10"
	case ValTypeBool:
		return "true"
	case ValTypeFloat:
		return "3.1415926"
	}

	return "test"
}

type Handler struct {
	Filename            string
	MethodName          string
	InputMsgStructName  string
	InputMsg            Message
	OutputMsgStructName string

	Service *Service
}

func (h *Handler) TestFileName() string {
	return strings.Replace(h.Filename, ".go", "_test.go", 1)
}

func CreateHandlerFileNaming(val string) options.HandlerFileNaming {
	if val == string(options.HandlerFileNamingWithoutDomain) {
		return options.HandlerFileNamingWithoutDomain
	}

	return options.HandlerFileNamingAsIs
}
