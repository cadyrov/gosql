package internal

import (
	"fmt"
	"strings"
)

// nolint:staticcheck
func sqlToGo(name string) string {
	words := strings.Split(name, "_")
	res := ""

	for i := range words {
		w := strings.Title(words[i])

		if strings.EqualFold(w, "id") || strings.EqualFold(w, "url") || strings.EqualFold(w, "sql") ||
			strings.EqualFold(w, "html") || strings.EqualFold(w, "http") || strings.EqualFold(w, "os") ||
			strings.EqualFold(w, "json") {
			w = strings.ToUpper(w)
		}

		res = fmt.Sprintf("%s%s", res, w)
	}

	return res
}
