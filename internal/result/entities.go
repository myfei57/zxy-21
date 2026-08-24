package result

import "time"

// Status is the QC lifecycle state of a result.
type Status string

const (
	StatusPending Status = "pending"
	StatusPassed  Status = "passed"
	StatusFailed  Status = "failed"
)

// Measurement is one analyte value produced by an instrument.
type Measurement struct {
	Analyte string  `json:"analyte"`
	Value   float64 `json:"value"`
	Unit    string  `json:"unit"`
}

// Result holds the raw measurements and the QC state of one detection.
type Result struct {
	ID           string        `json:"id"`
	SampleID     string        `json:"sample_id"`
	TaskID       string        `json:"task_id"`
	InstrumentID string        `json:"instrument_id"`
	Status       Status        `json:"status"`
	Values       []Measurement `json:"values"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
}

// QCRecord is the durable evidence of one QC judgement.
type QCRecord struct {
	ResultID    string    `json:"result_id"`
	RuleSetName string    `json:"rule_set_name"`
	Verdict     string    `json:"verdict"`
	Reason      string    `json:"reason"`
	JudgedAt    time.Time `json:"judged_at"`
	Operator    string    `json:"operator"`
}
