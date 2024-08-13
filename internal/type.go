package internal

import "google.golang.org/protobuf/reflect/protoreflect"

type ValType int8

const (
	ValTypeUndefined ValType = iota
	ValTypeString    ValType = iota
	ValTypeInt       ValType = iota
	ValTypeFloat     ValType = iota
	ValTypeBool      ValType = iota
)

func createValType(t protoreflect.Kind) ValType {
	switch t { //nolint: exhaustive // not need
	case protoreflect.StringKind:
		return ValTypeString
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind, protoreflect.Int64Kind,
		protoreflect.Sint64Kind, protoreflect.Sfixed64Kind, protoreflect.Uint32Kind, protoreflect.Fixed32Kind,
		protoreflect.Uint64Kind:
		return ValTypeInt
	case protoreflect.FloatKind, protoreflect.DoubleKind, protoreflect.Fixed64Kind:
		return ValTypeFloat
	case protoreflect.BoolKind:
		return ValTypeBool
	default:
		return ValTypeUndefined
	}
}
