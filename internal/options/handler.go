package options

type HandlerFileNaming string

const (
	HandlerFileNamingAsIs          HandlerFileNaming = "as_is"
	HandlerFileNamingWithoutDomain HandlerFileNaming = "without_domain"
)
