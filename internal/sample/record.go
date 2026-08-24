package sample

import "time"

// RegistrationRecord is the durable evidence that a sample entered the lab.
type RegistrationRecord struct {
	SampleID     string    `json:"sample_id"`
	NamespaceID  string    `json:"namespace_id"`
	Operator     string    `json:"operator"`
	RegisteredAt time.Time `json:"registered_at"`
	Note         string    `json:"note"`
}

// NewRegistrationRecord builds the registration evidence for one sample.
func NewRegistrationRecord(sampleID, namespaceID, operator, note string, at time.Time) RegistrationRecord {
	return RegistrationRecord{
		SampleID:     sampleID,
		NamespaceID:  namespaceID,
		Operator:     operator,
		RegisteredAt: at,
		Note:         note,
	}
}
