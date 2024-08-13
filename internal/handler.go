package internal

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
}

func CreateHandlerFileNaming(val string) HandlerFileNaming {
	if val == string(HandlerFileNamingWithoutDomain) {
		return HandlerFileNamingWithoutDomain
	}

	return HandlerFileNamingAsIs
}
