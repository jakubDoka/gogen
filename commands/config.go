package commands

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
)

// ConfigFile is name of config file
var ConfigFile = "config.json"

// Config stores all configuration
type Config struct {
	TemplateStart, TemplateEnd, TemplateRequest,
	PrefixSpecifier, GeneratedFile string
	TemplatePaths map[string]bool
}

// Cf stores configuration
var Cf Config = Config{
	TemplateRequest: "//gen",
	TemplateStart:   "//<<<",
	TemplateEnd:     "//>>>",
	PrefixSpecifier: "_prefix",
	GeneratedFile:   "gogen-output.go",
	TemplatePaths:   map[string]bool{},
}

// LoadConfig creates new if doesn't exist
func LoadConfig() {
	if Labels["-l"] {
		ConfigFile = path.Join(Dirs[0], "gogen-"+ConfigFile)
	}
	if !Exists(ConfigFile) {
		CheckError("unable to create default config", SaveConfig())
	} else {
		bytes, err := ioutil.ReadFile(ConfigFile)
		CheckError("cannot open config file", err)
		err = json.Unmarshal(bytes, &Cf)
		CheckError("cannot parse config, if you don't know how to"+
			"fix it, delete file and let default one be generated", err)
	}
}

// SaveConfig ...
func SaveConfig() error {
	f, err := os.Create(ConfigFile)

	defer f.Close()

	bytes, err := json.Marshal(&Cf)

	if err != nil {
		return err
	}

	_, err = f.Write(bytes)

	if err != nil {
		return err
	}
	return nil
}
