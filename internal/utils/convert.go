package utils

import (
	"fmt"
	"math"
	"time"
)

type DaySchema struct {
	Day       int64
	OpenHour  float64
	CloseHour float64
}

var weekday = map[string]int{
	"Sun":  int(time.Sunday),
	"Mon":  int(time.Monday),
	"Tue":  int(time.Tuesday),
	"Wed":  int(time.Wednesday),
	"Thur": int(time.Thursday),
	"Fri":  int(time.Friday),
	"Sat":  int(time.Saturday),
}

var weekdayStr = map[int]string{
	int(time.Sunday):    "Sun",
	int(time.Monday):    "Mon",
	int(time.Tuesday):   "Tue",
	int(time.Wednesday): "Wed",
	int(time.Thursday):  "Thur",
	int(time.Friday):    "Fri",
	int(time.Saturday):  "Sat",
}

func ParseTimeFormat(input string) ([]DaySchema, error) {
	runes := []rune(input)
	runes = append(runes, ' ')
	var (
		result  []DaySchema
		chars   []rune
		days    []int
		times   [][]rune
		isTime  bool
		isRange bool
	)
	for _, char := range runes {
		switch char {
		case ' ':
			if isTime {
				if len(chars) > 0 {
					times = append(times, chars)
				}
				if len(times) == 2 && isRange {
					openHour, err := time.Parse("15:04", string(times[0]))
					if err != nil {
						return nil, err
					}
					openH, err := time.ParseDuration(fmt.Sprintf("%dh%dm", openHour.Hour(), openHour.Minute()))
					if err != nil {
						return nil, err
					}
					closeHour, err := time.Parse("15:04", string(times[1]))
					if err != nil {
						return nil, err
					}
					closeH, err := time.ParseDuration(fmt.Sprintf("%dh%dm", closeHour.Hour(), closeHour.Minute()))
					if err != nil {
						return nil, err
					}
					closeHDetermine := closeH.Hours()
					if closeH.Hours() < openH.Hours() {
						closeHDetermine = closeH.Hours() + 24
					}
					for _, d := range days {
						result = append(result, DaySchema{
							Day:       int64(d),
							OpenHour:  math.Round(openH.Hours()*100) / 100,
							CloseHour: math.Round(closeHDetermine*100) / 100,
						})
					}
					isTime = false
					isRange = false
					days = []int{}
					times = [][]rune{}
				}
				chars = []rune{}
				continue
			} else {
				if len(chars) > 0 {
					days = append(days, weekday[string(chars)])
				}
				if len(days) == 2 && isRange {
					start := days[0]
					end := days[1]
					days = []int{}
					for i := start; i <= end; i++ {
						days = append(days, i)
					}
				}
				chars = []rune{}
				continue
			}
		case ',':
		case '/':
			isTime = false
			isRange = false
		case '-':
			isRange = true
		default:
			if isNumberChar(char) {
				isTime = true
			}
			chars = append(chars, char)
		}
	}
	return result, nil
}

func isNumberChar(s rune) bool {
	return '0' <= s && s <= '9'
}
