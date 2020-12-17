package parser

import (
	"encoding/json"
	"fmt"
	"gogen/dirs"
	"io/ioutil"
	"os"
	"path"
)

// ConnectionFile is name of a file where connection data is stored
var ConnectionFile = "connections.json"

// DataDir stores all data related to gogen
var DataDir = "gogen-data"

// Cons ...
var Cons = Connections{}

// LoadConnections ...
func LoadConnections() {
	p := path.Join(dirs.Gopath, DataDir, ConnectionFile)
	if !dirs.Exists(p) {
		return
	}
	bts, err := ioutil.ReadFile(p)
	if err != nil {
		fmt.Println("failed to load connections:", err)
		return
	}
	err = json.Unmarshal(bts, &Cons)
	if err != nil {
		fmt.Println("failed to parse connections:", err)
	}
}

// Connections stores all connections for bulk regenerations
type Connections map[string]map[string]bool

func (c Connections) Get(pack string) map[string]bool {
	val, ok := c[pack]
	if !ok {
		c[pack] = map[string]bool{}
		return c[pack]
	}

	return val
}

// Save saves Connections
func (c *Connections) Save() {
	bts, err := json.Marshal(c)
	if err != nil {
		panic(err)
	}
	p := path.Join(dirs.Gopath, DataDir)
	err = os.MkdirAll(p, os.ModeDir)
	if err != nil {
		panic(err)
	}
	f, err := os.Create(path.Join(p, ConnectionFile))
	if err != nil {
		panic(err)
	}
	defer f.Close()
	_, err = f.Write(bts)
	if err != nil {
		panic(err)
	}
}
