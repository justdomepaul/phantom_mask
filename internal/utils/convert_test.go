package utils

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type ParseSuite struct {
	suite.Suite
}

func (suite *ParseSuite) TestParseTimeFormat() {
	type want struct {
		Result []DaySchema
	}
	testCases := []struct {
		Label string
		Input string
		Want  want
	}{
		{
			Label: "Parse 1 format",
			Input: "Mon, Wed, Fri 08:00 - 12:00 / Tue, Thur 14:00 - 18:00",
			Want: want{
				Result: []DaySchema{
					{
						Day:       1,
						OpenHour:  8,
						CloseHour: 12,
					},
					{
						Day:       3,
						OpenHour:  8,
						CloseHour: 12,
					},
					{
						Day:       5,
						OpenHour:  8,
						CloseHour: 12,
					},
					{
						Day:       2,
						OpenHour:  14,
						CloseHour: 18,
					},
					{
						Day:       4,
						OpenHour:  14,
						CloseHour: 18,
					},
				},
			},
		},
		{
			Label: "Parse 2 format",
			Input: "Mon - Fri 08:00 - 17:00",
			Want: want{
				Result: []DaySchema{
					{
						Day:       1,
						OpenHour:  8,
						CloseHour: 17,
					},
					{
						Day:       2,
						OpenHour:  8,
						CloseHour: 17,
					},
					{
						Day:       3,
						OpenHour:  8,
						CloseHour: 17,
					},
					{
						Day:       4,
						OpenHour:  8,
						CloseHour: 17,
					},
					{
						Day:       5,
						OpenHour:  8,
						CloseHour: 17,
					},
				},
			},
		},
		{
			Label: "Parse 3 format",
			Input: "Mon - Fri 08:00 - 17:00 / Sat, Sun 08:00 - 12:00",
			Want: want{
				Result: []DaySchema{
					{
						Day:       1,
						OpenHour:  8,
						CloseHour: 17,
					},
					{
						Day:       2,
						OpenHour:  8,
						CloseHour: 17,
					},
					{
						Day:       3,
						OpenHour:  8,
						CloseHour: 17,
					},
					{
						Day:       4,
						OpenHour:  8,
						CloseHour: 17,
					},
					{
						Day:       5,
						OpenHour:  8,
						CloseHour: 17,
					},
					{
						Day:       6,
						OpenHour:  8,
						CloseHour: 12,
					},
					{
						Day:       0,
						OpenHour:  8,
						CloseHour: 12,
					},
				},
			},
		},
		{
			Label: "Parse 4 format",
			Input: "Mon - Wed 08:00 - 17:00 / Thur, Sat 20:00 - 02:00",
			Want: want{
				Result: []DaySchema{
					{
						Day:       1,
						OpenHour:  8,
						CloseHour: 17,
					},
					{
						Day:       2,
						OpenHour:  8,
						CloseHour: 17,
					},
					{
						Day:       3,
						OpenHour:  8,
						CloseHour: 17,
					},
					{
						Day:       4,
						OpenHour:  20,
						CloseHour: 26,
					},
					{
						Day:       6,
						OpenHour:  20,
						CloseHour: 26,
					},
				},
			},
		},
		{
			Label: "Parse 5 format",
			Input: "Mon, Wed, Fri 20:00 - 02:00",
			Want: want{
				Result: []DaySchema{
					{
						Day:       1,
						OpenHour:  20,
						CloseHour: 26,
					},
					{
						Day:       3,
						OpenHour:  20,
						CloseHour: 26,
					},
					{
						Day:       5,
						OpenHour:  20,
						CloseHour: 26,
					},
				},
			},
		},
	}
	for _, tc := range testCases {
		result, err := ParseTimeFormat(tc.Input)
		suite.NoError(err)
		suite.Equal(tc.Want.Result, result)
	}
}

func TestParseSuite(t *testing.T) {
	suite.Run(t, new(ParseSuite))
}

func BenchmarkParseTimeFormat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ParseTimeFormat("Mon - Wed 08:00 - 17:00 / Thur, Sat 20:00 - 02:00")
	}
}
