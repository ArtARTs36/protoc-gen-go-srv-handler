package entity

type SrvNaming string

type Services struct {
	Services []*Service
}

type Service struct {
	PackageName  string
	Name         string
	RPCName      string
	GoFileName   string
	TestFileName string

	Domain string

	APIImportPackage APIImportPackage
	Handlers         map[string]*Handler
}

type APIImportPackage struct {
	FullName               string
	Alias                  string
	AliasEqualsLastPackage bool
}
