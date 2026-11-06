package dal

import (
	"errors"
	"sync"
	"taskmanager/biz/model"
)

// returned when no task exists for a given ID
var ErrTaskNotFound = errors.New("task not found")

// In memory store for tasks ( can change later for json file/real db)
type TaskStore struct {
	mu     sync.RWMutex // Mutex to block concurrent access
	tasks  map[int64]*model.Task
	nextID int64 // auto-increment id
}

// Creates empty ready to use in memory store
func NewTaskStore() *TaskStore {
	return &TaskStore{
		tasks:  (make(map[int64]*model.Task)),
		nextID: 1,
	}
}

// Assigns a new ID, stores the task snd returns it
func (s *TaskStore) Create(t *model.Task) *model.Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	t.ID = s.nextID
	s.nextID++
	s.tasks[t.ID] = t
	return t
}

func (s *TaskStore) GetById(id int64) (*model.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	t, ok := s.tasks[id]

	if !ok {
		return nil, ErrTaskNotFound
	}

	return t, nil
}

func (s *TaskStore) List() []*model.Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]*model.Task, 0, len(s.tasks))
	for _, t := range s.tasks {
		out = append(out, t)
	}
	return out
}

func (s *TaskStore) Update(t *model.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.tasks[t.ID]; !ok {
		return ErrTaskNotFound
	}

	s.tasks[t.ID] = t
	return nil
}

func (s *TaskStore) Delete(id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.tasks[id]; !ok {
		return ErrTaskNotFound
	}

	delete(s.tasks, id)
	return nil
}
