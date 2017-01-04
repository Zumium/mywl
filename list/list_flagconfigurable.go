package list

import "flag"

func (this *WhiteList) InstallFlags(flagset *flag.FlagSet) {
	//WhiteList is FlagConfigurable
	flagset.StringVar(&this.whitelistPath, "listfile", "/etc/mywl/whitelist.txt", "The Whitelist File Path")
}
