package utils

import (
	"strconv"
	"strings"
)

func PrepareVarPSQL(statement string) (prepared string) {
	pieces := strings.Split(statement, "?")
	for i := range pieces {
		if i < (len(pieces) - 1) {
			pieces[i] += "$" + strconv.Itoa(i+1)
		}
	}

	prepared = strings.Join(pieces, "")

	return
}
