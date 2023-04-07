package query

import (
	"fmt"
	"strings"
)

// Формирует аргументы для prepared запроса,
// fieldsCount: количество полей,
// rowsCount: количество вставляемых строк с данными
func renderValues(fieldsCount, rowsCount int) string {
	var (
		rows    = make([]string, rowsCount)
		row     = make([]string, fieldsCount)
		counter = 1
	)

	for i := 0; i < rowsCount; i++ {
		for j := 0; j < fieldsCount; j++ {
			row[j] = fmt.Sprintf("$%d", j+counter)
		}

		rows[i] = fmt.Sprintf("(%s)", strings.Join(row, ","))
		counter += fieldsCount
	}

	return strings.Join(rows, ",")
}
