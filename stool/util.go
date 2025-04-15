package stool

import (
	"encoding/json"
)

func ToJsonString(object interface{}) string {
	js, _ := json.Marshal(object)
	return string(js)
}

func ToJsonStringClear(object interface{}) string {
	js, _ := json.MarshalIndent(object, "", "\t")
	return string(js)
}

func NotInList(i string, list []string) bool {
	if list == nil || len(list) == 0 {
		return true
	}

	for _, v := range list {
		if v == i {
			return false
		}
	}

	return true
}
