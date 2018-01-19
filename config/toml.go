package config

import (
	"github.com/BurntSushi/toml"
	"os"
)

type TomlLoader struct {
	Filename string
}

func (loader *TomlLoader) Load() (configuration *Configuration, err error) {
	var file *os.File
	file, err = os.Open(loader.Filename)
	if err != nil {
		return
	}
	defer file.Close()
	configuration = &Configuration{}
	_, err = toml.DecodeReader(file, configuration)
	return
}

func (loader *TomlLoader) Save(configuration *Configuration) (err error) {
	file, err := os.Create(loader.Filename)
	if err != nil && !os.IsExist(err) {
		return err
	}
	defer file.Close()
	return toml.NewEncoder(file).Encode(configuration)
}
