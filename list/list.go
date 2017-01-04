package list

import "container/list"

type WhiteList struct {
	whitelistPath string
	whitelist     *list.List
}
