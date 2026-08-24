package qc

import (
	"fmt"

	"lims/internal/errs"
)

// RuleSet describes the reference range of one analyte.
type RuleSet struct {
	Name    string  `json:"name"`
	Analyte string  `json:"analyte"`
	Min     float64 `json:"min"`
	Max     float64 `json:"max"`
	Unit    string  `json:"unit"`
}

// DefaultRuleSets returns the built-in QC rules used by the console.
func DefaultRuleSets() []RuleSet {
	return []RuleSet{
		{Name: "glucose-fasting", Analyte: "GLU", Min: 3.9, Max: 6.1, Unit: "mmol/L"},
		{Name: "wbc-count", Analyte: "WBC", Min: 3.5, Max: 9.5, Unit: "10^9/L"},
		{Name: "hgb-hemoglobin", Analyte: "HGB", Min: 115, Max: 175, Unit: "g/L"},
	}
}

// FindRule locates the QC rule for one analyte.
func FindRule(rules []RuleSet, analyte string) (RuleSet, error) {
	for _, rule := range rules {
		if rule.Analyte == analyte {
			return rule, nil
		}
	}
	return RuleSet{}, fmt.Errorf("%w: no QC rule for analyte %s", errs.ErrNotFound, analyte)
}
