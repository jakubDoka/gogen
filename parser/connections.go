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

// LoadConnections ...
func LoadConnections() Connections {
	c := Connections{}
	p := path.Join(dirs.Gopath, DataDir, ConnectionFile)
	if !dirs.Exists(p) {
		return c
	}
	bts, err := ioutil.ReadFile(p)
	if err != nil {
		fmt.Println("failed to load connections:", err)
		return c
	}
	err = json.Unmarshal(bts, &c)
	if err != nil {
		fmt.Println("failed to parse connections:", err)
	}

	return c
}

// Connections stores all connections for bulk regenerations
type Connections map[string]map[string]bool

// Get makes sure key always returns map and not nil, even if its empty map
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
