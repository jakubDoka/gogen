package dirs

import (
	"bufio"
	"gogen/str"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

// Gopath is value of GOPATH environment variable
var Gopath string

func init() {
	env, ok := os.LookupEnv("GOPATH")
	if !ok {
		panic("missing GOPATH")
	}
	Gopath = env
}

// PackPath returns go package path or error if it does not exist
func PackPath(imp string) (string, bool) {
	res := path.Join(Gopath, "src", imp)
	if Exists(res) {
		return res, true
	}

	return "", false
}

// PackImport return package import from path
func PackImport(p string) string {
	ln := len(path.Join(Gopath, "src")) + 1
	return p[ln:]
}

// PathName returns name of a base of a path without extention
func PathName(p string) string {
	p = NormPath(p)

	start, end := str.LastByte(p, '\\')+1, str.LastByte(p, '.')
	if end == -1 {
		return p[start:]
	}

	return p[start:end]
}

// NormPath returns path witch has all / replaced with \
func NormPath(p string) string {
	return strings.ReplaceAll(p, "/", "\\")
}

// ListFilePaths returns all paths to files in one directory.
// There is no recursion.
func ListFilePaths(p, ext string) (ps []string, err error) {
	return ListPaths(p, func(i os.FileInfo) bool {
		return !i.IsDir() && str.EndsWith(i.Name(), ext)
	})
}

// ListDirPaths returns all paths to dirs in one directory.
// There is no recursion.
func ListDirPaths(p string) (ps []string, err error) {
	return ListPaths(p, func(i os.FileInfo) bool {
		return i.IsDir()
	})
}

// ListPaths returns all paths to all items in directory, filtered
// There is no recursion.
func ListPaths(p string, filter func(os.FileInfo) bool) (ps []string, err error) {
	infos, err := ioutil.ReadDir(p)

	if err != nil {
		return
	}

	for _, i := range infos {
		if filter(i) {
			ps = append(ps, path.Join(p, i.Name()))
		}
	}

	return
}

// FileAsLines returns file as lines and excludes line s above and
// line with package and also returns name of package
func FileAsLines(p string) (lines Paragraph, name string, err error) {
	file, err := os.Open(p)
	if err != nil {
		return
	}
	defer file.Close()

	var afterPackage bool
	var i int

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		i++

		txt := scanner.Text()
		if afterPackage {
			lines = append(lines, Line{&p, i, txt})
		} else if str.StartsWith(txt, "package") {
			txt = str.RemInv(txt)
			name = txt[len("package"):]
			afterPackage = true
		}
	}

	err = scanner.Err()
	return
}

// Paragraph is group of lines
type Paragraph []Line

// NParagraph is only for debuging purposes
// it creates dummy paragraph from strings
func NParagraph(lines ...string) Paragraph {
	p := make(Paragraph, len(lines))
	for i := range lines {
		p[i] = Line{nil, i, lines[i]}
	}
	return p
}

// GetContent extracts all text to slice of strings
func (p Paragraph) GetContent() []string {
	res := make([]string, len(p))
	for i := range p {
		res[i] = p[i].Content
	}
	return res
}

// Copy ...
func (p Paragraph) Copy() Paragraph {
	np := make(Paragraph, len(p))
	copy(np, p)
	return np
}

// Line is file line, it stores its index and path for easy logging
type Line struct {
	Path    *string
	Idx     int
	Content string
}

func (l *Line) String() string {
	return l.Content
}

// CreateFile with initial content
func CreateFile(p, content string) error {
	f, err := os.Create(p)
	if err != nil {
		return err
	}

	defer f.Close()

	_, err = f.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}

// DeleteIfPresent deletes file or directory if it is present
func DeleteIfPresent(p string) error {
	if Exists(p) {
		return os.Remove(p)
	}
	return nil
}

// Exists returns whether file exist
func Exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
