package instrument

// Store persists instruments and their dispatch records.
type Store interface {
	SaveInstrument(instrument Instrument) error
	LoadInstrument(id string) (Instrument, error)
	ListInstruments() ([]Instrument, error)
	SaveDispatch(assignment InstrumentAssignment) error
	HasAssignment(taskID string) (bool, error)
	ListDispatches() ([]InstrumentAssignment, error)
}
