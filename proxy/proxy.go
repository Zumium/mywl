package proxy

type Proxy struct {
	name     string
	address  string
	protocol string
}

func NewProxy(name, protocol, address string) *Proxy {
	return &Proxy{name: name, address: address, protocol: protocol}
}

func (p *Proxy) ToProxyMethodString() string {
	return p.protocol + " " + p.address
}
