package stool

import (
	"fmt"
	"testing"
)

func TestGetGUID(t *testing.T) {
	i := 1
	m := make(map[string]struct{}, 0)
	for i < 5000 {
		c := Code8()
		if _, ok := m[c]; ok {
			fmt.Println(i)
			panic("xxx")
			return
		}

		m[c] = struct{}{}
		fmt.Println(c)
		i = i + 1
	}
}

func TestCode8Many(t *testing.T) {
	fmt.Println(Code8Many(4))
}

func TestGetGUID2(t *testing.T) {
	fmt.Println(GetGUID())
}

func TestToJsonStringClear(t *testing.T) {
	a := struct {
		A string `json:"a"`
	}{}
	fmt.Println(ToJsonStringClear(a))
}
