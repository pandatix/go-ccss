package goccss

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testsParseVector = map[string]struct {
	Vector       string
	ExpectedCCSS *CCSS
	ExpectedErr  error
}{
	"CCE-4675-5": {
		Vector: "AV:N/AC:L/Au:N/C:N/I:P/A:N/PL:ND/EM:P",
		ExpectedCCSS: &CCSS{
			u0: 0b10101000,
			u1: 0b01000010,
			u2: 0b00000000,
			u3: 0b00000000,
			u4: 0b00000000,
			u5: 0b00000000,
		},
		ExpectedErr: nil,
	},
	"CCE-4693-8 1": {
		Vector: "AV:L/AC:L/Au:N/C:N/I:N/A:P/PL:ND/EM:P",
		ExpectedCCSS: &CCSS{
			u0: 0b00101000,
			u1: 0b00010010,
			u2: 0b00000000,
			u3: 0b00000000,
			u4: 0b00000000,
			u5: 0b00000000,
		},
		ExpectedErr: nil,
	},
	"CCE-4693-8 2": {
		Vector: "AV:L/AC:L/Au:N/C:P/I:P/A:P/PL:ND/EM:A",
		ExpectedCCSS: &CCSS{
			u0: 0b00101001,
			u1: 0b01010000,
			u2: 0b00000000,
			u3: 0b00000000,
			u4: 0b00000000,
			u5: 0b00000000,
		},
		ExpectedErr: nil,
	},
	"CCE-2786-2": {
		Vector: "AV:L/AC:L/Au:N/C:N/I:N/A:P/PL:ND/EM:A",
		ExpectedCCSS: &CCSS{
			u0: 0b00101000,
			u1: 0b00010000,
			u2: 0b00000000,
			u3: 0b00000000,
			u4: 0b00000000,
			u5: 0b00000000,
		},
		ExpectedErr: nil,
	},
	"CCE-2363-0": {
		Vector: "AV:L/AC:H/Au:N/C:P/I:N/A:N/PL:ND/EM:A",
		ExpectedCCSS: &CCSS{
			u0: 0b00001001,
			u1: 0b00000000,
			u2: 0b00000000,
			u3: 0b00000000,
			u4: 0b00000000,
			u5: 0b00000000,
		},
		ExpectedErr: nil,
	},
	"CCE-2366-3": {
		Vector: "AV:L/AC:L/Au:N/C:N/I:N/A:C/PL:ND/EM:A",
		ExpectedCCSS: &CCSS{
			u0: 0b00101000,
			u1: 0b00100000,
			u2: 0b00000000,
			u3: 0b00000000,
			u4: 0b00000000,
			u5: 0b00000000,
		},
		ExpectedErr: nil,
	},
	"CCE-4208-5": {
		Vector: "AV:N/AC:H/Au:N/C:P/I:N/A:N/PL:ND/EM:A",
		ExpectedCCSS: &CCSS{
			u0: 0b10001001,
			u1: 0b00000000,
			u2: 0b00000000,
			u3: 0b00000000,
			u4: 0b00000000,
			u5: 0b00000000,
		},
		ExpectedErr: nil,
	},
	"CCE-2519-7 1": {
		Vector: "AV:N/AC:L/Au:N/C:N/I:N/A:P/PL:ND/EM:P",
		ExpectedCCSS: &CCSS{
			u0: 0b10101000,
			u1: 0b00010010,
			u2: 0b00000000,
			u3: 0b00000000,
			u4: 0b00000000,
			u5: 0b00000000,
		},
		ExpectedErr: nil,
	},
	"CCE-2519-7 2": {
		Vector: "AV:N/AC:H/Au:N/C:P/I:P/A:P/PL:U/EM:A",
		ExpectedCCSS: &CCSS{
			u0: 0b10001001,
			u1: 0b01011000,
			u2: 0b00000000,
			u3: 0b00000000,
			u4: 0b00000000,
			u5: 0b00000000,
		},
		ExpectedErr: nil,
	},
	"CCE-3171-6": {
		Vector: "AV:N/AC:L/Au:N/C:P/I:N/A:P/PL:ND/EM:A",
		ExpectedCCSS: &CCSS{
			u0: 0b10101001,
			u1: 0b00010000,
			u2: 0b00000000,
			u3: 0b00000000,
			u4: 0b00000000,
			u5: 0b00000000,
		},
		ExpectedErr: nil,
	},
	"CCE-3047-8 1": {
		Vector: "AV:L/AC:L/Au:N/C:N/I:N/A:P/PL:ND/EM:A",
		ExpectedCCSS: &CCSS{
			u0: 0b00101000,
			u1: 0b00010000,
			u2: 0b00000000,
			u3: 0b00000000,
			u4: 0b00000000,
			u5: 0b00000000,
		},
		ExpectedErr: nil,
	},
	"CCE-3047-8 2": {
		Vector: "AV:L/AC:L/Au:N/C:N/I:P/A:N/PL:ND/EM:P",
		ExpectedCCSS: &CCSS{
			u0: 0b00101000,
			u1: 0b01000010,
			u2: 0b00000000,
			u3: 0b00000000,
			u4: 0b00000000,
			u5: 0b00000000,
		},
		ExpectedErr: nil,
	},
	"CCE-4191-3 1": {
		Vector: "AV:L/AC:L/Au:N/C:N/I:N/A:P/PL:ND/EM:P",
		ExpectedCCSS: &CCSS{
			u0: 0b00101000,
			u1: 0b00010010,
			u2: 0b00000000,
			u3: 0b00000000,
			u4: 0b00000000,
			u5: 0b00000000,
		},
		ExpectedErr: nil,
	},
	"CCE-4191-3 2": {
		Vector: "AV:L/AC:L/Au:N/C:P/I:N/A:N/PL:ND/EM:A",
		ExpectedCCSS: &CCSS{
			u0: 0b00101001,
			u1: 0b00000000,
			u2: 0b00000000,
			u3: 0b00000000,
			u4: 0b00000000,
			u5: 0b00000000,
		},
		ExpectedErr: nil,
	},
	"CCE-3245-8": {
		Vector: "AV:N/AC:L/Au:N/C:P/I:P/A:P/PL:ND/EM:P",
		ExpectedCCSS: &CCSS{
			u0: 0b10101001,
			u1: 0b01010010,
			u2: 0b00000000,
			u3: 0b00000000,
			u4: 0b00000000,
			u5: 0b00000000,
		},
		ExpectedErr: nil,
	},
	"CCE-2776-3 1": {
		Vector: "AV:L/AC:L/Au:N/C:P/I:P/A:P/PL:U/EM:A",
		ExpectedCCSS: &CCSS{
			u0: 0b00101001,
			u1: 0b01011000,
			u2: 0b00000000,
			u3: 0b00000000,
			u4: 0b00000000,
			u5: 0b00000000,
		},
		ExpectedErr: nil,
	},
	"CCE-2776-3 2": {
		Vector: "AV:L/AC:L/Au:N/C:P/I:P/A:P/PL:U/EM:A/GEL:L/GRL:M",
		ExpectedCCSS: &CCSS{
			u0: 0b00101001,
			u1: 0b01011000,
			u2: 0b10010000,
			u3: 0b00000000,
			u4: 0b00000000,
			u5: 0b00000000,
		},
		ExpectedErr: nil,
	},
	"CCE-2776-3 3": {
		Vector: "AV:L/AC:L/Au:N/C:P/I:P/A:P/PL:U/EM:A/LVP:L/PTV:L/LRL:L/EC:C/EI:C/EA:C/CDP:L/CR:M/IR:M/AR:L",
		ExpectedCCSS: &CCSS{
			u0: 0b00101001,
			u1: 0b01011000,
			u2: 0b00000010,
			u3: 0b01010111,
			u4: 0b11101010,
			u5: 0b10010000,
		},
	},
	"CCE-2776-3 2&3": {
		Vector: "AV:L/AC:L/Au:N/C:P/I:P/A:P/PL:U/EM:A/GEL:L/GRL:M/LVP:L/PTV:L/LRL:L/EC:C/EI:C/EA:C/CDP:L/CR:M/IR:M/AR:L",
		ExpectedCCSS: &CCSS{
			u0: 0b00101001,
			u1: 0b01011000,
			u2: 0b10010010,
			u3: 0b01010111,
			u4: 0b11101010,
			u5: 0b10010000,
		},
	},
}

func TestParseVector(t *testing.T) {
	t.Parallel()

	for testname, tt := range testsParseVector {
		t.Run(testname, func(t *testing.T) {
			assert := assert.New(t)

			ccss, err := ParseVector(tt.Vector)

			assert.Equal(tt.ExpectedCCSS, ccss)
			assert.Equal(tt.ExpectedErr, err)

			if ccss != nil {
				newVec := ccss.Vector()
				assert.Equal(tt.Vector, newVec)
			}
		})
	}
}

func TestScores(t *testing.T) {
	t.Parallel()

	var tests = map[string]struct {
		CCSS                       *CCSS
		ExpectedBaseScore          float64
		ExpectedTemporalScore      float64
		ExpectedEnvironmentalScore float64
	}{
		"CCE-4675-5": {
			CCSS:                       must(ParseVector("AV:N/AC:L/Au:N/C:N/I:P/A:N/PL:ND/EM:P")),
			ExpectedBaseScore:          5.0,
			ExpectedTemporalScore:      5.0,
			ExpectedEnvironmentalScore: 5.0,
		},
		"CCE-4693-8 1": {
			CCSS:                       must(ParseVector("AV:L/AC:L/Au:N/C:N/I:N/A:P/PL:ND/EM:P")),
			ExpectedBaseScore:          2.1,
			ExpectedTemporalScore:      2.1,
			ExpectedEnvironmentalScore: 2.1,
		},
		"CCE-4693-8 2": {
			CCSS:                       must(ParseVector("AV:L/AC:L/Au:N/C:P/I:P/A:P/PL:ND/EM:A")),
			ExpectedBaseScore:          4.6,
			ExpectedTemporalScore:      4.6,
			ExpectedEnvironmentalScore: 4.6,
		},
		"CCE-2786-2": {
			CCSS:                       must(ParseVector("AV:L/AC:L/Au:N/C:N/I:N/A:P/PL:ND/EM:A")),
			ExpectedBaseScore:          2.1,
			ExpectedTemporalScore:      2.1,
			ExpectedEnvironmentalScore: 2.1,
		},
		"CCE-2363-0 1": {
			CCSS:                       must(ParseVector("AV:L/AC:H/Au:N/C:P/I:N/A:N/PL:ND/EM:A")),
			ExpectedBaseScore:          1.2,
			ExpectedTemporalScore:      1.2,
			ExpectedEnvironmentalScore: 1.2,
		},
		"CCE-2363-0 2": {
			CCSS:                       must(ParseVector("AV:L/AC:L/Au:N/C:N/I:N/A:P/PL:ND/EM:P")),
			ExpectedBaseScore:          2.1,
			ExpectedTemporalScore:      2.1,
			ExpectedEnvironmentalScore: 2.1,
		},
		"CCE-2366-3 1": {
			CCSS:                       must(ParseVector("AV:L/AC:L/Au:N/C:N/I:N/A:C/PL:ND/EM:A")),
			ExpectedBaseScore:          4.9,
			ExpectedTemporalScore:      4.9,
			ExpectedEnvironmentalScore: 4.9,
		},
		"CCE-2366-3 2": {
			CCSS:                       must(ParseVector("AV:L/AC:L/Au:N/C:N/I:N/A:P/PL:ND/EM:A")),
			ExpectedBaseScore:          2.1,
			ExpectedTemporalScore:      2.1,
			ExpectedEnvironmentalScore: 2.1,
		},
		"CCE-4208-5": {
			CCSS:                       must(ParseVector("AV:N/AC:H/Au:N/C:P/I:N/A:N/PL:ND/EM:A")),
			ExpectedBaseScore:          2.6,
			ExpectedTemporalScore:      2.6,
			ExpectedEnvironmentalScore: 2.6,
		},
		"CCE-2519-7 1": {
			CCSS:                       must(ParseVector("AV:N/AC:L/Au:N/C:N/I:N/A:P/PL:ND/EM:P")),
			ExpectedBaseScore:          5.0,
			ExpectedTemporalScore:      5.0,
			ExpectedEnvironmentalScore: 5.0,
		},
		"CCE-2519-7 2": {
			CCSS:                       must(ParseVector("AV:N/AC:H/Au:N/C:P/I:P/A:P/PL:U/EM:A")),
			ExpectedBaseScore:          5.1,
			ExpectedTemporalScore:      5.1,
			ExpectedEnvironmentalScore: 5.1,
		},
		"CCE-3171-6": {
			CCSS:                       must(ParseVector("AV:N/AC:L/Au:N/C:P/I:N/A:P/PL:ND/EM:A")),
			ExpectedBaseScore:          6.4,
			ExpectedTemporalScore:      6.4,
			ExpectedEnvironmentalScore: 6.4,
		},
		"CCE-3047-8 1": {
			CCSS:                       must(ParseVector("AV:L/AC:L/Au:N/C:N/I:N/A:P/PL:ND/EM:A")),
			ExpectedBaseScore:          2.1,
			ExpectedTemporalScore:      2.1,
			ExpectedEnvironmentalScore: 2.1,
		},
		"CCE-3047-8 2": {
			CCSS:                       must(ParseVector("AV:L/AC:L/Au:N/C:N/I:P/A:N/PL:ND/EM:P")),
			ExpectedBaseScore:          2.1,
			ExpectedTemporalScore:      2.1,
			ExpectedEnvironmentalScore: 2.1,
		},
		"CCE-4191-3 1": {
			CCSS:                       must(ParseVector("AV:L/AC:L/Au:N/C:N/I:N/A:P/PL:ND/EM:P")),
			ExpectedBaseScore:          2.1,
			ExpectedTemporalScore:      2.1,
			ExpectedEnvironmentalScore: 2.1,
		},
		"CCE-4191-3 2": {
			CCSS:                       must(ParseVector("AV:L/AC:L/Au:N/C:P/I:N/A:N/PL:ND/EM:A")),
			ExpectedBaseScore:          2.1,
			ExpectedTemporalScore:      2.1,
			ExpectedEnvironmentalScore: 2.1,
		},
		"CCE-3245-8": {
			CCSS:                       must(ParseVector("AV:N/AC:L/Au:N/C:P/I:P/A:P/PL:ND/EM:P")),
			ExpectedBaseScore:          7.5,
			ExpectedTemporalScore:      7.5,
			ExpectedEnvironmentalScore: 7.5,
		},
		"CCE-2776-3 1": {
			CCSS:                       must(ParseVector("AV:L/AC:L/Au:N/C:P/I:P/A:P/PL:U/EM:A")),
			ExpectedBaseScore:          4.6,
			ExpectedTemporalScore:      4.6,
			ExpectedEnvironmentalScore: 4.6,
		},
		"CCE-2776-3 2": {
			CCSS:                  must(ParseVector("AV:L/AC:L/Au:N/C:P/I:P/A:P/PL:U/EM:A/GEL:L/GRL:M")),
			ExpectedBaseScore:     4.6,
			ExpectedTemporalScore: 3.7,
			// XXX results are unsure due to intrication of temporal and environmental metric groups. Set to what was computed.
			ExpectedEnvironmentalScore: 4.3,
		},
		"CCE-2776-3 3": {
			CCSS:                       must(ParseVector("AV:L/AC:L/Au:N/C:P/I:P/A:P/PL:U/EM:A/GEL:L/GRL:M/LVP:L/PTV:L/LRL:L/EC:C/EI:C/EA:C/CDP:L/CR:M/IR:M/AR:L")),
			ExpectedBaseScore:          4.6,
			ExpectedTemporalScore:      3.7,
			ExpectedEnvironmentalScore: 6.1,
		},
	}

	for testname, tt := range tests {
		t.Run(testname, func(t *testing.T) {
			assert := assert.New(t)

			baseScore := tt.CCSS.BaseScore()
			temporalScore := tt.CCSS.TemporalScore()
			environmentalScore := tt.CCSS.EnvironmentalScore()

			assert.Equal(tt.ExpectedBaseScore, baseScore)
			assert.Equal(tt.ExpectedTemporalScore, temporalScore)
			assert.Equal(tt.ExpectedEnvironmentalScore, environmentalScore)
		})
	}
}

func must[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}
