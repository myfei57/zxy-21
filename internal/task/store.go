package task

// Store persists tasks and their assignment records.
type Store interface {
	SaveTask(task Task) error
	LoadTask(id string) (Task, error)
	ListTasks() ([]Task, error)
	SaveAssignment(record AssignmentRecord) error
	LoadAssignment(taskID string) (AssignmentRecord, error)
	ListAssignments() ([]AssignmentRecord, error)
}
