package stool

import (
	"time"
)

const OneDaysSeconds int64 = 24 * 60 * 60

func TimeNowSecond1970() int64 {
	return time.Now().Unix()
}

func TimeNowMilliSecond1970() int64 {
	return time.Now().UnixNano() / 1e6
}

func GetNowISO8601() string {
	return time.Now().UTC().Format("2006-01-02T15:04:05Z")
}

func TimeNowYMDHMS(timeZoneOffset int) string {
	cstZone := time.FixedZone("CST", timeZoneOffset*3600)
	return time.Now().In(cstZone).Format("2006-01-02 15:04:05")
}

func UnixTimeFormatToYMDHMS(second int64, timeZoneOffset int) string {
	cstZone := time.FixedZone("CST", timeZoneOffset*3600)
	return time.Unix(second, 0).In(cstZone).Format("2006-01-02 15:04:05")
}

func UnixTimeFormat(second int64, format string, timeZoneOffset int) string {
	cstZone := time.FixedZone("CST", timeZoneOffset*3600)
	return time.Unix(second, 0).In(cstZone).Format(format)
}

func UnixTimeFormatToYMD(second int64, timeZoneOffset int) string {
	cstZone := time.FixedZone("CST", timeZoneOffset*3600)
	return time.Unix(second, 0).In(cstZone).Format("2006-01-02")
}

func YMDHMSToUnixTime(date string, timeZoneOffset int) int64 {
	cstZone := time.FixedZone("CST", timeZoneOffset*3600)
	tm, _ := time.ParseInLocation("2006-01-02 15:04:05", date, cstZone)
	return tm.Unix()
}

// GetDayZeroTime 获取以当天为基数的某天零点时间戳
func GetDayZeroTime(delta int, timeZoneOffset int) int64 {
	cstZone := time.FixedZone("CST", timeZoneOffset*3600)
	timeStr := time.Now().In(cstZone).Format("2006-01-02")
	t2, _ := time.ParseInLocation("2006-01-02", timeStr, cstZone)
	return t2.AddDate(0, 0, delta).Unix()
}

// GetDayTime 获取以当天为基数的某天某时某分某秒时间戳
func GetDayTime(hour, min, sec int, timeZoneOffset int) int64 {
	cstZone := time.FixedZone("CST", timeZoneOffset*3600)
	timeStr := time.Now().In(cstZone).Format("2006-01-02")
	t2, _ := time.ParseInLocation("2006-01-02", timeStr, cstZone)
	duration := time.Duration(hour*60*60+min*60+sec) * time.Second
	return t2.Add(duration).Unix()
}
