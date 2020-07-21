package model

import (
	"github.com/jochenboesmans/gedcom-parser/proto"
	"sync"
)

type ConcurrencySafeGedcom struct {
	Gedcom proto.Gedcom
	Lock   sync.RWMutex
}
