package internal

type SrvNaming string

const (
	SrvNamingAsIs        SrvNaming = "as_is"
	SrvNamingJustService SrvNaming = "just_service"
)

func CreateSrvNaming(val string) SrvNaming {
	if val == string(SrvNamingJustService) {
		return SrvNamingJustService
	}

	return SrvNamingAsIs
}

type Services struct {
	Services []*Service
}

type Service struct {
	PackageName string
	Name        string
	RPCName     string
	PbFileName  string

	APIImportPackage APIImportPackage
	Handlers         []*Handler
}

type APIImportPackage struct {
	FullName               string
	Alias                  string
	AliasEqualsLastPackage bool
}

type Handler struct {
	Filename            string
	MethodName          string
	InputMsgStructName  string
	OutputMsgStructName string
}
