package config

import (
	"github.com/BurntSushi/toml"
	"os"
)

// TOML implementation of the Loader interface
type TomlLoader struct {
	Filename string
}

// Blocking method which returns a pointer to the loaded configuration and a non-wrapped error if something goes wrong.
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

// Blocking method which accepts a pointer to the loaded configuration and a non-wrapped error if something goes wrong. It saves the configuration values to the file name provided by the TomlLoader.
func (loader *TomlLoader) Save(configuration *Configuration) (err error) {
	file, err := os.Create(loader.Filename)
	if err != nil && !os.IsExist(err) {
		return err
	}
	defer file.Close()
	return toml.NewEncoder(file).Encode(configuration)
}
