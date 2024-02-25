package cron

import (
	"errors"
	"github.com/robfig/cron/v3"
	"sync"
)

var c *cron.Cron
var tasks = make(map[string]cron.EntryID)
var lock sync.RWMutex

func init() {
	c = cron.New()
}

func AddTask(name string, cronStr string, f func()) error {
	e, err := c.AddFunc(cronStr, f)
	if err != nil {
		return err
	}
	lock.Lock()
	tasks[name] = e
	lock.Unlock()
	return nil
}

func DelTask(name string) error {
	lock.RLock()
	if e, ok := tasks[name]; ok {
		lock.RUnlock()
		lock.Lock()
		c.Remove(e)
		lock.Unlock()
	} else {
		lock.RUnlock()
		return errors.New("not found the task")
	}
	return nil
}

func Start() error {
	c.Start()
	return nil
}

func Close() error {
	c.Stop()
	return nil
}
