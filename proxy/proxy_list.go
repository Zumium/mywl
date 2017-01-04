package proxy

import (
	"bufio"
	"container/list"
	"errors"
	"flag"
	"github.com/Zumium/mywl/common"
	"os"
	"strings"
)

type ProxyList struct {
	current   *Proxy
	proxylist *list.List
}

var proxyConfigFilePath string

func (pl *ProxyList) Add(name, protocol, address string) {
	pl.proxylist.PushBack(NewProxy(name, protocol, address))
}

func (pl *ProxyList) Find(name string) (common.Proxy, error) {
	for e := pl.proxylist.Front(); e != nil; e = e.Next() {
		proxy, _ := e.Value.(*Proxy)
		if proxy.name == name {
			return proxy, nil
		}
	}
	return nil, errors.New("No such proxy named " + name)
}

func (pl *ProxyList) Del(name string) error {
	for e := pl.proxylist.Front(); e != nil; e = e.Next() {
		proxy, _ := e.Value.(*Proxy)
		if proxy.name == name {
			pl.proxylist.Remove(e)
			return nil
		}
	}
	return errors.New("No such proxy named " + name)
}

func (pl *ProxyList) Len() int {
	return pl.proxylist.Len()
}

func (pl *ProxyList) SetCurrent(name string) error {
	if proxy, err := pl.Find(name); err != nil {
		return err
	} else {
		p, _ := proxy.(*Proxy)
		pl.current = p
	}
	return nil
}

func (pl *ProxyList) GetCurrent() common.Proxy {
	return pl.current
}

func (pl *ProxyList) Set(name, protocol, address string) error {
	if proxy, err := pl.Find(name); err != nil {
		return err
	} else {
		p, _ := proxy.(*Proxy)
		if protocol != "" {
			p.protocol = protocol
		}
		if address != "" {
			p.address = address
		}
	}
	return nil
}

func (pl *ProxyList) InstallFlags(flagset *flag.FlagSet) {
	flagset.StringVar(&proxyConfigFilePath, "proxyfile", "/etc/mywl/proxyconfigs.txt", "file that saves proxy configurations")
}

func (pl *ProxyList) Init() error {
	//var absPath string
	//if path, err := filepath.Abs(proxyConfigFilePath); err != nil {
	//	return err
	//} else {
	pl.proxylist = list.New()
	//	absPath = path
	//}
	var proxylistFile *os.File
	defer func() {
		if proxylistFile != nil {
			proxylistFile.Close()
		}
	}()
	//if file, err := os.Open(absPath); err != nil {
	if file, err := os.Open(proxyConfigFilePath); err != nil {
		return nil
	} else {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			segs := strings.Split(scanner.Text(), " ")
			pl.proxylist.PushBack(NewProxy(segs[0], segs[1], segs[2]))
			if len(segs) == 4 && segs[3] == "current" {
				pl.SetCurrent(segs[0])
			}
		}
		if err = scanner.Err(); err != nil {
			return err
		}
	}
	return nil
}

func (pl *ProxyList) ForEach(f func(each common.Proxy)) {
	for e := pl.proxylist.Front(); e != nil; e = e.Next() {
		p, _ := e.Value.(*Proxy)
		f(p)
	}
}
