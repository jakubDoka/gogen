package dirs

import (
	"bufio"
	"errors"
	"gogen/str"
	"io/ioutil"
	"os"
	"path"
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
func PackPath(imp string) (res string, err error) {
	res = path.Join(Gopath, "src", imp)
	if Exists(res) {
		return
	}

	return "", errors.New("package does not exist")
}

// ListFilePaths returns all paths to files in one directory.
// There is no recursion.
func ListFilePaths(p string) (ps []string, err error) {
	return ListPaths(p, func(i os.FileInfo) bool {
		return !i.IsDir()
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
func FileAsLines(p string) (lines []string, name string, err error) {
	file, err := os.Open(p)
	if err != nil {
		return
	}
	defer file.Close()

	var afterPackage bool

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		txt := scanner.Text()
		if afterPackage {
			lines = append(lines, txt)
		} else if str.StartsWith(txt, "package") {
			txt = str.RemInv(txt)
			name = txt[len("package"):]
			afterPackage = true
		}

	}

	err = scanner.Err()
	return
}

// Exists returns whether file exist
func Exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
