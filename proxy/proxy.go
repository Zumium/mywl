package proxy

import "text/template"

const proxyJsonTemplateString = `{ name : "{{.Name}}", protocol : "{{.Protocol}}", address : "{{.Address}}" }`

type Proxy struct {
	name     string
	address  string
	protocol string
}

var proxyJsonTemplate = template.Must(template.New("ProxyJsonTemplate").Parse(proxyJsonTemplateString))

func NewProxy(name, protocol, address string) *Proxy {
	return &Proxy{name: name, address: address, protocol: protocol}
}

func (p *Proxy) ToProxyMethodString() string {
	return p.protocol + " " + p.address
}

func (p *Proxy) ToMap() map[string]string {
	return map[string]string{
		"Name":     p.name,
		"Protocol": p.protocol,
		"Address":  p.address,
	}
}
