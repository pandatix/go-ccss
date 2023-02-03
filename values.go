package goccss

// The following values are the enumerations of each metric possible values.
// Do not change their order as it is vital to the implementation.
//
// Those are used to get a highly memory-performant implementation.
// The 6 bytes are used as follows, using the uint8 type.
//
//    u0       u1       u2       u3       u4       u5
// /------\ /------\ /------\ /------\ /------\ /------\
// ........ ........ ........ ........ ........ ........
// \/\/\/\/ \/\/\/|\__/\_/\_/ \/\_/\/\_/\/\_/\/ \/\/
// AV |Au C  I A |EM | GRL | PTV | EC EI |CDP | IR |
//   AC         PL  GEL   LVP   LRL     EA   CR   AR

// Base

const (
	av_l uint8 = iota
	av_a
	av_n
)

const (
	ac_h uint8 = iota
	ac_m
	ac_l
)

const (
	au_m uint8 = iota
	au_s
	au_n
)

const (
	cia_n uint8 = iota
	cia_p
	cia_c
)

const (
	pl_nd uint8 = iota
	pl_r
	pl_u
	pl_a
)

const (
	em_a uint8 = iota
	em_p
)

// Temporal

const (
	gel_nd uint8 = iota
	gel_n
	gel_l
	gel_m
	gel_h
)

const (
	grl_nd uint8 = iota
	grl_h
	grl_m
	grl_l
	grl_n
)

// Environmental

const (
	lvp_nd uint8 = iota
	lvp_n
	lvp_l
	lvp_m
	lvp_h
)

const (
	ptv_nd uint8 = iota
	ptv_l
	ptv_m
	ptv_h
)

const (
	lrl_nd uint8 = iota
	lrl_n
	lrl_l
	lrl_m
	lrl_h
)

const (
	ecia_nd uint8 = iota
	ecia_n
	ecia_p
	ecia_c
)

const (
	cdp_nd uint8 = iota
	cdp_n
	cdp_l
	cdp_lm
	cdp_mh
	cdp_h
)

const (
	ciar_nd uint8 = iota
	ciar_l
	ciar_m
	ciar_h
)
