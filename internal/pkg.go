package internal

type PkgNaming string

const (
	PkgNamingAsIs                 PkgNaming = "as_is"
	PkgNamingWithoutServiceSuffix PkgNaming = "without_service_suffix"
)

func CreatePkgNaming(val string) PkgNaming {
	if val == string(PkgNamingWithoutServiceSuffix) {
		return PkgNamingWithoutServiceSuffix
	}

	return PkgNamingAsIs
}
