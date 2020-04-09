package envconfig

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"github.com/short-d/app/fw"
)

type EnvConfig struct {
	environment fw.Environment
}

func (e EnvConfig) ParseConfigFromEnv(config interface{}) error {
	configVal := reflect.ValueOf(config)
	if configVal.Kind() != reflect.Ptr || configVal.IsNil() {
		return errors.New("config must be a pointer")
	}

	elem := configVal.Elem()
	if elem.Kind() != reflect.Struct {
		return errors.New("config must be a struct")
	}

	numFields := elem.NumField()
	configType := elem.Type()

	for idx := 0; idx < numFields; idx++ {
		field := configType.Field(idx)
		envName, ok := field.Tag.Lookup("env")
		if !ok {
			continue
		}
		defaultVal := field.Tag.Get("default")
		envVal := e.environment.GetEnv(envName, defaultVal)
		err := setFieldValue(field, elem.Field(idx), envVal)
		if err != nil {
			return err
		}
	}
	return nil
}

func setFieldValue(field reflect.StructField, fieldValue reflect.Value, newValue string) error {
	kind := field.Type.Kind()
	switch kind {
	case reflect.String:
		fieldValue.SetString(newValue)
		return nil
	case reflect.Int:
		num, err := strconv.Atoi(newValue)
		if err != nil {
			return err
		}
		fieldValue.SetInt(int64(num))
		return nil
	case reflect.Bool:
		boolean, err := strconv.ParseBool(newValue)
		if err != nil {
			return err
		}
		fieldValue.SetBool(boolean)
		return nil
	default:
		return fmt.Errorf("unexpected field type: %s", kind)
	}
}

func NewEnvConfig(environment fw.Environment) EnvConfig {
	return EnvConfig{environment: environment}
}
