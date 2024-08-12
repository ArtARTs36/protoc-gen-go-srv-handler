package internal

type Services struct {
	Services []*Service
}

type Service struct {
	PackageName string
	Name        string
	PbFileName  string

	ApiImportPackage ApiImportPackage
	Handlers         []*Handler
}

type ApiImportPackage struct {
	FullName string
	Alias    string
}

type Handler struct {
	Filename            string
	MethodName          string
	InputMsgStructName  string
	OutputMsgStructName string
}
