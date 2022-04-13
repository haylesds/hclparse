package main

import (
	"github.com/haylesds/hclparse"
)

func main() {
	p := hclparse.NewParser("terraform_files.txt")
	p.FindHclObjects()

	//p.PrintCsv("test.csv")
	p.OutputAllResources("test.tf")
}
