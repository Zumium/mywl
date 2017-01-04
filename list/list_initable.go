package list

import (
	"bufio"
	"container/list"
	"os"
)

func (this *WhiteList) Init() error {
	//WhiteList is Initable
	//convert path to absolute file path
	//absPath, _ := filepath.Abs(this.whitelistPath)
	//this.whitelistPath = absPath
	//Init the linkedlist
	this.whitelist = list.New()
	//load up list file
	if err := this.loadList(); err != nil {
		return err
	}
	return nil
}

func (this *WhiteList) loadList() error {
	//load list items from file to 'whitelist' field
	//step 1: open the given file
	var whitelistFile *os.File
	defer func() {
		if whitelistFile != nil {
			whitelistFile.Close()
		}
	}()
	if file, err := os.Open(this.whitelistPath); err != nil {
		whitelistFile = nil
		return err
	} else {
		whitelistFile = file
	}
	//step 2: Create line reader
	scanner := bufio.NewScanner(whitelistFile)
	//step 3: Start reading
	for scanner.Scan() {
		this.whitelist.PushBack(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	//Reading completed
	return nil
}
