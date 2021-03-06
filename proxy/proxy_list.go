package proxy

import (
	"bufio"
	"container/list"
	"errors"
	"flag"
	"fmt"
	"github.com/Zumium/mywl/common"
	"os"
	"strings"
	"text/template"
)

type ProxyList struct {
	current   *Proxy
	proxylist *list.List
}

var proxyConfigFilePath string
var proxyRecordTemplate = template.Must(template.New("ProxyRecord").Parse(`{{.Name}} {{.Protocol}} {{.Address}}{{if .Current}} current{{end}}` + "\n"))

func (pl *ProxyList) Add(name, protocol, address string) {
	pl.proxylist.PushBack(NewProxy(name, protocol, address))

	if err := pl.Save(); err != nil {
		fmt.Fprintf(os.Stderr, "error occurd on proxylist's saving process: %s\n", err.Error())
	}
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
			if proxy == pl.current {
				pl.current = NewProxy("DIRECT", "DIRECT", "")
			}
			pl.proxylist.Remove(e)

			if err := pl.Save(); err != nil {
				fmt.Fprintf(os.Stderr, "error occurd on proxylist's saving process: %s\n", err.Error())
			}

			return nil
		}
	}
	return errors.New("No such proxy named " + name)
}

func (pl *ProxyList) Len() int {
	return pl.proxylist.Len()
}

func (pl *ProxyList) SetCurrent(name string) error {
	proxy, err := pl.Find(name)
	if err != nil {
		return err
	}
	p, _ := proxy.(*Proxy)
	pl.current = p

	if err := pl.Save(); err != nil {
		fmt.Fprintf(os.Stderr, "error occurd on proxylist's saving process: %s\n", err.Error())
	}

	return nil
}

func (pl *ProxyList) GetCurrent() common.Proxy {
	return pl.current
}

func (pl *ProxyList) Set(name, protocol, address string) error {
	proxy, err := pl.Find(name)
	if err != nil {
		return err
	}
	p, _ := proxy.(*Proxy)
	change := false
	if protocol != "" {
		p.protocol = protocol
		change = true
	}
	if address != "" {
		p.address = address
		change = true
	}

	if change {
		if err := pl.Save(); err != nil {
			fmt.Fprintf(os.Stderr, "error occurd on proxylist's saving process: %s\n", err.Error())
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

func (pl *ProxyList) InstallFlags(flagset *flag.FlagSet) {
	flagset.StringVar(&proxyConfigFilePath, "proxyfile", "/etc/mywl/proxyconfigs.txt", "file that saves proxy configurations")
}

func (pl *ProxyList) Init() error {
	pl.proxylist = list.New()
	return pl.Load()
}

func (pl *ProxyList) Load() error {
	pl.proxylist.Init()

	file, err := os.Open(proxyConfigFilePath)
	if err != nil {
		return err
	}
	defer func() {
		if file != nil {
			file.Close()
		}
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		segs := strings.Split(scanner.Text(), " ")
		p := NewProxy(segs[0], segs[1], segs[2])
		pl.proxylist.PushBack(p)
		if len(segs) == 4 && segs[3] == "current" {
			pl.current = p
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	//In case that there is no "current" proxy
	if pl.current == nil {
		pl.current = NewProxy("DIRECT", "DIRECT", "")
	}
	return nil
}

func (pl *ProxyList) Save() error {
	file, err := os.Create(proxyConfigFilePath)
	if err != nil {
		return err
	}
	defer func() {
		if file != nil {
			file.Close()
		}
	}()
	//writer := bufio.NewWriter(file)
	for e := pl.proxylist.Front(); e != nil; e = e.Next() {
		p, _ := e.Value.(*Proxy)
		//proxyRecordTemplate.Execute(writer,map[string]string{"Name": p.name, "Protocol": p.protocol, "Address": p.address, "Current": pl.current.name})
		current := ""
		if p.name == pl.current.name {
			current = "true"
		}
		proxyRecordTemplate.Execute(file, map[string]string{"Name": p.name, "Protocol": p.protocol, "Address": p.address, "Current": current})
	}
	//if err := writer.Flush(); err != nil {
	//	return err
	//}
	return nil
}
