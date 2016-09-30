package config

import (
	"os"
	"path"
	"fmt"
	"strings"
	"encoding/json"
	"reflect"
	"io/ioutil"
	"errors"
)

//config 类
type Config struct {

}

func NewConfig() *Config {
	return &Config{}
}

//获取开发环境
func (this *Config) Env() string {
	if env := os.Getenv("GO_ENV"); env != "" {
		return env 
	}

	return "development"
}

func (this *Config) GetConfigByEnv(file, env string) (string, error) {
	var envFile string 
	var extname string
	extname = path.Ext(file)

	if extname == "" {
		envFile = fmt.Sprintf("%v.%v", file, env)
	} else {
		envFile = fmt.Sprintf("%v.%v%v", strings.TrimSuffix(file, extname), env, extname)
	}

	if fileInfo, err := os.Stat(envFile); err == nil && fileInfo.Mode().IsRegular() {
		return envFile, nil
	} else {
		return "", fmt.Errorf("fail to find file %v", file)
	}

}

func (this *Config) GetConfigurations(files ...string) []string {
	var results []string
	env := this.Env()
	for i := len(files) - 1; i >= 0; i-- {
		var foundFile bool
		var file = files[i]
		if fileInfo, err := os.Stat(file); err == nil && fileInfo.Mode().IsRegular() {
			foundFile = true
			results = append(results, file)
		}

		if file, err := this.GetConfigByEnv(file, env); err == nil {
			foundFile = true
			results = append(results, file)
		} 

		if !foundFile {
			if example, err := this.GetConfigByEnv(file, "example"); err == nil {
				fmt.Printf("Failed to find configuration %v, using example file %v\n", file, example)
				results = append(results, example)
			} else {
				fmt.Printf("Failed to find configuration %v\n", file)
			}
		}
	}

	return results
}

func (this *Config) GetPrefix(config interface{}) string {
	if prefix := os.Getenv("CONFIGOR_ENV_PRE"); prefix != "" {
		return prefix
	}

	return "configor"
}

func (this *Config) Load(config interface{}, files ...string) error {
	for _, file := range this.GetConfigurations(files...) {
		if err := this.load(config, file); err != nil {
			return err
		}
	}

	if prefix := this.GetPrefix(config); prefix == "-" {
		return this.processTags(config)
	} else {
		return this.processTags(config, prefix)
	}
}

func (this *Config) load(config interface{}, file string) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, config)	
}

func (this *Config) processTags(config interface{}, prefix ...string) error {
	configValue := reflect.Indirect(reflect.ValueOf(config))
	if configValue.Kind() != reflect.Struct {
		return errors.New("invalid config, should be struct")
	}
	configType := configValue.Type()
	for i := 0; i < configType.NumField(); i++ {
		fieldStruct := configType.Field(i)
		field := configValue.Field(i)
		if field.Kind() == reflect.Struct {
			if err := this.processTags(field.Addr().Interface(), append(prefix, fieldStruct.Name)...); err != nil {
				return err
			}
		}
	}

	return nil 
}	

















