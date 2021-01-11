package gedcom

import (
	"fmt"
	"github.com/jochenboesmans/gedcom-parser/util"
)

func interpretSexStructure(line *Line) (string, error) {
	genderFull, ok := util.GenderFullByLetter[line.Value()]
	if !ok {
		return "", fmt.Errorf("invalid sex letter value: %s", line.Value())
	}
	return genderFull, nil
}
