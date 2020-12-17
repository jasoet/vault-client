package util

import (
	"fmt"
	"strings"
)

//ToArrStr Convert slice of interface{} to array of string
func ToArrStr(input []interface{}) (result []string) {
	result = []string{}
	for _, v := range input {
		result = append(result, fmt.Sprint(v))
	}
	return
}

//ToArrStrPrefix Convert slice of interface{} to array of string with prefix added
func ToArrStrPrefix(input []interface{}, prefix string) (result []string) {
	result = []string{}
	for _, v := range input {
		result = append(result, fmt.Sprintf("%v%v", prefix, v))
	}
	return
}

//ToArrStrPrefixPath Convert slice of interface{} to array of string with prefix added
//Prefix will modify to include `/` if not available
func ToArrStrPrefixPath(input []interface{}, prefix string) (result []string) {
	if !strings.HasSuffix(prefix, "/") {
		prefix = fmt.Sprintf("%v/", prefix)
	}

	return ToArrStrPrefix(input, prefix)
}
