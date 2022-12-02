package Utils

import "strings"

// SplitQueries is a simple function that splits a text string of queries typically from a .sql file and returns a
// slice of formatted queries that has had the "\n" and spaces removed from each end. It does it's initial split by
// using the ";" character, so have each query end with it.
func SplitQueries(queries string) (formattedQueries []string) {
	splitQueries := strings.Split(queries, ";")

	for _, query := range splitQueries {
		if query == "" {
			continue
		}
		query = strings.Trim(query, "\n")
		query = strings.TrimSpace(query)
		formattedQueries = append(formattedQueries, query+";")
	}

	return
}
