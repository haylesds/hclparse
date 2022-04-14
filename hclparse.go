package hclparse

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func ImportFromList(list string) *MyFileList {
	var fs []*MyFile

	file, _ := os.Open(list)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fs = append(fs, &MyFile{
			Path:         scanner.Text(),
			ContentBytes: nil,
			Content: &MyContent{
				Lines: make(map[int]string, 0),
			},
			Modules:   make(map[int]*HCLModule),
			Data:      make(map[int]*HCLData),
			Variables: make(map[int]*HCLVariable),
			Providers: make(map[int]*HCLProvider),
			Backend:   make(map[int]*HCLBackend),
			Resource:  make(map[int]*HCLResource),
			Locals:    make(map[int]*HCLLocals),
		})
	}

	return &MyFileList{Files: fs}
}

// ### MyContent

type MyContent struct {
	Lines map[int]string
}

// ### MyFile

type MyFile struct {
	Path         string
	ContentBytes []byte     `json:"-"`
	Content      *MyContent `json:"-"`
	Modules      map[int]*HCLModule
	Variables    map[int]*HCLVariable
	Providers    map[int]*HCLProvider
	Data         map[int]*HCLData
	Resource     map[int]*HCLResource
	Backend      map[int]*HCLBackend
	Locals       map[int]*HCLLocals
}

func (mf *MyFile) GetContent() {
	content, _ := os.ReadFile(mf.Path)
	mf.ContentBytes = content
}

func (mf *MyFile) ParseContent() {

	file := bytes.NewReader(mf.ContentBytes)
	scanner := bufio.NewScanner(file)

	i := 0

	for scanner.Scan() {
		i++
		mf.Content.Lines[i] = scanner.Text()
	}
}

func (mf *MyFile) PrintContent() {
	fmt.Println(string(mf.ContentBytes))
}

func (mf *MyFile) PrintContentLines() {

	l := len(mf.Content.Lines)
	i := 1

	for i <= l {
		toWrite := fmt.Sprintf("%d: %s", i, mf.Content.Lines[i])

		if i < 10 && l >= 10 {
			toWrite = fmt.Sprintf("0%s", toWrite)
		}

		if i < 100 && l >= 100 {
			toWrite = fmt.Sprintf("0%s", toWrite)
		}

		if i < 1000 && l >= 1000 {
			toWrite = fmt.Sprintf("0%s", toWrite)
		}

		fmt.Println(toWrite)

		i++
	}

	fmt.Print("\n")
}

func (mf *MyFile) FindHclObjects() {

	l := len(mf.Content.Lines)
	i := 1

	for i <= l {
		line := mf.Content.Lines[i]

		if IsModule(line) {
			x := NewHCLModule()
			x.StartLine = i
			x.EndLine = FindClose(mf.Content.Lines, i)
			x.Label = FindQuote(line)

			j := i + 1

			for j <= x.EndLine {
				if IsProperty(mf.Content.Lines[j]) {
					k, v := HandleProperty(mf.Content.Lines[j])
					if v == "[" {
						v = "list"
					}
					x.Properties[k] = v
				}
				j++
			}

			if val, ok := x.Properties["version"]; ok {
				x.Version = strings.Replace(val, "\"", "", -1)
			}

			if val, ok := x.Properties["source"]; ok {
				x.Source = strings.Replace(val, "\"", "", -1)
			}

			mf.Modules[i] = x
		}

		if IsVariable(line) {
			x := NewHCLVariable()
			x.StartLine = i
			x.EndLine = FindClose(mf.Content.Lines, i)
			x.Label = FindQuote(line)

			j := i + 1

			for j < x.EndLine {
				x.Values = append(x.Values, mf.Content.Lines[j])
				j++
			}

			mf.Variables[i] = x
		}

		if IsProvider(line) {
			x := NewHCLProvider()
			x.StartLine = i
			x.EndLine = FindClose(mf.Content.Lines, i)
			x.Label = FindQuote(line)

			j := i + 1

			for j <= x.EndLine {
				if IsProperty(mf.Content.Lines[j]) {
					k, v := HandleProperty(mf.Content.Lines[j])
					if v == "[" {
						v = "list"
					}
					x.Properties[k] = v
				}
				j++
			}

			mf.Providers[i] = x
		}

		if IsData(line) {
			x := NewHCLData()
			x.StartLine = i
			x.EndLine = FindClose(mf.Content.Lines, i)
			x.Type = FindQuote(line)
			x.Label = FindQuote(strings.Replace(line, "\"x.Type\"", "", -1))

			j := i + 1

			for j <= x.EndLine {
				if IsProperty(mf.Content.Lines[j]) {
					k, v := HandleProperty(mf.Content.Lines[j])
					if v == "[" {
						v = "list"
					}
					x.Properties[k] = v
				}
				j++
			}

			mf.Data[i] = x
		}

		if IsResource(line) {
			x := NewHCLResource()
			x.StartLine = i
			x.EndLine = FindClose(mf.Content.Lines, i)
			x.Type = FindQuote(line)
			x.Label = FindSecondQuote(line)

			j := i + 1

			for j < x.EndLine {
				x.Values = append(x.Values, mf.Content.Lines[j])
				j++
			}

			mf.Resource[i] = x
		}

		if IsBackend(line) {
			x := NewHCLBackend()
			x.StartLine = i
			x.EndLine = FindClose(mf.Content.Lines, i)
			x.Label = FindQuote(line)

			j := i + 1

			for j <= x.EndLine {
				if IsProperty(mf.Content.Lines[j]) {
					k, v := HandleProperty(mf.Content.Lines[j])
					if v == "[" {
						v = "list"
					}
					x.Properties[k] = v
				}
				j++
			}

			mf.Backend[i] = x
		}

		if IsLocals(line) {
			x := NewHCLLocals()
			x.StartLine = i
			x.EndLine = FindClose(mf.Content.Lines, i)

			j := i + 1

			for j <= x.EndLine {
				if IsProperty(mf.Content.Lines[j]) {
					k, v := HandleProperty(mf.Content.Lines[j])
					if v == "[" {
						v = "list"
					}
					x.Properties[k] = v
				}
				j++
			}

			mf.Locals[i] = x
		}

		i++
	}

	//fmt.Print("\n")
}

// #### MyFileList

type MyFileList struct {
	Files []*MyFile
}

func (mfl *MyFileList) ToString() {
	for _, f := range mfl.Files {
		fmt.Println(f.Path)
	}
}

// ### MyParser

type MyParser struct {
	FileList *MyFileList
}

func NewParser(list string) *MyParser {
	p := &MyParser{
		FileList: ImportFromList(list),
	}

	for _, f := range p.FileList.Files {
		f.GetContent()
		f.ParseContent()
	}

	return p
}

func (p *MyParser) FindHclObjects() {
	for _, f := range p.FileList.Files {
		f.FindHclObjects()
	}
}

func (p *MyParser) FilesWithModules() []*MyFile {

	var output []*MyFile

	for _, f := range p.FileList.Files {
		if len(f.Modules) > 0 {
			output = append(output, f)
		}
	}

	return output
}

func (p *MyParser) PrintCsv(path string) error {

	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	defer f.Close()

	headers := []string{
		"file",
		"modules",
		"variables",
		"data",
		"resources",
		"backend",
		"providers",
		"locals",
	}

	output := [][]string{
		headers,
	}

	for _, f := range p.FileList.Files {
		modules := ""
		variables := ""
		data := ""
		resources := ""
		backend := ""
		providers := ""
		locals := ""

		for _, m := range f.Modules {

			v := ""

			if m.Version == "" {
				v = "?"
			} else {
				v = m.Version
			}

			modules = fmt.Sprintf("%s\n%s - %s (v%s)", modules, m.Label, m.Source, v)
		}

		for _, v := range f.Variables {
			variables = fmt.Sprintf("%s\n%s", variables, v.Label)
		}

		for _, d := range f.Data {
			data = fmt.Sprintf("%s\n%s", data, d.Label)
		}

		for _, r := range f.Resource {
			resources = fmt.Sprintf("%s\n%s: %s", resources, r.Type, r.Label)
		}

		for _, b := range f.Backend {
			backend = fmt.Sprintf("%s\n%s", backend, b.Label)
		}

		for _, p := range f.Providers {
			providers = fmt.Sprintf("%s\n%s", providers, p.Label)
		}

		for _, l := range f.Locals {
			for k, v := range l.Properties {
				locals = fmt.Sprintf("%s\n%s = %s", locals, k, v)
			}
		}

		row := []string{
			f.Path,
			strings.TrimSpace(modules),
			strings.TrimSpace(variables),
			strings.TrimSpace(data),
			strings.TrimSpace(resources),
			strings.TrimSpace(backend),
			strings.TrimSpace(providers),
			strings.TrimSpace(locals),
		}

		output = append(output, row)

	}

	writer := csv.NewWriter(f)
	err = writer.WriteAll(output)
	if err != nil {
		return err
	}

	return nil
}

func (p *MyParser) OutputAllResources(path string) error {

	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	defer f.Close()

	output := []string{}

	for _, f := range p.FileList.Files {
		resources := []string{}

		for _, l := range f.Resource {
			resources = append(resources, fmt.Sprintf("resource \"%s\" \"%s\" {", l.Type, l.Label))
			for _, v := range l.Values {
				resources = append(resources, v)
			}
			resources = append(resources, "}")
		}

		output = append(output, resources...)

	}

	_, err = f.WriteString(Flatten(output))
	if err != nil {
		return err
	}

	return nil
}

// ### HCL Objects

type HCLModule struct {
	Label      string
	Source     string
	Version    string
	StartLine  int
	EndLine    int
	Properties map[string]string
}

type HCLVariable struct {
	Label     string
	Values    []string
	StartLine int
	EndLine   int
}

type HCLProvider struct {
	Label      string
	StartLine  int
	EndLine    int
	Properties map[string]string
}

type HCLData struct {
	Label      string
	Type       string
	StartLine  int
	EndLine    int
	Properties map[string]string
}

type HCLResource struct {
	Label     string
	Type      string
	StartLine int
	EndLine   int
	Values    []string
}

type HCLBackend struct {
	Label      string
	StartLine  int
	EndLine    int
	Properties map[string]string
}

type HCLLocals struct {
	StartLine  int
	EndLine    int
	Properties map[string]string
}

func NewHCLModule() *HCLModule {
	return &HCLModule{
		Properties: make(map[string]string),
	}
}

func NewHCLVariable() *HCLVariable {
	return &HCLVariable{
		Values: make([]string, 0),
	}
}

func NewHCLProvider() *HCLProvider {
	return &HCLProvider{
		Properties: make(map[string]string),
	}
}

func NewHCLData() *HCLData {
	return &HCLData{
		Properties: make(map[string]string),
	}
}

func NewHCLResource() *HCLResource {
	return &HCLResource{
		Values: make([]string, 0),
	}
}

func NewHCLBackend() *HCLBackend {
	return &HCLBackend{
		Properties: make(map[string]string),
	}
}

func NewHCLLocals() *HCLLocals {
	return &HCLLocals{
		Properties: make(map[string]string),
	}
}

// ### UTILS

func ToJson(i interface{}) []byte {
	j, _ := json.MarshalIndent(i, "", " ")
	return j
}

func FindQuote(s string) string {
	newString := ""
	first := strings.Index(s, "\"")
	newString = s[first+1:]
	last := strings.Index(newString, "\"")
	newString = newString[:last]
	return newString
}

func FindSecondQuote(s string) string {
	newString := ""

	first := strings.Index(s, "\"")
	newString = s[first+1:]

	second := strings.Index(newString, "\"")
	newString = newString[second+1:]

	third := strings.Index(newString, "\"")
	newString = newString[third+1:]

	fourth := strings.Index(newString, "\"")
	newString = newString[:fourth]

	return newString
}

func FindClose(lines map[int]string, start int) int {

	opens := CountOpens(lines[start])
	closes := CountCloses(lines[start])

	current := start

	if opens == 0 {
		fmt.Println("Can't find open.")
	} else if closes >= opens {
		return current
	} else {
		for closes < opens {
			current++
			opens += CountOpens(lines[current])
			closes += CountCloses(lines[current])
		}
		return current
	}

	return 0
}

func CountOpens(s string) int {
	return strings.Count(s, "{")
}

func CountCloses(s string) int {
	return strings.Count(s, "}")
}

func IsModule(s string) bool {
	output := false

	if strings.Contains(s, "module \"") && strings.Contains(s, "{") {
		output = true
	}

	return output
}

func IsVariable(s string) bool {
	output := false

	if strings.Contains(s, "variable \"") && strings.Contains(s, "{") {
		output = true
	}

	return output
}

func IsProvider(s string) bool {
	output := false

	if strings.Contains(s, "provider \"") && strings.Contains(s, "{") {
		output = true
	}

	return output
}

func IsData(s string) bool {
	output := false

	if strings.Contains(s, "data \"") && strings.Contains(s, "{") {
		output = true
	}

	return output
}

func IsResource(s string) bool {
	output := false

	if strings.Contains(s, "resource \"") && strings.Contains(s, "{") {
		output = true
	}

	return output
}

func IsBackend(s string) bool {
	output := false

	if strings.Contains(s, "backend \"") && strings.Contains(s, "{") {
		output = true
	}

	return output
}

func IsLocals(s string) bool {
	output := false

	if strings.Contains(s, "locals ") && strings.Contains(s, "{") {
		output = true
	}

	return output
}

func IsProperty(s string) bool {
	output := false

	if strings.Contains(s, " = ") {
		output = true
	}

	return output
}

func HandleProperty(s string) (key, value string) {
	ssplit := strings.Split(s, "=")

	key = strings.TrimSpace(ssplit[0])
	value = strings.TrimSpace(ssplit[1])

	return key, value
}

func Flatten(s []string) string {
	output := ""
	for _, v := range s {
		output = fmt.Sprintf("%s\n%s", output, v)
	}
	output = strings.TrimSpace(output)
	return output
}
