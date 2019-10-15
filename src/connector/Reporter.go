package connector

import (
	"log"
	"golang.org/x/net/context"
)

const RANDOM_REPORTER_GETTER_SQL = "SELECT * FROM `<project>.<data_set>.reporter` ORDER BY RAND() LIMIT 1"

type Reporter struct {
	Name             string `json:"name,omitempty"`
	Icon             string `json:"name,omitempty"`
	Talk             string `json:"name,omitempty"`
}

func getTodayReporter(ctx context.Context) (Reporter, error) {
	it, readErr := readQuery(ctx, RANDOM_REPORTER_GETTER_SQL)
	if readErr != nil {
		return Reporter{}, readErr
	}
	var line Reporter
	err := it.Next(&line)
	if err != nil {
		log.Println("Failed to Iterate Query:", err)
	}
	return line, nil
}
