package main

import (
	"github.com/dshayles/hclparse"
)

func main() {
	p := hclparse.NewParser("terraform_files.txt")
	p.FindHclObjects()

	p.PrintCsv("test.csv")
}
