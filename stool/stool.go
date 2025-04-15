package stool

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"
)

func Round(f float64, p int) float64 {
	shift := math.Pow(10, float64(p))
	return math.Round(f*shift) / shift
}

func UUID() string {
	return GetGUID()
}

func GetString(s *string) string {
	if s == nil {
		return ""
	}

	return *s
}

func GetStringPointer(s string) *string {
	if s == "" {
		return nil
	}

	return &s
}

func GetNearNum(n float64, y int64) int64 {
	temp := math.Ceil(n)
	temp2 := int64(temp)
	if y == 0 {
		return temp2
	}

	if temp2%y > 0 {
		return ((temp2 / y) + 1) * y
	}

	return temp2
}

func VerifyEmailFormat(email string) bool {
	pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`

	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

func VerifyMobileFormat(mobileNum string) bool {
	if strings.HasPrefix(mobileNum, "86") {
		mobileNum = strings.TrimPrefix(mobileNum, "86")
	} else {
		mobileNum = strings.TrimPrefix(mobileNum, "93")
	}

	regular := "^(13[0-9]|14[01456879]|15[0-35-9]|16[2567]|17[0-8]|18[0-9]|19[0-35-9])\\d{8}$"

	reg := regexp.MustCompile(regular)
	return reg.MatchString(mobileNum)
}

func MapSubMap(a map[string]int64, b map[string]int64) map[string]int64 {
	temp := map[string]int64{}
	for k, v := range a {
		existInB := false
		for kk, _ := range b {
			if k == kk {
				existInB = true
				break
			}
		}

		if !existInB {
			temp[k] = v
		}
	}

	return temp
}

func Percent(num int) bool {
	if num <= 0 {
		return false
	}

	if num >= 100 {
		return true
	}

	prizeList := make([]string, 0)
	magicFactor := 100

	for i := 0; i < magicFactor*num; i++ {
		prizeList = append(prizeList, "1")
	}

	for i := 0; i < magicFactor*(100-num); i++ {
		prizeList = append(prizeList, "0")
	}

	Shuffle(prizeList)
	Shuffle(prizeList)
	Shuffle(prizeList)

	if prizeList[0] == "1" {
		return true
	} else {
		return false
	}
}

func Shuffle(nums []string) []string {
	if len(nums) <= 1 {
		return nums
	}

	rand2 := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	rand2.Seed(time.Now().UTC().UnixNano())
	for i := len(nums); i > 0; i-- {
		last := i - 1
		idx := rand2.Intn(i)
		nums[last], nums[idx] = nums[idx], nums[last]
	}
	return nums
}

var charsRandomNumbers = []string{
	"0", "1", "2", "3", "4", "5",
	"6", "7", "8", "9", "A", "B", "C", "D", "E", "F", "G", "H", "I",
	"J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V",
	"W", "X", "Y", "Z",
}

func GetRandomNumbers() string {
	s := ""
	rand2 := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	for i := 0; i < 6; i++ {
		rand2.Seed(time.Now().UTC().UnixNano())
		s = fmt.Sprintf("%s%s", s, charsRandomNumbers[rand2.Intn(len(charsRandomNumbers))])
	}

	return s
}

func ReadFromFile(filepath string) ([]byte, error) {
	return os.ReadFile(filepath)
}
