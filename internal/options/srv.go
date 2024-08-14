package options

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
