package http

import (
	"flag"
	"github.com/Zumium/mywl/common"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
	"sync"
	"text/template"
)

const pacTemplate = `
var PROXY_METHOD = {{.Proxymethod}};
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

//type currentSetting struct {
//	Proxymethod string
//	Liststring  string
//}

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
		//err := t.Execute(c.Response(), &currentSetting{b.proxylist.GetCurrent().ToProxyMethodString(), b.whitelist.ToJsArray()})
		err := t.Execute(c.Response(), map[string]string{"Proxymethod": b.proxylist.GetCurrent().ToProxyMethodString(), "Liststring": b.whitelist.ToJsArray()})
		if err != nil {
			return err
		}
		return nil
	})

	//GET all proxy settings
	newServer.httpserver.GET("/proxies", func(c echo.Context) error {
		proxiesArray := make([]map[string]string, 0, b.proxylist.Len())
		b.proxylist.ForEach(func(each common.Proxy) {
			proxiesArray = append(proxiesArray, each.ToMap())
		})
		if err := c.JSON(http.StatusOK, proxiesArray); err != nil {
			return err
		}
		return nil
	})

	//GET current proxy setting
	newServer.httpserver.GET("/proxies/current", func(c echo.Context) error {
		if err := c.JSON(http.StatusOK, b.proxylist.GetCurrent().ToMap()); err != nil {
			return err
		}
		return nil
	})

	//GET whitelist
	newServer.httpserver.GET("/whitelist", func(c echo.Context) error {
		c.Response().Header().Set("Content-Type", "application/json")
		if err := c.String(http.StatusOK, b.whitelist.ToJsArray()); err != nil {
			return err
		}
		return nil
	})

	return newServer
}
