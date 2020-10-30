package commands

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

// string constats
const (
	TemplateAnnotation = "//gogen_template"
	IgnoreAnnotation   = "//gogen_ignore"
	GeneratedComment   = "//This file wos generated by gogen (golang code generator)" //TODO add repository
)

// PrintAllCommands prints all commands with additional information
func PrintAllCommands() {
	fmt.Println("So far i offer folloving...")
	for _, v := range CommandHandler.commands {
		fmt.Printf("\t%s - %s - %s\n", v.name, v.argsStruct, v.description)
	}
	Exit("i hope like it...UWU")
}

// PrintCommandHelp ...
func PrintCommandHelp(command Command) {
	fmt.Println("Some deeper descriptions...")
	for _, v := range command.help {
		fmt.Printf("\t%s %s\n", command.name, v)
	}
	Exit("If this no help, nothing any help.")
}

// CreateTemplatesInDir collects gen requests from all files and creates them in separate file
func CreateTemplatesInDir(dir string) {
	infos, err := ioutil.ReadDir(dir)
	if Warming("skipping directory, it cannot be red", err) {
		return
	}

	requests := []Rules{}
	var pack string
	for _, inf := range infos {

		path := path.Join(dir, inf.Name())
		if !IsGoFile(path) {
			continue
		}

		requests = append(requests, CollectFileTemplateRequests(path, &pack)...)
	}

	if pack == "" {
		Terminate("unable to resolve package name in directory: " + dir)
	}

	genFilePath := path.Join(dir, Cf.GeneratedFile)
	exists := Exists(genFilePath)
	if len(requests) == 0 {
		if exists {
			err := os.Remove(genFilePath)
			Warming("unable to remove redundant template file", err)
		}
		return
	} else if !exists {
		_, err := os.Create(genFilePath)
		CheckError("unable to create template file", err)
	}

	requests = FilterDuplipcates(requests)
	for _, r := range requests {
		HandleSyntaxError(r.ErrData)
	}
	code := ""
	imports := map[string]bool{}
	for _, r := range requests {
		temp, _ := Templates[r.name]
		MergeImports(imports, temp.imports)
		code += temp.Generate(r.args)
	}

	err = ioutil.WriteFile(genFilePath, []byte(pack+FormatImports(imports)+code), 0644)
	Warming("unable to modify template file", err)
	return
}

// CollectFileTemplateRequests ...
func CollectFileTemplateRequests(path string, pack *string) (requests []Rules) {
	file, err := os.Open(path)
	if Warming("skipping file, it cannot be opened", err) {
		return
	}

	defer file.Close()
	scanner := bufio.NewScanner(file)
	lineIdx := 0
	for scanner.Scan() {
		lineIdx++
		line := scanner.Text()
		if !strings.Contains(line, Cf.TemplateRequest) {
			if *pack == "" && StartsWith(Bulkremove("\n\r\t", line), "package") {
				fmt.Println(line)
				*pack = line + "\n"
			}
			continue
		}
		start, end := len(Cf.TemplateRequest), strings.Index(line, ">")+1
		if end == 0 {
			continue
		}
		req := NewRules(line[start:end], ErrData{path, lineIdx, -1})

		if req.err == -1 {
			if val, ok := Templates[req.name]; ok {
				if len(val.args) != len(req.args) {
					req.err = IncorrectArgs
				}
			} else {
				req.err = Undefined
			}
		}

		requests = append(requests, req)

	}
	fmt.Println(requests)
	return
}

// MergeImports appends one map of imports to another
func MergeImports(a, b map[string]bool) {
	for k := range b {
		a[k] = true
	}
}

// FormatImports turns imports to correct go sintax
func FormatImports(a map[string]bool) string {
	if len(a) == 0 {
		return ""
	}
	str := "\nimport (\n"
	for k := range a {
		str += "\t" + k + "\n"
	}
	return str + ")"
}

// Bulkremove removes all characters contained in string from original string
func Bulkremove(chars string, orig string) string {
	for _, r := range chars {
		orig = strings.Replace(orig, string(r), "", -1)
	}
	return orig
}

// ExtractImports parses string to slice of package names
func ExtractImports(str string) (imports []string) {
	parts := strings.Split(str, "import")

	if len(parts) < 2 {

		return
	}

	for _, i := range parts[1:] {
		start, end := strings.IndexRune(i, '(')+1, strings.IndexRune(i, ')')
		var str []string
		if end == -1 || start > strings.IndexRune(i, '\n') {
			str = []string{i[1:strings.IndexRune(i, '\n')]}
		} else {
			str = strings.Split(i[start:end], "\n")
		}

		for _, v := range str {
			v = Bulkremove(" \r\t\n", v)
			if len(v) == 0 {
				continue
			}
			imports = append(imports, v)
		}

	}
	fmt.Printf("%q", imports)
	return
}

// StartsWith for string
func StartsWith(str, sub string) bool {
	return len(str) >= len(sub) && str[:len(sub)] == sub
}

// ParseTemplatesInDir parses all files that are annotated in given directory, recursively
func ParseTemplatesInDir(dir string) (res []Template) {
	infos, err := ioutil.ReadDir(dir)
	if Warming("skipping directory, it cannot be read", err) {
		return
	}

	for _, info := range infos {
		if err != nil || info.IsDir() || !IsGoFile(info.Name()) {
			continue
		}

		res = append(res, ParseTemplateFile(path.Join(dir, info.Name()))...)
	}
	return
}

// ParseTemplateFile takes all templates from file and saves them to heap
func ParseTemplateFile(path string) (res []Template) {
	bytes, err := ioutil.ReadFile(path)
	if Warming("skipping file, it cannot be opened, nor read", err) {
		return
	}

	content := string(bytes)
	if !StartsWith(content, TemplateAnnotation) {
		if !StartsWith(content, IgnoreAnnotation) {
			Warming("skipping file, you can annotate file with '"+IgnoreAnnotation+"' if its intentional",
				errors.New("file is not annotated"))
		}
		return
	}

	templates := strings.Split(content, Cf.TemplateStart)
	if len(templates) == 0 {
		return
	}

	imports := ExtractImports(templates[0])

	for _, t := range templates[1:] {
		idx := strings.Index(t, ">") + 1
		key := t[:idx]
		start, end := strings.Index(t, "\n")+1, strings.Index(t, Cf.TemplateEnd)
		if end == -1 {
			end = len(t)
		}
		temp := NewTemplate(key, t[start:end], imports, ErrData{path, GetErrorLine(path, key, content), -1})
		res = append(res, temp)
	}
	return
}
