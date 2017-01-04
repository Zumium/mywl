package common

import "flag"

type FlagConfigurable interface {
	InstallFlags(flagset *flag.FlagSet)
}

type Initable interface {
	Init() error
}

type WhiteList interface {
	Add(url string)
	Del(url string)
	Has(url string) bool
	ToJsArray() string
}

type Proxy interface {
	ToProxyMethodString() string
	ToMap() map[string]string
}

type ProxyList interface {
	Add(name, protocol, address string)
	Find(name string) (Proxy, error)
	Del(name string) error
	Len() int
	SetCurrent(name string) error
	GetCurrent() Proxy
	Set(name, protocol, address string) error
	ForEach(f func(each Proxy))
}

type Server interface {
	Start()
}

type Persistable interface {
	Load() error
	Save() error
}
