package qc

import (
	"fmt"

	"lims/internal/result"
)

// Verdict is the outcome of a QC judgement.
type Verdict string

const (
	VerdictPass Verdict = "pass"
	VerdictFail Verdict = "fail"
)

// Judge compares measurements against the reference range.
func Judge(rule RuleSet, values []result.Measurement) (Verdict, string) {
	for _, measurement := range values {
		if measurement.Analyte != rule.Analyte {
			continue
		}
		if measurement.Value < rule.Min || measurement.Value > rule.Max {
			return VerdictFail, fmt.Sprintf(
				"%s %.2f %s outside reference range [%.2f, %.2f]",
				measurement.Analyte, measurement.Value, measurement.Unit, rule.Min, rule.Max,
			)
		}
		return VerdictPass, "measurement within reference range"
	}
	return VerdictFail, fmt.Sprintf("analyte %s was not measured", rule.Analyte)
}
