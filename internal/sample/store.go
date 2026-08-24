package sample

// Store persists samples and their registration records.
type Store interface {
	SaveSample(sample Sample) error
	LoadSample(id string) (Sample, error)
	ListSamples() ([]Sample, error)
	SaveRegistration(record RegistrationRecord) error
	LoadRegistration(sampleID string) (RegistrationRecord, error)
	ListRegistrations() ([]RegistrationRecord, error)
}
