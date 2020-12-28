package mariadb

import (
	"encoding/json"
	"time"

	"github.com/KimJeongChul/go-mariadb-monitor/broker"
	"github.com/KimJeongChul/go-mariadb-monitor/logger"
)

const packageName = "mariadb"

type MariaDBProfiler struct {
	period        int
	mariaDBClient *MariaDBClient // Redis client
	broker        *broker.Broker
}

// NewMariaDBProfiler Create new MariaDB Profiler
func NewMariaDBProfiler(period int, mariaDBClient *MariaDBClient, broker *broker.Broker) *MariaDBProfiler {
	mp := &MariaDBProfiler{
		period:        period,
		mariaDBClient: mariaDBClient,
		broker:        broker,
	}
	return mp
}

func (mp MariaDBProfiler) Start() {
	ticker := time.NewTicker(time.Duration(mp.period) * time.Second)
	for {
		select {
		case <-ticker.C:
			status, cErr := mp.mariaDBClient.getStatus()
			logger.LogI(packageName, "MariaDBProfiler:Start", "Status:", status)
			if cErr != nil {
				logger.LogE(packageName, "MariaDBProfiler:Start", "getStatus error:", cErr.Error())
				continue
			}
			msgUpdateMariaDBStats, err := json.Marshal(map[string]string{
				"method":                       "updateMariaDBStats",
				"threadsConnected":             status.ThreadsConnected,
				"threadsRunning":               status.ThreadsRunning,
				"connections":                  status.Connections,
				"innodbBufferPoolReadRequests": status.InnodbBufferPoolReadRequests,
				"innodbBufferPoolReads":        status.InnodbBufferPoolReads,
				"innodbRowLockWaits":           status.InnodbRowLockWaits,
				"innodbRowLockCurrentWaits":    status.InnodbRowLockCurrentWaits,
				"innodbRowLockTime":            status.InnodbRowLockTime,
				"memoryUsed":                   status.MemoryUsed,
				"bytesReceived":                status.BytesReceived,
				"bytesSent":                    status.BytesSent,
			})
			if err == nil {
				mp.broker.Messages <- msgUpdateMariaDBStats
			} else {
				logger.LogE(packageName, "MariaDBProfiler:Start", "json marshal error:", err.Error())
			}
		}
	}
}
