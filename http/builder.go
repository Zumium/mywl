package http

import (
	"encoding/json"
	"flag"
	"github.com/Zumium/mywl/common"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
	"sync"
	"text/template"
)

const pacTemplate = `
var PROXY_METHOD = "{{.Proxymethod}}";
var RULES = [
[".cn"],
{{.Liststring}}
];
function FindProxyForURL(url, host) {

    function check_ipv4() {
        // check if the ipv4 format (TODO: ipv6)
        var re_ipv4 = /^\d+\.\d+\.\d+\.\d+$/g;
        if (re_ipv4.test(host)) {
            // in theory, we can add chnroutes test here.
            // but that is probably too much an overkill.
            return true;
        }
    }

    function isDomain(domain) {
        var host_length, domain_length;
        return ((domain[0] === '.') ? (host === domain.slice(1) || ((host_length = host.length) >= (domain_length = domain.length) && host.slice(host_length - domain_length) === domain)) : (host === domain));
    }

    function rule_filter(callback) {
        // IMPORTANT: Respect the order of RULES.
        for (var j = 0; j < RULES.length; j++) {
            var rules=RULES[j]
            for (var i = 0; i < rules.length; i++) {
               if (callback(rules[i]) === true) {
                   return true;
               }
            }
        }
        return false;
    }

    // skip local hosts
    if (isPlainHostName(host) === true || check_ipv4() === true || rule_filter(isDomain) === true) {
        return "DIRECT";

    } else {
            // if none of above cases, it is always safe to use the proxy
            return PROXY_METHOD;
    }

}`

type ServerBuilder struct {
	bindPort  int
	bindAddr  string
	whitelist common.WhiteList
	proxylist common.ProxyList
}

var serverBuilderInstance *ServerBuilder
var once sync.Once

func GetBuilder() *ServerBuilder {
	once.Do(func() {
		serverBuilderInstance = new(ServerBuilder)
	})
	return serverBuilderInstance
}

func (b *ServerBuilder) InstallFlags(flagset *flag.FlagSet) {
	flagset.IntVar(&b.bindPort, "port", 7000, "the port that the server binds to")
	flagset.StringVar(&b.bindAddr, "address", "127.0.0.1", "the address that the server binds to")
}

func (b *ServerBuilder) SetPort(port int) *ServerBuilder {
	b.bindPort = port
	return b
}

func (b *ServerBuilder) SetAddr(addr string) *ServerBuilder {
	b.bindAddr = addr
	return b
}

func (b *ServerBuilder) SetWhiteList(l common.WhiteList) *ServerBuilder {
	b.whitelist = l
	return b
}

func (b *ServerBuilder) SetProxyList(pl common.ProxyList) *ServerBuilder {
	b.proxylist = pl
	return b
}

func (b *ServerBuilder) Build() common.Server {
	newServer := new(Server)
	newServer.httpserver = echo.New()
	newServer.listenAddr = b.bindAddr + ":" + strconv.Itoa(b.bindPort)

	t, err := template.New("PacTemplate").Parse(pacTemplate)
	if err != nil {
		panic(err)
	}

	//Get PAC file
	newServer.httpserver.GET("/pac", func(c echo.Context) error {
		c.Response().Header().Set("Content-Type", "application/x-ns-proxy-autoconfig")
		return t.Execute(c.Response(), map[string]string{"Proxymethod": b.proxylist.GetCurrent().ToProxyMethodString(), "Liststring": b.whitelist.ToJsArray()})
	})

	//GET all proxy settings
	newServer.httpserver.GET("/proxies", func(c echo.Context) error {
		proxiesArray := make([]map[string]string, 0, b.proxylist.Len())
		b.proxylist.ForEach(func(each common.Proxy) {
			proxiesArray = append(proxiesArray, each.ToMap())
		})
		return c.JSON(http.StatusOK, proxiesArray)
	})

	//GET current proxy setting
	newServer.httpserver.GET("/proxies/current", func(c echo.Context) error {
		return c.JSON(http.StatusOK, b.proxylist.GetCurrent().ToMap())
	})

	//GET :name proxy
	newServer.httpserver.GET("/proxies/:name", func(c echo.Context) error {
		p, err := b.proxylist.Find(c.Param("name"))
		if err != nil {
			return c.JSON(http.StatusNotFound, map[string]string{"Message": err.Error()})
		}
		return c.JSON(http.StatusOK, p.ToMap())
	})

	//PATCH :name proxy
	newServer.httpserver.PATCH("/proxies/:name", func(c echo.Context) error {
		if c.Request().Header.Get("Content-Type") != "application/json" {
			return c.JSON(http.StatusUnsupportedMediaType, map[string]string{"Message": "Must be 'application/json' MIME type"})
		}
		bodyDecoder := json.NewDecoder(c.Request().Body)
		p := new(ProxiesPostBody)
		if err := bodyDecoder.Decode(p); err != nil {
			return err
		}
		if err := b.proxylist.Set(c.Param("name"), p.Protocol, p.Address); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"Message": err.Error()})
		}
		c.Response().WriteHeader(http.StatusOK)
		return nil
	})

	//Delete /proxies/:name
	newServer.httpserver.DELETE("/proxies/:name", func(c echo.Context) error {
		if err := b.proxylist.Del(c.Param("name")); err != nil {
			return c.JSON(http.StatusNotFound, map[string]string{"Message": err.Error()})
		}
		c.Response().WriteHeader(http.StatusOK)
		return nil
	})

	//GET whitelist
	newServer.httpserver.GET("/whitelist", func(c echo.Context) error {
		c.Response().Header().Set("Content-Type", "application/json")
		return c.String(http.StatusOK, b.whitelist.ToJsArray())
	})

	//Query if a url exists in the whitelist
	newServer.httpserver.GET("/whitelist/:url", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]bool{"Exists": b.whitelist.Has(c.Param("url"))})
	})

	newServer.httpserver.POST("/proxies", func(c echo.Context) error {
		if c.Request().Header.Get("Content-Type") != "application/json" {
			return c.JSON(http.StatusUnsupportedMediaType, map[string]string{"Message": "Must be 'application/json' MIME type"})
		}
		bodyDecoder := json.NewDecoder(c.Request().Body)
		newProxy := new(ProxiesPostBody)
		if err := bodyDecoder.Decode(newProxy); err != nil {
			return err
		}
		if p, _ := b.proxylist.Find(newProxy.Name); p != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"Message": "Already exists"})
		}
		b.proxylist.Add(newProxy.Name, newProxy.Protocol, newProxy.Address)
		c.Response().WriteHeader(http.StatusCreated)
		return nil
	})

	newServer.httpserver.PATCH("/proxies/current", func(c echo.Context) error {
		if c.Request().Header.Get("Content-Type") != "application/json" {
			c.JSON(http.StatusUnsupportedMediaType, map[string]string{"Message": "Must be 'application/json' MIME type"})
			return nil
		}
		bodyDecoder := json.NewDecoder(c.Request().Body)
		switchName := new(ProxiesCurrentPatchBody)
		if err := bodyDecoder.Decode(switchName); err != nil {
			return err
		}
		if err := b.proxylist.SetCurrent(switchName.Name); err != nil {
			return c.JSON(http.StatusNotFound, map[string]string{"Message": err.Error()})
		} else {
			c.Response().WriteHeader(http.StatusOK)
		}
		return nil
	})

	newServer.httpserver.PATCH("/whitelist", func(c echo.Context) error {
		if c.Request().Header.Get("Content-Type") != "application/json" {
			return c.JSON(http.StatusUnsupportedMediaType, map[string]string{"Message": "Must be 'application/json' MIME type"})
		}
		bodyDecoder := json.NewDecoder(c.Request().Body)
		listOperation := new(WhitelistPatchBody)
		if err := bodyDecoder.Decode(listOperation); err != nil {
			return err
		}
		var err error
		switch listOperation.Operation {
		case "Add":
			if exists := b.whitelist.Has(listOperation.Url); exists {
				err = c.JSON(http.StatusForbidden, map[string]string{"Message": "Already exists"})
			} else {
				b.whitelist.Add(listOperation.Url)
				c.Response().WriteHeader(http.StatusOK)
			}
		case "Delete":
			if exists := b.whitelist.Has(listOperation.Url); exists {
				b.whitelist.Del(listOperation.Url)
				c.Response().WriteHeader(http.StatusOK)
			} else {
				err = c.JSON(http.StatusForbidden, map[string]string{"Message": "Doesn't exists"})
			}
		default:
			err = c.JSON(http.StatusBadRequest, map[string]string{"Message": "No such operation"})
		}
		return err
	})

	return newServer
}
