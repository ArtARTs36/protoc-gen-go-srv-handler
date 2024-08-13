package internal

import "google.golang.org/protobuf/reflect/protoreflect"

type ValType int8

const (
	ValTypeString ValType = iota
	ValTypeInt    ValType = iota
	ValTypeFloat  ValType = iota
	ValTypeBool   ValType = iota
)

func createValType(t protoreflect.Kind) ValType {
	switch t {
	case protoreflect.StringKind:
		return ValTypeString
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind, protoreflect.Int64Kind,
		protoreflect.Sint64Kind, protoreflect.Sfixed64Kind, protoreflect.Uint32Kind, protoreflect.Fixed32Kind,
		protoreflect.Uint64Kind:
		return ValTypeInt
	case protoreflect.FloatKind, protoreflect.DoubleKind:
		return ValTypeFloat
	case protoreflect.BoolKind:
		return ValTypeBool
	default:
		return ValTypeString
	}
}
