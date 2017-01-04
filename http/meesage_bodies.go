package http

type ProxiesPostBody struct {
	Name     string
	Protocol string
	Address  string
}

type ProxiesCurrentPatchBody struct {
	Name string
}

type WhitelistPatchBody struct {
	Operation string
	Url       string
}
