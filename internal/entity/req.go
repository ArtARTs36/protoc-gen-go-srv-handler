package entity

type RequestValidatorType string

const (
	RequestValidatorTypeNo   = RequestValidatorType("no")
	RequestValidatorTypeOzzo = RequestValidatorType("ozzo")
)

type RequestValidatorFields string

const (
	RequestValidatorFieldsNonOptional = RequestValidatorFields("non_optional")
)

type RequestValidator struct {
	Type   RequestValidatorType
	Fields RequestValidatorFields
}

func CreateRequestValidator(val string) RequestValidatorType {
	if val == string(RequestValidatorTypeOzzo) {
		return RequestValidatorTypeOzzo
	}

	return RequestValidatorTypeNo
}

func CreateRequestValidatorFields(_ string) RequestValidatorFields {
	return RequestValidatorFieldsNonOptional
}
