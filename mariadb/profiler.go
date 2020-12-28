package mariadb

import (
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
	//funcName := "profiler:Start"
	ticker := time.NewTicker(time.Duration(mp.period) * time.Second)
	for {
		select {
		case <-ticker.C:
			status, cErr := mp.mariaDBClient.getStatus()
			logger.LogI(packageName, "Start", "Status:", status)
			if cErr != nil {
				continue
			}

		}
	}
}
