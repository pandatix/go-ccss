# Go-CCSS

[![reference](https://godoc.org/github.com/pandatix/go-ccss/v5?status.svg=)](https://pkg.go.dev/github.com/pandatix/go-ccss)
[![go report](https://goreportcard.com/badge/github.com/pandatix/go-ccss)](https://goreportcard.com/report/github.com/pandatix/go-ccss)
[![Coverage Status](https://coveralls.io/repos/github/pandatix/go-ccss/badge.svg?branch=main)](https://coveralls.io/github/pandatix/go-ccss?branch=main)
[![CI](https://github.com/pandatix/go-ccss/actions/workflows/ci.yaml/badge.svg)](https://github.com/pandatix/go-ccss/actions?query=workflow%3Aci+)
[![CodeQL](https://github.com/pandatix/go-ccss/actions/workflows/codeql-analysis.yaml/badge.svg)](https://github.com/pandatix/go-ccss/actions/workflows/codeql-analysis.yaml)

Go-CCSS is a low-allocation Go module made to manipulate Common Configuration Scoring System (CCSS), standardized by the [NIST-IR 7502](https://csrc.nist.gov/publications/detail/nistir/7502/final). It is the only implementation known today.

The CCSS is a set of measures of the severity of software security configuration issues.

## Summary

 - [How to use](#how-to-use)
 - [Feedbacks](#feedbacks)

## How to use

The following code gives an example on how to use the present Go module.

It parses a CCSS vector, then computes its base score and displays it.

```go
package main

import (
	"fmt"
	"log"

	goccss "github.com/pandatix/go-ccss"
)

func main() {
	vec := "AV:N/AC:L/Au:N/C:P/I:P/A:P/PL:ND/EM:P"
	ccss, err := goccss.ParseVector(vec)
	if err != nil {
		log.Fatal(err)
	}
	baseScore := ccss.BaseScore()
	fmt.Printf("%s -> base score: %.1f\n", vec, baseScore)
}
```

## Feedbacks

The idea behind CCSS is to score configuration issues that lies between software flaws, covered by CVSS, and misuse vulnerabilities, covered by CMSS. This is a pretty good concept, but don't seem to be used a lot effectively. This is coroborated by the empty responses to Github queries (like `"Common Configuration Scoring System"`, `"CCSS"` and `"NIST-IR 7502`) and the poor discussions in the scientific field.

Here follows a set of feedbacks on the NIST-IR 7502, inherited from a reader and developer point-of-views:
 - the CCSS formulas suffers from the same issues than CVSS: the linear regression can't provide explainability, there is no details on how they were created so the methodology can't by criticised... See [1] for more infos. Moreover, in Section 1.2, 3rd paragraph, it is stated that "an organization could use vulnerability measurements as part of determining the relative importance of particular settings and identifying the settings causing the greatest increase in risk". Nevertheless, as the score can't be explained due to the regression that constitutes the formulas, those measurements are not exploitable. It's pointed out that they are indicators and small differences don't change everything. The specification also asserts that "writing CCSS data is not yet directly useful to organizations for decision making" which raises questions on its usability and the real value of it. The CCSS is still a work to be done, as it needs additional efforts in order to be usable by industry and researchers. Finally, adjusting attributes of a CCSS vector to decrease a score is a poor methodology that should not be used. It should be prefered a methodology where considering the assets and possible threats, as in [2]. It shows that a score could enable comparisons and decisions, but security can't rely only on it as the assets and their threats are the core or the problem and have to get the focus.
 - the EM (Exploitation Method) of Base Metrics is not part of the score. This goes on the opposite of the scoring idea, as a metric value should affect the score in order to categorize each combination, for comparison purposes as it is when involved in decision making.
 - the Au (Authentication) of Base Metrics only "measures the number of times an exploiter must authenticate to a target in
order to exploit a vulnerability", but does not "gauge the strength or complexity of the authentication process". This does not reflect if the authentication process is weak, as a poor 2FA would be considered stronger than a password but could effectively be the opposite. Indeed, if the 2FA choices make that an attacker could intercept and pass it (using SMS, a key...), it could be considered weaker that a single password of 4k characters on a 64-chars alphabet (64^4k combinations).
Another example could be a 2FA with password and e-mail. In case the password is shared (frequent case), then the attacker only needs it to bypass the 2FA, so the authentication complexity only rely on the password which is a single authentication vector.
 - the GEL (General Exploit Level) of Temporal Metrics is based on the assumption that someone could know at any time how often an exploit exist and is used. This kind of information is hard to get, so the metric is trusted hard to effectively emit : exploits could not be made public, informations on attacks may be missing...
Moreover, it is part of the environmental score computing, so variations of it will change not only temporal metrics but also environmental, which is against the concept of independent metric groups. This imply that with a vector composed of base and temporal metrics, the environmental will have a different score than the base, which behaves differently than CVSS v2 (obviously serves as a foundation for this specification).
 - the GRL (General Remediation Level) of Temporal Metrics has a value depending on the decrease of the incidence of exploitation, while this is trusted hard to measure. It is the same problem that lies under CVSS/CMSS/CCSS : how can you give a quantitative result (a score) from qualitative elements (an enumeration) while having explainability ? Indeed, how can you quantitatively measure the impact of a qualitative decision (deploy a firewall, train people...) effectively and in a short time (to be usable) ? The answers to those questions are trusted not trivial, despite the specification words. This also applies to LVP (Local Vulnerability Prevalence) and TRL (Top Remediation Level).
 - the CDP (Collateral Damage Potential) of Local Impact Metrics defines an enumeration that is left to the appreciation of the user. Those values are later interpreted as numerical values by the specification. This implies that the output are non-reproductible and not explainable between environments as the specification does not give a framework in which a user have to work, but rather depends on the quality of every operators evaluations. For instance, a company can decide the High value must never be used, removing a whole set of combinations and so a score. Nevertheless, this is leveraged by the fact that the Local Impact is local to an environment. It remains inconsistent by design between groups as comparison is not possible.
 - section 2.3.3.3 states that the CCSS does not enforce a methodology as the FIPS 199, which is a double-edged decision. First, it makes the organization free of defining its own, which is good to match its use cases and existing methodology if any. On the other side, it does not give one to the newcomers, or even guide them deeply. The CCSS specification could give a basic approach that could then be conjugated to specific needs and tactics.
 - the specifications illustrates an analyst could generate multiple vectors depending on a configuration. This is indeed usefull to enumerate cases, but it won't work sooner or later by the number of combinations.
 - section 3.1 paragraph 3 refers to a future work that never seems to have been published.
 - section 2.1.1.1 Exploitation Method is the first Base metric defined in the specifiation, despite being the last one in the vector. This section could be put at the end.
 - the Base metric PL (Privilege Level) is not defined in its sections, but is defined in Table 14 and the vector samples.
 - there is not rating defined as for CVSS (NONE, LOW, MEDIUM, HIGH, CRITICAL). As for this last, it could enable to transcribe promptly the severity.
 - section 2.1.2 in its last paragraph defines the Base metric Privilege Level (PL), and states that it is used by Environmental metrics as defined in Section 2.3.3.1. Nevertheless, this last section does not refers to this metric. In addition, it is also not present in the equations defined by Section 3.2.
 - vector samples are not enoughly numerous. This can't enable a good understanding of the understanding and the mandatory caracteristic of attributes. Moreover, this impacts the code coverage as it does not provide enough test cases.
 - CCE-2776-3 3rd sample does not explicitly state it needs the temporal metrics defined before, which is necessary to produce the expected result. The whole vector should be proposed in order to better understand this.

References:
 - [1] J. Spring, E. Hatleback, A. Householder, A. Manion and D. Shick, "Time to Change the CVSS?," in IEEE Security & Privacy, vol. 19, no. 2, pp. 74-78, March-April 2021, doi: 10.1109/MSEC.2020.3044475.
 - [2] Nan Messe (Zhang), "Security by Design : An asset-based approach to bridge the gap between architects and security experts", 2021, HAL Id: tel-03407189.
