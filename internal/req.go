package internal

type RequestValidator string

const (
	RequestValidatorNo   = RequestValidator("no")
	RequestValidatorOzzo = RequestValidator("ozzo")
)

func CreateRequestValidator(val string) RequestValidator {
	if val == string(RequestValidatorOzzo) {
		return RequestValidatorOzzo
	}

	return RequestValidatorNo
}
