package list

import (
	"bytes"
	"fmt"
	"os"
)

func (this *WhiteList) Add(url string) {
	this.whitelist.PushFront(url)

	if err := this.Save(); err != nil {
		fmt.Fprintf(os.Stderr, "error occurd on whitelist's saving process: %s\n", err.Error())
	}
}

func (this *WhiteList) Del(url string) {
	for e := this.whitelist.Front(); e != nil; e = e.Next() {
		str, _ := e.Value.(string)
		if str == url {
			this.whitelist.Remove(e)
			if err := this.Save(); err != nil {
				fmt.Fprintf(os.Stderr, "error occurd on whitelist's saving process: %s\n", err.Error())
			}
			return
		}
	}
}

func (this *WhiteList) Has(url string) bool {
	for e := this.whitelist.Front(); e != nil; e = e.Next() {
		str, _ := e.Value.(string)
		if str == url {
			return true
		}
	}
	return false
}

func (this *WhiteList) ToJsArray() string {
	var buffer bytes.Buffer
	lastElem := this.whitelist.Back()
	buffer.WriteString("[\n")
	for e := this.whitelist.Front(); e != nil; e = e.Next() {
		str, _ := e.Value.(string)
		buffer.WriteString("\"")
		buffer.WriteString(str)
		buffer.WriteString("\"")
		if e != lastElem {
			buffer.WriteString(",\n")
		}
	}
	buffer.WriteString("\n]")
	return buffer.String()
}
