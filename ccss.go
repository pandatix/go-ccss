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
		base: base{},
		temporal: temporal{
			generalExploitLevel:     "ND",
			generalRemediationLevel: "ND",
		},
		environmental: environmental{
			localVulnerabilityPrevalence: "ND",
			perceivedTargetValue:         "ND",
			localRemediationLevel:        "ND",
			environmentConfidentiality:   "ND",
			environmentIntegrity:         "ND",
			environmentAvailability:      "ND",
			collateralDamagePotential:    "ND",
			confidentialityRequirement:   "ND",
			integrityRequirement:         "ND",
			availabilityRequirement:      "ND",
		},
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

			if curr == 20 {
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
	base
	temporal
	environmental
}

type base struct {
	// AV -> [L,A,N]. Mandatory
	accessVector string
	// AC -> [H,M,L]. Mandatory
	accessComplexity string
	// Au -> [M,S,N]. Mandatory
	authentication string
	// C -> [N,P,C]. Mandatory
	confidentiality string
	// I -> [N,P,C]. Mandatory
	integrity string
	// A -> [N,P,C]. Mandatory
	availability string

	// The following metrics are not used in the equations (Section 3.2)

	// PL -> [R,U,A,ND]. Mandatory
	privilegeLevel string
	// EM -> [A,P]. Mandatory
	exploitationMethod string
}

type temporal struct {
	// GEL -> [N,L,M,H,ND]. Not mandatory
	generalExploitLevel string
	// GRL -> [H,M,L,N,ND]. Not mandatory
	generalRemediationLevel string
}

type environmental struct {
	// LVP -> [N,L,M,H,ND]. Not mandatory
	localVulnerabilityPrevalence string
	// PTV -> [L,M,H,ND]. Not mandatory
	perceivedTargetValue string
	// LRL -> [N,L,M,H,ND]. Not mandatory
	localRemediationLevel string
	// EC -> [N,P,C,ND]. Not mandatory
	environmentConfidentiality string
	// EI -> [N,P,C,ND]. Not mandatory
	environmentIntegrity string
	// EA -> [N,P,C,ND]. Not mandatory
	environmentAvailability string
	// CDP -> [N,L,LM,MH,H,ND]. Not mandatory
	collateralDamagePotential string
	// CR -> [L,M,H,ND]. Not mandatory
	confidentialityRequirement string
	// IR -> [L,M,H,ND]. Not mandatory
	integrityRequirement string
	// AR -> [L,M,H,ND]. Not mandatory
	availabilityRequirement string
}

func (ccss CCSS) Get(abv string) (string, error) {
	switch abv {
	// Base
	case "AV":
		return ccss.accessVector, nil
	case "AC":
		return ccss.accessComplexity, nil
	case "Au":
		return ccss.authentication, nil
	case "C":
		return ccss.confidentiality, nil
	case "I":
		return ccss.integrity, nil
	case "A":
		return ccss.availability, nil
	case "PL":
		return ccss.privilegeLevel, nil
	case "EM":
		return ccss.exploitationMethod, nil

	// Temporal
	case "GEL":
		return ccss.generalExploitLevel, nil
	case "GRL":
		return ccss.generalRemediationLevel, nil

	// Environmental
	case "LVP":
		return ccss.localVulnerabilityPrevalence, nil
	case "PTV":
		return ccss.perceivedTargetValue, nil
	case "LRL":
		return ccss.localRemediationLevel, nil
	case "EC":
		return ccss.environmentConfidentiality, nil
	case "EI":
		return ccss.environmentIntegrity, nil
	case "EA":
		return ccss.environmentAvailability, nil
	case "CDP":
		return ccss.collateralDamagePotential, nil
	case "CR":
		return ccss.confidentialityRequirement, nil
	case "IR":
		return ccss.integrityRequirement, nil
	case "AR":
		return ccss.availabilityRequirement, nil

	default:
		return "", &ErrInvalidMetric{Abv: abv}
	}
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
		if err := validate(value, []string{"L", "A", "N"}); err != nil {
			return err
		}
		ccss.accessVector = value
	case "AC":
		if err := validate(value, []string{"H", "M", "L"}); err != nil {
			return err
		}
		ccss.accessComplexity = value
	case "Au":
		if err := validate(value, []string{"M", "S", "N"}); err != nil {
			return err
		}
		ccss.authentication = value
	case "C":
		if err := validate(value, []string{"N", "P", "C"}); err != nil {
			return err
		}
		ccss.confidentiality = value
	case "I":
		if err := validate(value, []string{"N", "P", "C"}); err != nil {
			return err
		}
		ccss.integrity = value
	case "A":
		if err := validate(value, []string{"N", "P", "C"}); err != nil {
			return err
		}
		ccss.availability = value
	case "PL":
		if err := validate(value, []string{"R", "U", "A", "ND"}); err != nil {
			return err
		}
		ccss.privilegeLevel = value
	case "EM":
		if err := validate(value, []string{"A", "P"}); err != nil {
			return err
		}
		ccss.exploitationMethod = value

	// Temporal
	case "GEL":
		if err := validate(value, []string{"N", "L", "M", "H", "ND"}); err != nil {
			return err
		}
		ccss.generalExploitLevel = value
	case "GRL":
		if err := validate(value, []string{"H", "M", "L", "N", "ND"}); err != nil {
			return err
		}
		ccss.generalRemediationLevel = value

	// Environmental
	case "LVP":
		if err := validate(value, []string{"N", "L", "M", "H", "ND"}); err != nil {
			return err
		}
		ccss.localVulnerabilityPrevalence = value
	case "PTV":
		if err := validate(value, []string{"L", "M", "H", "ND"}); err != nil {
			return err
		}
		ccss.perceivedTargetValue = value
	case "LRL":
		if err := validate(value, []string{"N", "L", "M", "H", "ND"}); err != nil {
			return err
		}
		ccss.localRemediationLevel = value
	case "EC":
		if err := validate(value, []string{"N", "P", "C", "ND"}); err != nil {
			return err
		}
		ccss.environmentConfidentiality = value
	case "EI":
		if err := validate(value, []string{"N", "P", "C", "ND"}); err != nil {
			return err
		}
		ccss.environmentIntegrity = value
	case "EA":
		if err := validate(value, []string{"N", "P", "C", "ND"}); err != nil {
			return err
		}
		ccss.environmentAvailability = value
	case "CDP":
		if err := validate(value, []string{"N", "L", "LM", "MH", "H", "ND"}); err != nil {
			return err
		}
		ccss.collateralDamagePotential = value
	case "CR":
		if err := validate(value, []string{"L", "M", "H", "ND"}); err != nil {
			return err
		}
		ccss.confidentialityRequirement = value
	case "IR":
		if err := validate(value, []string{"L", "M", "H", "ND"}); err != nil {
			return err
		}
		ccss.integrityRequirement = value
	case "AR":
		if err := validate(value, []string{"L", "M", "H", "ND"}); err != nil {
			return err
		}
		ccss.availabilityRequirement = value

	default:
		return &ErrInvalidMetric{Abv: abv}
	}
	return nil
}

func validate(value string, enabled []string) error {
	// Check is valid
	for _, enbl := range enabled {
		if value == enbl {
			return nil
		}
	}
	return ErrInvalidMetricValue
}

func (ccss CCSS) BaseScore() float64 {
	impact := ccss.Impact()
	return roundTo1Decimal(((0.6 * impact) + (0.4 * ccss.Exploitability()) - 1.5) * fImpact(impact))
}

func (ccss CCSS) Impact() float64 {
	return 10.41 * (1 - (1-cia(ccss.base.confidentiality))*(1-cia(ccss.base.integrity))*(1-cia(ccss.base.availability)))
}

func (ccss CCSS) Exploitability() float64 {
	return 20 * accessVector(ccss.base.accessVector) * authentication(ccss.base.authentication) * accessComplexity(ccss.base.accessComplexity)
}

func (ccss CCSS) TemporalScore() float64 {
	tmpExplt := math.Min(10, ccss.Exploitability()*generalExploitLevel(ccss.temporal.generalExploitLevel)*generalRemediationLevel(ccss.temporal.generalRemediationLevel))
	return roundTo1Decimal(((0.6 * ccss.Impact()) + (0.4 * tmpExplt) - 1.5) * fImpact(ccss.Impact()))
}

func (ccss CCSS) EnvironmentalScore() float64 {
	c := mod(ccss.base.confidentiality, ccss.environmental.environmentConfidentiality)
	i := mod(ccss.base.integrity, ccss.environmental.environmentIntegrity)
	a := mod(ccss.base.availability, ccss.environmental.environmentAvailability)

	lclExpltLvl := localVulnerabilityPrevalence(ccss.environmental.localVulnerabilityPrevalence) * perceivedTargetValue(ccss.environmental.perceivedTargetValue)
	envImpact := math.Min(10, 10.41*(1-(1-cia(c)*ciar(ccss.environmental.confidentialityRequirement))*(1-cia(i)*ciar(ccss.environmental.integrityRequirement))*(1-cia(a)*ciar(ccss.environmental.availabilityRequirement)))*collateralDamagePotential(ccss.environmental.collateralDamagePotential))
	envExplt := math.Min(10, ccss.Exploitability()*generalExploitLevel(ccss.temporal.generalExploitLevel)*lclExpltLvl*localRemediationLevel(ccss.environmental.localRemediationLevel))
	return roundTo1Decimal(((0.6 * envImpact) + (0.4 * envExplt) - 1.5) * fImpact(ccss.Impact()))
}

// Helpers to compute CCSS scores.

func mod(base, env string) string {
	if env != "ND" {
		return env
	}
	return base
}

func fImpact(v float64) float64 {
	if v == 0 {
		return 0
	}
	return 1.176
}

func accessVector(v string) float64 {
	switch v {
	case "L":
		return 0.395
	case "A":
		return 0.646
	case "N":
		return 1.0
	default:
		panic(ErrInvalidMetricValue)
	}
}

func authentication(v string) float64 {
	switch v {
	case "M":
		return 0.45
	case "S":
		return 0.56
	case "N":
		return 0.704
	default:
		panic(ErrInvalidMetricValue)
	}
}

func accessComplexity(v string) float64 {
	switch v {
	case "H":
		return 0.35
	case "M":
		return 0.61
	case "L":
		return 0.71
	default:
		panic(ErrInvalidMetricValue)
	}
}

func cia(v string) float64 {
	switch v {
	case "N":
		return 0.0
	case "P":
		return 0.275
	case "C":
		return 0.660
	default:
		panic(ErrInvalidMetricValue)
	}
}

func generalExploitLevel(v string) float64 {
	switch v {
	case "N":
		return 0.6
	case "L":
		return 0.8
	case "M":
		return 1.0
	case "H":
		return 1.2
	case "ND":
		return 1.0
	default:
		panic(ErrInvalidMetricValue)
	}
}

func generalRemediationLevel(v string) float64 {
	switch v {
	case "H":
		return 0.4
	case "M":
		return 0.6
	case "L":
		return 0.8
	case "N":
		return 1.0
	case "ND":
		return 1.0
	default:
		panic(ErrInvalidMetricValue)
	}
}

func localVulnerabilityPrevalence(v string) float64 {
	switch v {
	case "N":
		return 0.6
	case "L":
		return 0.8
	case "M":
		return 1.0
	case "H":
		return 1.2
	case "ND":
		return 1.0
	default:
		panic(ErrInvalidMetricValue)
	}
}

func perceivedTargetValue(v string) float64 {
	switch v {
	case "L":
		return 0.8
	case "M":
		return 1.0
	case "H":
		return 1.2
	case "ND":
		return 1.0
	default:
		panic(ErrInvalidMetricValue)
	}
}

func ciar(v string) float64 {
	switch v {
	case "L":
		return 0.5
	case "M":
		return 1.0
	case "H":
		return 1.51
	case "ND":
		return 1.0
	default:
		panic(ErrInvalidMetricValue)
	}
}

func collateralDamagePotential(v string) float64 {
	switch v {
	case "N":
		return 1.0
	case "L":
		return 1.25
	case "LM":
		return 1.5
	case "MH":
		return 1.75
	case "H":
		return 2.0
	case "ND":
		return 1.0
	default:
		panic(ErrInvalidMetricValue)
	}
}

func localRemediationLevel(v string) float64 {
	switch v {
	case "N":
		return 1.0
	case "L":
		return 0.8
	case "M":
		return 0.6
	case "H":
		return 0.4
	case "ND":
		return 1.0
	default:
		panic(ErrInvalidMetricValue)
	}
}

// this helper is not specified, so we literally round the value
// to 1 decimal.
func roundTo1Decimal(x float64) float64 {
	return math.RoundToEven(x*10) / 10
}
