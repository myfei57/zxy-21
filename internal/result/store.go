package result

// Store persists results, QC records and station acknowledgements.
type Store interface {
	SaveResult(result Result) error
	LoadResult(id string) (Result, error)
	LoadResultStatus(resultID string) (string, error)
	ListResults() ([]Result, error)
	SaveQCRecord(record QCRecord) error
	LoadQCRecord(resultID string) (QCRecord, error)
	ListQCRecords() ([]QCRecord, error)
	SaveAck(ack AckRecord) error
	LoadAck(sampleID string) (AckRecord, error)
	ListAcks() ([]AckRecord, error)
}
