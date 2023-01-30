package goccss_test

import (
	"testing"

	goccss "github.com/pandatix/go-ccss"
)

var Gccss *goccss.CCSS
var Gerr error

func BenchmarkParseVector_Base(b *testing.B) {
	benchmarkParseVector("AV:L/AC:L/Au:N/C:P/I:P/A:P/PL:U/EM:A", b)
}

func BenchmarkParseVector_WithTempAndEnv(b *testing.B) {
	benchmarkParseVector("AV:L/AC:L/Au:N/C:P/I:P/A:P/PL:U/EM:A/GEL:L/GRL:M/LVP:L/PTV:L/LRL:L/EC:C/EI:C/EA:C/CDP:L/CR:M/IR:M/AR:L", b)
}

func benchmarkParseVector(vector string, b *testing.B) {
	var ccss *goccss.CCSS
	var err error
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ccss, err = goccss.ParseVector(vector)
	}
	Gccss = ccss
	Gerr = err
}

var Gstr string

func BenchmarkCCSSVector(b *testing.B) {
	ccss, _ := goccss.ParseVector("AV:L/AC:L/Au:N/C:P/I:P/A:P/PL:U/EM:A/GEL:L/GRL:M/LVP:L/PTV:L/LRL:L/EC:C/EI:C/EA:C/CDP:L/CR:M/IR:M/AR:L")
	var str string
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		str = ccss.Vector()
	}
	Gstr = str
}

var Gget string

func BenchmarkCCSSGet(b *testing.B) {
	const abv = "Au"
	ccss, _ := goccss.ParseVector("AV:L/AC:L/Au:N/C:P/I:P/A:P/PL:U/EM:A/GEL:L/GRL:M/LVP:L/PTV:L/LRL:L/EC:C/EI:C/EA:C/CDP:L/CR:M/IR:M/AR:L")
	var get string
	var err error
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		get, err = ccss.Get(abv)
	}
	Gget = get
	Gerr = err
}

func BenchmarkCCSSSet(b *testing.B) {
	const abv = "Au"
	const value = "S"
	ccss, _ := goccss.ParseVector("AV:L/AC:L/Au:N/C:P/I:P/A:P/PL:U/EM:A/GEL:L/GRL:M/LVP:L/PTV:L/LRL:L/EC:C/EI:C/EA:C/CDP:L/CR:M/IR:M/AR:L")
	var err error
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err = ccss.Set(abv, value)
	}
	Gerr = err
}

var Gscore float64

func BenchmarkCCSSBaseScore(b *testing.B) {
	var score float64
	ccss, _ := goccss.ParseVector("AV:L/AC:L/Au:N/C:P/I:P/A:P/PL:U/EM:A/GEL:L/GRL:M/LVP:L/PTV:L/LRL:L/EC:C/EI:C/EA:C/CDP:L/CR:M/IR:M/AR:L")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		score = ccss.BaseScore()
	}
	Gscore = score
}

func BenchmarkCCSSTemporalScore(b *testing.B) {
	var score float64
	ccss, _ := goccss.ParseVector("AV:L/AC:L/Au:N/C:P/I:P/A:P/PL:U/EM:A/GEL:L/GRL:M/LVP:L/PTV:L/LRL:L/EC:C/EI:C/EA:C/CDP:L/CR:M/IR:M/AR:L")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		score = ccss.TemporalScore()
	}
	Gscore = score
}

func BenchmarkCCSSEnvironmentalScore(b *testing.B) {
	var score float64
	ccss, _ := goccss.ParseVector("AV:L/AC:L/Au:N/C:P/I:P/A:P/PL:U/EM:A/GEL:L/GRL:M/LVP:L/PTV:L/LRL:L/EC:C/EI:C/EA:C/CDP:L/CR:M/IR:M/AR:L")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		score = ccss.EnvironmentalScore()
	}
	Gscore = score
}
