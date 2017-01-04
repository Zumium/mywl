package proxy

import "sync"

var proxylistInstance *ProxyList
var once sync.Once

func GetProxyList() *ProxyList {
  once.Do(func() {
    proxylistInstance = new(ProxyList)
  })
  return proxylistInstance
}
