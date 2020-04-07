package periodic

import (
	"errors"
	"github.com/google/uuid"
	"sync"
	"time"
)

type Operator func(data ...interface{})

type Periodic struct {
	lock     sync.Mutex
	capacity int
	tasks    map[string]chan bool
}

func (p *Periodic) Repeat(t time.Duration, o Operator, data ...interface{}) (string, error) {
	p.lock.Lock()
	defer p.lock.Unlock()

	if len(p.tasks) < p.capacity {

		retVal := make(chan bool)
		go func() {
			ticker := time.NewTicker(t)
			for {
				select {
				case <-ticker.C:
					o(data)
				case <-retVal:
					ticker.Stop()
					return
				}
			}
		}()
		id := uuid.New().String()
		p.tasks[id] = retVal

		return id, nil
	}
	return "", errors.New("Max capacity reached")
}

func (p *Periodic) Once(t time.Duration, o Operator, data ...interface{}) (string, error) {
	p.lock.Lock()
	defer p.lock.Unlock()

	if len(p.tasks) < p.capacity {

		retVal := make(chan bool)
		go func() {
			timer := time.NewTimer(t)
			for {
				select {
				case <-timer.C:
					o(data)
				case <-retVal:
					timer.Stop()
					return
				}
			}
		}()
		id := uuid.New().String()
		p.tasks[id] = retVal

		return id, nil
	}

	return "", errors.New("Max capacity reached")
}

func (p *Periodic) Stop(id string) (bool, error) {
	p.lock.Lock()
	defer p.lock.Unlock()

	if taskChan, ok := p.tasks[id]; ok {
		taskChan <- true
		delete(p.tasks, id)
		return true, nil
	}
	return false, errors.New("Task do not exists")
}

func (p *Periodic) Len() int {
	return len(p.tasks)
}

func (p *Periodic) Capacity() int {
	return p.capacity
}

func New(c int) *Periodic {
	return &Periodic{
		capacity: c,
		tasks:    map[string]chan bool{},
	}
}
