package index

import "context"

// Manager allows for managing active index operations conveniently.
// Must be initialised using NewManager at app startup.
type Manager struct {
	operations map[string]*Operation
}

// NewManager intialises a new index manager.
func NewManager() *Manager {
	return &Manager{
		operations: make(map[string]*Operation),
	}
}

// GetOperation fetches the operation with corresponding pid.
func (m *Manager) GetOperation(pid string) (*Operation, bool) {
	o, ok := m.operations[pid]
	return o, ok
}

// CancelOperation cancels the active operation and deletes from active operations map.
// NOTE: does not delete from database or set status to paused.
func (m *Manager) CancelOperation(pid string) bool {
	o, ok := m.GetOperation(pid)
	if !ok {
		return false
	}

	o.cancelFunc()
	delete(m.operations, pid)

	return true
}

// InsertOperation adds an operation to the actove operations map.
func (m *Manager) InsertOperation(o *Operation) {
	m.operations[o.ID] = o
}

// RunOperation starts the index operation from the crrent message in the background.
func (m *Manager) RunOperation(ctx context.Context, o *Operation) {
	m.InsertOperation(o)
	go o.run(ctx)
}
