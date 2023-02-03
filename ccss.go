package goccss

import (
	"math"
	"strings"
	"sync"
	"unsafe"
)

var order = [][]string{
	{"AV", "AC", "Au", "C", "I", "A", "PL", "EM"}, // Base metrics
	{"GEL", "GRL"}, // Temporal metrics
	{"LVP", "PTV", "LRL", "EC", "EI", "EA", "CDP", "CR", "IR", "AR"}, // Environmental metrics
}

// ParseVector parses a CCSS vector.
func ParseVector(vector string) (*CCSS, error) {
	// Split parts
	partsPtr := splitPool.Get()
	defer splitPool.Put(partsPtr)
	pts := partsPtr.([]string)
	ei := split(pts, vector)
	pts = pts[:ei+1]

	// Work on each CCSS part
	ccss := &CCSS{
		u0: 0,
		u1: 0,
		u2: 0,
		u3: 0,
		u4: 0,
		u5: 0, // last 4 bits not used
	}

	slci := 0
	i := 0
	for _, pt := range pts {
		abv, v, _ := strings.Cut(pt, ":")
		tgt := ""
		switch slci {
		case 0, 2:
			tgt = order[slci][i]
		case 1:
			tgt = order[1][i]
			if i == 0 && tgt != abv {
				slci++
				tgt = order[2][0]
			}
		default:
			return nil, ErrInvalidMetricValue
		}
		if abv != tgt {
			return nil, ErrInvalidMetricOrder
		}

		if err := ccss.Set(abv, v); err != nil {
			return nil, err
		}

		// Go to next element in slice, or next slice if fully consumed
		i++
		if i == len(order[slci]) {
			slci++
			i = 0
		}
	}
	// Check whole last metric group is specified in vector (=> i == 0)
	if i != 0 {
		return nil, ErrTooShortVector
	}
	// Check whole last metric group is specified in vector (=> i == 0)
	if i != 0 {
		return nil, ErrTooShortVector
	}
	return ccss, nil
}

var splitPool = sync.Pool{
	New: func() any {
		buf := make([]string, 20)
		return buf
	},
}

func split(dst []string, vector string) int {
	start := 0
	curr := 0
	l := len(vector)
	i := 0
	for ; i < l; i++ {
		if vector[i] == '/' {
			dst[curr] = vector[start:i]

			start = i + 1
			curr++

			if curr == 19 {
				break
			}
		}
	}
	dst[curr] = vector[start:]
	return curr
}

func (ccss CCSS) Vector() string {
	l := lenVec(&ccss)
	b := make([]byte, 0, l)

	// Base
	app(&b, "AV:", ccss.get("AV"))
	app(&b, "/AC:", ccss.get("AC"))
	app(&b, "/Au:", ccss.get("Au"))
	app(&b, "/C:", ccss.get("C"))
	app(&b, "/I:", ccss.get("I"))
	app(&b, "/A:", ccss.get("A"))
	app(&b, "/PL:", ccss.get("PL"))
	app(&b, "/EM:", ccss.get("EM"))

	// Temporal
	gel, grl := ccss.get("GEL"), ccss.get("GRL")
	if gel != "ND" || grl != "ND" {
		app(&b, "/GEL:", gel)
		app(&b, "/GRL:", grl)
	}

	// Environmental
	lvp, ptv, lrl, ec, ei, ea, cdp, cr, ir, ar := ccss.get("LVP"), ccss.get("PTV"), ccss.get("LRL"), ccss.get("EC"), ccss.get("EI"), ccss.get("EA"), ccss.get("CDP"), ccss.get("CR"), ccss.get("IR"), ccss.get("AR")
	if lvp != "ND" || ptv != "ND" || lrl != "ND" || ec != "ND" || ei != "ND" || ea != "ND" || cdp != "ND" || cr != "ND" || ir != "ND" || ar != "ND" {
		app(&b, "/LVP:", lvp)
		app(&b, "/PTV:", ptv)
		app(&b, "/LRL:", lrl)
		app(&b, "/EC:", ec)
		app(&b, "/EI:", ei)
		app(&b, "/EA:", ea)
		app(&b, "/CDP:", cdp)
		app(&b, "/CR:", cr)
		app(&b, "/IR:", ir)
		app(&b, "/AR:", ar)
	}

	return unsafe.String(&b[0], l)
}

func lenVec(ccss *CCSS) int {
	// Base:
	// - AV, AC, Au, EM: 4
	// - C, I, A: 3
	// - PL: 3 + len(v)
	// - separators: 7
	// Total: 4*4 + 3*3 + 10 + len(v)
	pl := ccss.get("PL")
	l := 4*4 + 3*3 + 10 + len(pl)

	// Temporal:
	// - GEL, GRL: 4 + len(v)
	// - separators: 2
	// Total: 5*2 + 2 + 2*len(v)
	gel, grl := ccss.get("GEL"), ccss.get("GRL")
	if gel != "ND" || grl != "ND" {
		l += 4*2 + 2 + len(gel) + len(grl)
	}

	// Environmental:
	// - LVP, PTV, LRL, CDP: 4 + len(v)
	// - EC, EI, EA, CR, IR, AR: 3 + len(v)
	// - separators: 10
	// Total: 4*4 + 6*3 + 10 + 10*len(v)
	lvp, ptv, lrl, ec, ei, ea, cdp, cr, ir, ar := ccss.get("LVP"), ccss.get("PTV"), ccss.get("LRL"), ccss.get("EC"), ccss.get("EI"), ccss.get("EA"), ccss.get("CDP"), ccss.get("CR"), ccss.get("IR"), ccss.get("AR")
	if lvp != "ND" || ptv != "ND" || lrl != "ND" || ec != "ND" || ei != "ND" || ea != "ND" || cdp != "ND" || cr != "ND" || ir != "ND" || ar != "ND" {
		l += 4*4 + 6*3 + 10 + len(lvp) + len(ptv) + len(lrl) + len(ec) + len(ei) + len(ea) + len(cdp) + len(cr) + len(ir) + len(ar)
	}

	return l
}

func app(b *[]byte, pre, v string) {
	*b = append(*b, pre...)
	*b = append(*b, v...)
}

type CCSS struct {
	u0, u1, u2, u3, u4, u5 uint8
}

func (ccss CCSS) Get(abv string) (r string, err error) {
	switch abv {
	// Base
	case "AV":
		v := (ccss.u0 & 0b11000000) >> 6
		switch v {
		case av_l:
			r = "L"
		case av_a:
			r = "A"
		case av_n:
			r = "N"
		}
	case "AC":
		v := (ccss.u0 & 0b00110000) >> 4
		switch v {
		case ac_h:
			r = "H"
		case ac_m:
			r = "M"
		case ac_l:
			r = "L"
		}
	case "Au":
		v := (ccss.u0 & 0b00001100) >> 2
		switch v {
		case au_m:
			r = "M"
		case au_s:
			r = "S"
		case au_n:
			r = "N"
		}
	case "C":
		v := ccss.u0 & 0b00000011
		switch v {
		case cia_n:
			r = "N"
		case cia_p:
			r = "P"
		case cia_c:
			r = "C"
		}
	case "I":
		v := (ccss.u1 & 0b11000000) >> 6
		switch v {
		case cia_n:
			r = "N"
		case cia_p:
			r = "P"
		case cia_c:
			r = "C"
		}
	case "A":
		v := (ccss.u1 & 0b00110000) >> 4
		switch v {
		case cia_n:
			r = "N"
		case cia_p:
			r = "P"
		case cia_c:
			r = "C"
		}
	case "PL":
		v := (ccss.u1 & 0b00001100) >> 2
		switch v {
		case pl_nd:
			r = "ND"
		case pl_r:
			r = "R"
		case pl_u:
			r = "U"
		case pl_a:
			r = "A"
		}
	case "EM":
		v := (ccss.u1 & 0b00000010) >> 1
		switch v {
		case em_a:
			r = "A"
		case em_p:
			r = "P"
		}

	// Temporal
	case "GEL":
		v := ((ccss.u1 & 0b00000001) << 2) | ((ccss.u2 & 0b11000000) >> 6)
		switch v {
		case gel_nd:
			r = "ND"
		case gel_n:
			r = "N"
		case gel_l:
			r = "L"
		case gel_m:
			r = "M"
		case gel_h:
			r = "H"
		}
	case "GRL":
		v := (ccss.u2 & 0b00111000) >> 3
		switch v {
		case grl_nd:
			r = "ND"
		case grl_h:
			r = "H"
		case grl_m:
			r = "M"
		case grl_l:
			r = "L"
		case grl_n:
			r = "N"
		}

	// Environmental
	case "LVP":
		v := ccss.u2 & 0b00000111
		switch v {
		case lvp_nd:
			r = "ND"
		case lvp_n:
			r = "N"
		case lvp_l:
			r = "L"
		case lvp_m:
			r = "M"
		case lvp_h:
			r = "H"
		}
	case "PTV":
		v := (ccss.u3 & 0b11000000) >> 6
		switch v {
		case ptv_nd:
			r = "ND"
		case ptv_l:
			r = "L"
		case ptv_m:
			r = "M"
		case ptv_h:
			r = "H"
		}
	case "LRL":
		v := (ccss.u3 & 0b00111000) >> 3
		switch v {
		case lrl_nd:
			r = "ND"
		case lrl_n:
			r = "N"
		case lrl_l:
			r = "L"
		case lrl_m:
			r = "M"
		case lrl_h:
			r = "H"
		}
	case "EC":
		v := (ccss.u3 & 0b00000110) >> 1
		switch v {
		case ecia_nd:
			r = "ND"
		case ecia_n:
			r = "N"
		case ecia_p:
			r = "P"
		case ecia_c:
			r = "C"
		}
	case "EI":
		v := ((ccss.u3 & 0b00000001) << 1) | ((ccss.u4 & 0b10000000) >> 7)
		switch v {
		case ecia_nd:
			r = "ND"
		case ecia_n:
			r = "N"
		case ecia_p:
			r = "P"
		case ecia_c:
			r = "C"
		}
	case "EA":
		v := (ccss.u4 & 0b01100000) >> 5
		switch v {
		case ecia_nd:
			r = "ND"
		case ecia_n:
			r = "N"
		case ecia_p:
			r = "P"
		case ecia_c:
			r = "C"
		}
	case "CDP":
		v := (ccss.u4 & 0b00011100) >> 2
		switch v {
		case cdp_nd:
			r = "ND"
		case cdp_n:
			r = "N"
		case cdp_l:
			r = "L"
		case cdp_lm:
			r = "LM"
		case cdp_mh:
			r = "MH"
		case cdp_h:
			r = "H"
		}
	case "CR":
		v := ccss.u4 & 0b00000011
		switch v {
		case ciar_nd:
			r = "ND"
		case ciar_l:
			r = "L"
		case ciar_m:
			r = "M"
		case ciar_h:
			r = "H"
		}
	case "IR":
		v := (ccss.u5 & 0b11000000) >> 6
		switch v {
		case ciar_nd:
			r = "ND"
		case ciar_l:
			r = "L"
		case ciar_m:
			r = "M"
		case ciar_h:
			r = "H"
		}
	case "AR":
		v := (ccss.u5 & 0b00110000) >> 4
		switch v {
		case ciar_nd:
			r = "ND"
		case ciar_l:
			r = "L"
		case ciar_m:
			r = "M"
		case ciar_h:
			r = "H"
		}

	default:
		return "", &ErrInvalidMetric{Abv: abv}
	}
	return
}

// get is used for internal purposes only.
func (ccss CCSS) get(abv string) string {
	str, err := ccss.Get(abv)
	if err != nil {
		panic(err)
	}
	return str
}

func (ccss *CCSS) Set(abv, value string) error {
	switch abv {
	// Base
	case "AV":
		v, err := validate(value, []string{"L", "A", "N"})
		if err != nil {
			return err
		}
		ccss.u0 = (ccss.u0 & 0b00111111) | (v << 6)
	case "AC":
		v, err := validate(value, []string{"H", "M", "L"})
		if err != nil {
			return err
		}
		ccss.u0 = (ccss.u0 & 0b11001111) | (v << 4)
	case "Au":
		v, err := validate(value, []string{"M", "S", "N"})
		if err != nil {
			return err
		}
		ccss.u0 = (ccss.u0 & 0b11110011) | (v << 2)
	case "C":
		v, err := validate(value, []string{"N", "P", "C"})
		if err != nil {
			return err
		}
		ccss.u0 = (ccss.u0 & 0b11111100) | v
	case "I":
		v, err := validate(value, []string{"N", "P", "C"})
		if err != nil {
			return err
		}
		ccss.u1 = (ccss.u1 & 0b00111111) | (v << 6)
	case "A":
		v, err := validate(value, []string{"N", "P", "C"})
		if err != nil {
			return err
		}
		ccss.u1 = (ccss.u1 & 0b11001111) | (v << 4)
	case "PL":
		v, err := validate(value, []string{"ND", "R", "U", "A"})
		if err != nil {
			return err
		}
		ccss.u1 = (ccss.u1 & 0b11110011) | (v << 2)
	case "EM":
		v, err := validate(value, []string{"A", "P"})
		if err != nil {
			return err
		}
		ccss.u1 = (ccss.u1 & 0b11111101) | (v << 1)

	// Temporal
	case "GEL":
		v, err := validate(value, []string{"ND", "N", "L", "M", "H"})
		if err != nil {
			return err
		}
		ccss.u1 = (ccss.u1 & 0b11111110) | ((v & 0b100) >> 2)
		ccss.u2 = (ccss.u2 & 0b00111111) | ((v & 0b011) << 6)
	case "GRL":
		v, err := validate(value, []string{"ND", "H", "M", "L", "N"})
		if err != nil {
			return err
		}
		ccss.u2 = (ccss.u2 & 0b11000111) | (v << 3)

	// Environmental
	case "LVP":
		v, err := validate(value, []string{"ND", "N", "L", "M", "H"})
		if err != nil {
			return err
		}
		ccss.u2 = (ccss.u2 & 0b11111000) | v
	case "PTV":
		v, err := validate(value, []string{"ND", "L", "M", "H"})
		if err != nil {
			return err
		}
		ccss.u3 = (ccss.u3 & 0b00111111) | (v << 6)
	case "LRL":
		v, err := validate(value, []string{"ND", "N", "L", "M", "H"})
		if err != nil {
			return err
		}
		ccss.u3 = (ccss.u3 & 0b11000111) | (v << 3)
	case "EC":
		v, err := validate(value, []string{"ND", "N", "P", "C"})
		if err != nil {
			return err
		}
		ccss.u3 = (ccss.u3 & 0b11111001) | (v << 1)
	case "EI":
		v, err := validate(value, []string{"ND", "N", "P", "C"})
		if err != nil {
			return err
		}
		ccss.u3 = (ccss.u3 & 0b11111110) | ((v & 0b10) >> 1)
		ccss.u4 = (ccss.u4 & 0b01111111) | ((v & 0b01) << 7)
	case "EA":
		v, err := validate(value, []string{"ND", "N", "P", "C"})
		if err != nil {
			return err
		}
		ccss.u4 = (ccss.u4 & 0b10011111) | (v << 5)
	case "CDP":
		v, err := validate(value, []string{"ND", "N", "L", "LM", "MH", "H"})
		if err != nil {
			return err
		}
		ccss.u4 = (ccss.u4 & 0b11100011) | (v << 2)
	case "CR":
		v, err := validate(value, []string{"ND", "L", "M", "H"})
		if err != nil {
			return err
		}
		ccss.u4 = (ccss.u4 & 0b11111100) | v
	case "IR":
		v, err := validate(value, []string{"ND", "L", "M", "H"})
		if err != nil {
			return err
		}
		ccss.u5 = (ccss.u5 & 0b00110000) | (v << 6)
	case "AR":
		v, err := validate(value, []string{"ND", "L", "M", "H"})
		if err != nil {
			return err
		}
		ccss.u5 = (ccss.u5 & 0b11000000) | (v << 4)

	default:
		return &ErrInvalidMetric{Abv: abv}
	}
	return nil
}

// validate returns the index of value in enabled if matches.
// enabled values have to match the values.go constants order.
func validate(value string, enabled []string) (i uint8, err error) {
	// Check is valid
	for _, enbl := range enabled {
		if value == enbl {
			return i, nil
		}
		i++
	}
	return 0, ErrInvalidMetricValue
}

func (ccss CCSS) BaseScore() float64 {
	impact := ccss.Impact()
	return roundTo1Decimal(((0.6 * impact) + (0.4 * ccss.Exploitability()) - 1.5) * fImpact(impact))
}

func (ccss CCSS) Impact() float64 {
	c := ccss.u0 & 0b00000011
	i := (ccss.u1 & 0b11000000) >> 6
	a := (ccss.u1 & 0b00110000) >> 4
	return 10.41 * (1 - (1-cia(c))*(1-cia(i))*(1-cia(a)))
}

func (ccss CCSS) Exploitability() float64 {
	av := (ccss.u0 & 0b11000000) >> 6
	ac := (ccss.u0 & 0b00110000) >> 4
	au := (ccss.u0 & 0b00001100) >> 2
	return 20 * accessVector(av) * authentication(au) * accessComplexity(ac)
}

func (ccss CCSS) TemporalScore() float64 {
	gel := ((ccss.u1 & 0b00000001) << 2) | ((ccss.u2 & 0b11000000) >> 6)
	grl := (ccss.u2 & 0b00111000) >> 3
	tmpExplt := math.Min(10, ccss.Exploitability()*generalExploitLevel(gel)*generalRemediationLevel(grl))
	return roundTo1Decimal(((0.6 * ccss.Impact()) + (0.4 * tmpExplt) - 1.5) * fImpact(ccss.Impact()))
}

func (ccss CCSS) EnvironmentalScore() float64 {
	c := mod(ccss.u0&0b00000011, (ccss.u3&0b00000110)>>1)
	i := mod((ccss.u1&0b11000000)>>6, ((ccss.u3&0b00000001)<<1)|((ccss.u4&0b10000000)>>7))
	a := mod((ccss.u1&0b00110000)>>4, (ccss.u4&0b01100000)>>5)
	lvp := ccss.u2 & 0b00000111
	ptv := (ccss.u3 & 0b11000000) >> 6
	cr := ccss.u4 & 0b00000011
	ir := (ccss.u5 & 0b11000000) >> 6
	ar := (ccss.u5 & 0b00110000) >> 4
	cdp := (ccss.u4 & 0b00011100) >> 2
	gel := ((ccss.u1 & 0b00000001) << 2) | ((ccss.u2 & 0b11000000) >> 6)
	lrl := (ccss.u3 & 0b00111000) >> 3

	lclExpltLvl := localVulnerabilityPrevalence(lvp) * perceivedTargetValue(ptv)
	envImpact := math.Min(10, 10.41*(1-(1-cia(c)*ciar(cr))*(1-cia(i)*ciar(ir))*(1-cia(a)*ciar(ar)))*collateralDamagePotential(cdp))
	envExplt := math.Min(10, ccss.Exploitability()*generalExploitLevel(gel)*lclExpltLvl*localRemediationLevel(lrl))
	return roundTo1Decimal(((0.6 * envImpact) + (0.4 * envExplt) - 1.5) * fImpact(ccss.Impact()))
}

// Helpers to compute CCSS scores.

func mod(base, modified uint8) uint8 {
	// If "modified" is different of 0, it is different of "X"
	// => shift to one before (skip X index)
	if modified != 0 {
		return modified - 1
	}
	return base
}

func fImpact(v float64) float64 {
	if v == 0 {
		return 0
	}
	return 1.176
}

func accessVector(v uint8) float64 {
	switch v {
	case av_l:
		return 0.395
	case av_a:
		return 0.646
	case av_n:
		return 1.0
	default:
		panic(ErrInvalidMetricValue)
	}
}

func authentication(v uint8) float64 {
	switch v {
	case au_m:
		return 0.45
	case au_s:
		return 0.56
	case au_n:
		return 0.704
	default:
		panic(ErrInvalidMetricValue)
	}
}

func accessComplexity(v uint8) float64 {
	switch v {
	case ac_h:
		return 0.35
	case ac_m:
		return 0.61
	case ac_l:
		return 0.71
	default:
		panic(ErrInvalidMetricValue)
	}
}

func cia(v uint8) float64 {
	switch v {
	case cia_n:
		return 0.0
	case cia_p:
		return 0.275
	case cia_c:
		return 0.660
	default:
		panic(ErrInvalidMetricValue)
	}
}

func generalExploitLevel(v uint8) float64 {
	switch v {
	case gel_n:
		return 0.6
	case gel_l:
		return 0.8
	case gel_m, gel_nd:
		return 1.0
	case gel_h:
		return 1.2
	default:
		panic(ErrInvalidMetricValue)
	}
}

func generalRemediationLevel(v uint8) float64 {
	switch v {
	case grl_h:
		return 0.4
	case grl_m:
		return 0.6
	case grl_l:
		return 0.8
	case grl_n, grl_nd:
		return 1.0
	default:
		panic(ErrInvalidMetricValue)
	}
}

func localVulnerabilityPrevalence(v uint8) float64 {
	switch v {
	case lvp_n:
		return 0.6
	case lvp_l:
		return 0.8
	case lvp_m, lvp_nd:
		return 1.0
	case lvp_h:
		return 1.2
	default:
		panic(ErrInvalidMetricValue)
	}
}

func perceivedTargetValue(v uint8) float64 {
	switch v {
	case ptv_l:
		return 0.8
	case ptv_m, ptv_nd:
		return 1.0
	case ptv_h:
		return 1.2
	default:
		panic(ErrInvalidMetricValue)
	}
}

func ciar(v uint8) float64 {
	switch v {
	case ciar_l:
		return 0.5
	case ciar_m, ciar_nd:
		return 1.0
	case ciar_h:
		return 1.51
	default:
		panic(ErrInvalidMetricValue)
	}
}

func collateralDamagePotential(v uint8) float64 {
	switch v {
	case cdp_n, cdp_nd:
		return 1.0
	case cdp_l:
		return 1.25
	case cdp_lm:
		return 1.5
	case cdp_mh:
		return 1.75
	case cdp_h:
		return 2.0
	default:
		panic(ErrInvalidMetricValue)
	}
}

func localRemediationLevel(v uint8) float64 {
	switch v {
	case lrl_n, lrl_nd:
		return 1.0
	case lrl_l:
		return 0.8
	case lrl_m:
		return 0.6
	case lrl_h:
		return 0.4
	default:
		panic(ErrInvalidMetricValue)
	}
}

// this helper is not specified, so we literally round the value
// to 1 decimal.
func roundTo1Decimal(x float64) float64 {
	return math.RoundToEven(x*10) / 10
}
