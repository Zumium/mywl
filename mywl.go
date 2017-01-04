package main

import (
	"flag"
	"fmt"
	"github.com/Zumium/mywl/common"
	"github.com/Zumium/mywl/http"
	"github.com/Zumium/mywl/list"
	"github.com/Zumium/mywl/proxy"
	"os"
)

func installFlags(flagConfigurables []common.FlagConfigurable, flagset *flag.FlagSet) {
	for _, mod := range flagConfigurables {
		mod.InstallFlags(flagset)
	}
}

func initModules(initables []common.Initable) {
	for _, mod := range initables {
		if err := mod.Init(); err != nil {
			fmt.Fprintf(os.Stderr, "error occurd on initing process: %s", err.Error())
			os.Exit(-1)
		}
	}
}

func save(persistables []common.Persistable) {
	for _, mod := range persistables {
		if err := mod.Save(); err != nil {
			fmt.Fprintf(os.Stderr, "error occurd on exiting process: %s", err.Error())
			os.Exit(-1)
		}
	}
}

func main() {
	//Parse Commandline Args
	flagset := flag.NewFlagSet("mywl", flag.ExitOnError)
	installFlags([]common.FlagConfigurable{list.GetInstance(), proxy.GetProxyList(), http.GetBuilder()}, flagset)
	flagset.Parse(os.Args[1:])
	//Init modules
	initModules([]common.Initable{list.GetInstance(), proxy.GetProxyList()})
	//Assemble http server
	server := http.GetBuilder().SetWhiteList(list.GetInstance()).SetProxyList(proxy.GetProxyList()).Build()
	//start server
	server.Start()

	//Going to exit
	save([]common.Persistable{list.GetInstance(), proxy.GetProxyList()})
	fmt.Println("Shutting down...")
	os.Exit(0)
}
