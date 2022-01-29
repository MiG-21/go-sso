package event

import (
	"fmt"
	"sync"

	"github.com/rs/zerolog"
)

type (
	Event interface {
		Name() string
		Data() interface{}
	}

	Listener func(interface{}) error

	Service struct {
		mutex     *sync.Mutex
		quit      chan bool
		events    chan Event
		listeners map[string][]Listener
		logger    *zerolog.Logger
	}
)

func (s *Service) Listen() error {
	// check for mutex
	if s.mutex == nil {
		s.mutex = &sync.Mutex{}
	}

	if s.events != nil {
		return fmt.Errorf("listener already inititated")
	}

	// create the observer channels.
	s.quit = make(chan bool)
	s.events = make(chan Event)

	// run the observer.
	return s.eventLoop()
}

func (s *Service) Shutdown() error {
	// shutdown event loop
	if s.events != nil {
		// send a quit signal.
		s.quit <- true

		// shutdown channels.
		close(s.quit)
		close(s.events)

		s.listeners = nil
	}

	return nil
}

func (s *Service) AddListener(n string, l Listener) {
	// check for mutex
	if s.mutex == nil {
		s.mutex = &sync.Mutex{}
	}

	// Lock:
	// 1. operations on array listeners
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.listeners == nil {
		s.listeners = make(map[string][]Listener)
	}

	s.listeners[n] = append(s.listeners[n], l)
}

func (s *Service) Emit(event Event) {
	s.events <- event
}

// handleEvent send event to the observer listeners.
func (s *Service) handleEvent(event Event) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if listeners, ok := s.listeners[event.Name()]; ok {
		for _, listener := range listeners {
			go func(l Listener) {
				if err := l(event.Data()); err != nil {
					s.logger.Err(err)
				}
			}(listener)
		}
	}
}

func (s *Service) eventLoop() error {
	// run observer.
	go func() {
		for {
			select {
			case event := <-s.events:
				s.handleEvent(event)
			case <-s.quit:
				return
			}
		}
	}()

	return nil
}
