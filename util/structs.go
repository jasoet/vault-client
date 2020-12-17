package util

import (
	"github.com/mitchellh/mapstructure"
	"reflect"
	"time"
)

/*
This function will help you to convert your object from struct to map[string]interface{} based on your JSON tag in your structs.
Example how to use posted in sample_test.go file.
Credit: https://gist.github.com/bxcodec/c2a25cfc75f6b21a0492951706bc80b8
*/
func StructToMap(item interface{}) map[string]interface{} {

	res := map[string]interface{}{}
	if item == nil {
		return res
	}
	v := reflect.TypeOf(item)
	reflectValue := reflect.ValueOf(item)
	reflectValue = reflect.Indirect(reflectValue)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	for i := 0; i < v.NumField(); i++ {
		tag := v.Field(i).Tag.Get("json")
		field := reflectValue.Field(i).Interface()
		if tag != "" && tag != "-" {
			if v.Field(i).Type.Kind() == reflect.Struct {
				res[tag] = StructToMap(field)
			} else {
				res[tag] = field
			}
		}
	}
	return res
}

//MapToStruct used to convert Map to Struct, mapping uses `json` tag, will also decode string to time with `time.RFC3339Nano` layout
//See https://github.com/mitchellh/mapstructure/blob/master/mapstructure_test.go for mapstructure library usage example.
func MapToStruct(input interface{}, result interface{}) (err error) {
	config := &mapstructure.DecoderConfig{TagName: "json", Result: result, DecodeHook: mapstructure.StringToTimeHookFunc(time.RFC3339Nano)}
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return
	}

	return decoder.Decode(input)
}
