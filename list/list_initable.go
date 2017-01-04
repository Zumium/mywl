package list

import (
	"bufio"
	"container/list"
	"os"
)

func (this *WhiteList) Init() error {
	//WhiteList is Initable
	//Init the linkedlist
	this.whitelist = list.New()
	//load up list file
	return this.Load()
}

func (this *WhiteList) Load() error {
	this.whitelist.Init()
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

func (this *WhiteList) Save() error {
	var file *os.File
	defer func() {
		if file != nil {
			file.Close()
		}
	}()

	if f, err := os.Create(this.whitelistPath); err != nil {
		return err
	} else {
		file = f
	}

	writer := bufio.NewWriter(file)
	for e := this.whitelist.Front(); e != nil; e = e.Next() {
		s, _ := e.Value.(string)
		writer.WriteString(s)
		writer.WriteString("\n")
	}
	if err := writer.Flush(); err != nil {
		return err
	}

	return nil
}
