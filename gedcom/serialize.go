package gedcom

import (
	"encoding/json"
	"github.com/golang/protobuf/proto"
)

func (gedcom *ConcurrencySafeGedcom) ToJSON() (*[]byte, error) {
	gedcomJson, err := json.Marshal(&gedcom.Gedcom)
	if err != nil {
		return nil, err
	}
	return &gedcomJson, nil
}

func (gedcom *ConcurrencySafeGedcom) ToProto() (*[]byte, error) {
	gedcomProto, err := proto.Marshal(&gedcom.Gedcom)
	if err != nil {
		return nil, err
	}
	return &gedcomProto, nil
}
