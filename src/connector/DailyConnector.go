package connector

import (
	"net/http"
)

const(
	PROJECT_NORMAL_PREFIX = "Sprint"
)

func DailyConnector(w http.ResponseWriter, r *http.Request) {
	connect(w, r, PROJECT_NORMAL_PREFIX, false)
}