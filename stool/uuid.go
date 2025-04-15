package stool

import (
	"fmt"
	"github.com/gofrs/uuid"
	"math/rand"
	"strings"
	"time"
)

var chars = []string{
	"a", "b", "c", "d", "e", "f",
	"g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s",
	"t", "u", "v", "w", "x", "y", "z", "0", "1", "2", "3", "4", "5",
	"6", "7", "8", "9", "A", "B", "C", "D", "E", "F", "G", "H", "I",
	"J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V",
	"W", "X", "Y", "Z",
}

func GetGUID() (valueGUID string) {
	objID, _ := uuid.NewV4()
	objIdStr := objID.String()
	objIdStr = strings.Replace(objIdStr, "-", "", -1)
	valueGUID = objIdStr
	return valueGUID
}

func Code8() string {
	objID, _ := uuid.NewV4()
	b := objID.Bytes()
	if len(b) < 16 {
		return ""
	}

	cool := make([]string, 8)
	cool[0] = chars[(b[0]+b[1])%62]
	cool[1] = chars[(b[2]+b[3])%62]
	cool[2] = chars[(b[4]+b[5])%62]
	cool[3] = chars[(b[6]+b[7])%62]
	cool[4] = chars[(b[8]+b[9])%62]
	cool[5] = chars[(b[10]+b[11])%62]
	cool[6] = chars[(b[12]+b[13])%62]
	cool[7] = chars[(b[14]+b[15])%62]

	code := strings.Join(cool, "")
	return code
}

func Code8Many(num int) (list []string) {
	i := 0

	for i < num {
		c := Code8()
		if c == "" {
			continue
		}
		list = append(list, c)
		i = i + 1
	}

	return list
}

func GetRandomNumber(n int) string {
	if n <= 0 {
		return GetGUID()
	}

	s := ""
	rand.Seed(time.Now().UnixMicro())
	for i := 0; i < n; i++ {
		s = fmt.Sprintf("%s%d", s, rand.Intn(10))
	}

	return s
}
