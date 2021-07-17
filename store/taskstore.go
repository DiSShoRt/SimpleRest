package taskstore

import (
	"fmt"
	"sync"
	"time"
)


type Task struct {
	ID   int       `json:"id"`
	Text string    `json:"text"`
	Tags []string  `json:"tags"`
	Due  time.Time `json:"due"`

}


type TaskStore struct {
	mux  sync.Mutex
	tasks map[int]Task
	nextID int
}

func New() *TaskStore {
	ts := &TaskStore{}
	ts.tasks = make(map[int]Task)
	ts.nextID = 0
	return ts
}

func (ts *TaskStore) CreateTask(tx string, tags []string, due time.Time) int {
	ts.mux.Lock()

	task := Task{
		ID: ts.nextID,
		Text: tx,
		Due: due,
	}

	task.Tags = make([]string, len(tags))
	copy(task.Tags, tags)
	defer ts.mux.Unlock()
	ts.nextID++
	return task.ID
}


func (ts *TaskStore) GetTask(id int) (Task, error) {
	ts.mux.Lock()
	defer ts.mux.Unlock()

	t, ok := ts.tasks[id]
	if ok {
		return t, nil
	} else { 
		return Task{}, fmt.Errorf("Please change input id = %d, task not found", id)
	}
}


func (ts *TaskStore) DeleteTask(id int) error {
	ts.mux.Lock()
	defer ts.mux.Unlock()

	if _, ok := ts.tasks[id]; !ok {
		return fmt.Errorf("Please change input id = %d, task not found", id)

	} else { 

		delete(ts.tasks, id)
		return nil
	}
}


func (ts *TaskStore) DeleteAllTask() error {
	ts.mux.Lock()
	defer ts.mux.Unlock()

	ts.tasks = make(map[int]Task)
	return nil
}


func (ts *TaskStore) GetAllTask() []Task {
	ts.mux.Lock()
	defer ts.mux.Unlock()

	all := make([]Task, 0, len(ts.tasks))
	for _, task := range ts.tasks {
		all = append(all, task)
	}
	return all
}


func (ts *TaskStore) GetTaskByTags(tag string) []Task {
	ts.mux.Lock()
	defer ts.mux.Unlock()

	var tasks []Task

	for _,  task := range ts.tasks {
		for _, tasktag := range task.Tags {
			if tasktag == tag {
				tasks = append(tasks, task)
				break
			}
		}
	}
	return tasks
}


func (ts *TaskStore) GetTaskByDue(year int, mn time.Month, day int) []Task {
	ts.mux.Lock()
	defer ts.mux.Unlock()

	var tasks []Task

	for _, task := range ts.tasks {
		y, m, d := task.Due.Date()
		if y == year && m == mn && d == day {
			tasks = append(tasks, task)
		}
	}
	return tasks
}