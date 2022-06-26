package internal

import (
	"database/sql"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

func getSuspendedSessions(conn *sql.DB) []prometheus.Metric {
	var metrics []prometheus.Metric

	rows := performQuery(
		"SELECT wait_type, COUNT(*) AS cnt FROM sys.dm_exec_requests WHERE session_id > 50 AND status = 'suspended' GROUP BY wait_type;",
		conn,
	)

	for rows.Next() {
		var waitTypes string
		var count int64
		err := rows.Scan(&waitTypes, &count)
		if err != nil {
			logrus.Errorf("Failed to scan with error: %s", err)
		}
		metrics = append(metrics, returnMetric(
			"sql_suspended_sessions",
			"Current suspended user sessions",
			"wait_type",
			waitTypes,
			float64(count),
		))
	}
	err := rows.Err()
	if err != nil {
		logrus.Errorf("Scan failed %s:", err)
	}
	return metrics
}
