package util

import (
	"errors"
	"sync"
)

type Stack struct {
	Lock sync.Mutex
	Elements []string
}

func New() *Stack {
	return &Stack{sync.Mutex{}, []string{}}
}

func (s *Stack) Push (v string) {
	s.Lock.Lock()
	defer s.Lock.Unlock()

	s.Elements = append(s.Elements, v)
}

func (s *Stack) Pop () (string, error) {
	s.Lock.Lock()
	defer s.Lock.Unlock()

	stackSize := len(s.Elements)
	if stackSize == 0 {
		return "", errors.New("stack is empty")
	}

	topElement := s.Elements[stackSize-1]
	s.Elements = s.Elements[:stackSize-1]
	return topElement, nil
}
