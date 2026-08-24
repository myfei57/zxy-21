package result

import (
	"errors"
	"fmt"

	"lims/internal/errs"
	"lims/internal/task"
)

// ResultView is what the console page renders for one result.
type ResultView struct {
	Result        Result         `json:"result"`
	QCRecord      *QCRecord      `json:"qc_record,omitempty"`
	TaskState     task.TaskState `json:"task_state"`
	DisplayStatus string         `json:"display_status"`
}

// View reads the live result record so the page never shows a stale verdict.
func View(store Store, taskStore task.Store, resultID string) (ResultView, error) {
	current, err := store.LoadResult(resultID)
	if err != nil {
		return ResultView{}, err
	}
	state, err := task.StateFor(taskStore, store, current.TaskID)
	if err != nil {
		return ResultView{}, fmt.Errorf("load task state: %w", err)
	}
	record, err := store.LoadQCRecord(resultID)
	var qcRecord *QCRecord
	if err == nil {
		qcRecord = &record
	} else if !errors.Is(err, errs.ErrNotFound) {
		return ResultView{}, fmt.Errorf("load QC record: %w", err)
	}
	return ResultView{
		Result:        current,
		QCRecord:      qcRecord,
		TaskState:     state,
		DisplayStatus: state.ResultStatus,
	}, nil
}
