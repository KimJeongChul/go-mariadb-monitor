package mariadb

import (
	cerror "github.com/KimJeongChul/go-mariadb-monitor/error"
	"github.com/KimJeongChul/go-mariadb-monitor/logger"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MariaDBClient struct {
	ServerAddress string
	ServerPort    string
	UserName      string
	Password      string
	DbName        string
	db            *gorm.DB
}

func NewMariaDBClient(addr string, port string, user string, password string, dbName string) (*MariaDBClient, *cerror.CError) {
	var err error
	mc := &MariaDBClient{
		ServerAddress: addr,
		ServerPort:    port,
		UserName:      user,
		Password:      password,
		DbName:        dbName,
	}

	dsn := mc.UserName + ":" + mc.Password + "@tcp(" + mc.ServerAddress + ":" + mc.ServerPort + ")/" + mc.DbName + "?charset=utf8mb4&parseTime=True&loc=Local"
	logger.LogI("mariadb", "NewMariaDbClient", "dsn:", dsn)

	// MariaDB Connect
	mc.db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.LogE("mariadb", "NewMariaDbClient", "gorm open error:", err)
		return nil, cerror.NewCError(cerror.MARIA_DB_ERR, "database open error:"+err.Error())
	}
	return mc, nil
}

type Status struct {
	ThreadsConnected             string
	ThreadsRunning               string
	Connections                  string
	InnodbBufferPoolReadRequests string
	InnodbBufferPoolReads        string
	InnodbRowLockWaits           string
	InnodbRowLockCurrentWaits    string
	InnodbRowLockTime            string
	MemoryUsed                   string
	BytesReceived                string
	BytesSent                    string
}

// Schema: status
type Table struct {
	Variable_name string
	Value         string
}

func (mc MariaDBClient) getStatus() (*Status, *cerror.CError) {
	status := &Status{}
	rows, err := mc.db.Raw("SHOW STATUS").Rows()
	defer rows.Close()
	if err != nil {
		cErr := cerror.NewCError(cerror.MARIA_DB_ERR, "show status err:"+err.Error())
		return nil, cErr
	}
	var result Table
	for rows.Next() {
		mc.db.ScanRows(rows, &result)
		logger.LogI("mariadb", "getStatus", "Column Name:", result.Variable_name, "Colume Value:", result.Value)
		switch result.Variable_name {
		case "Threads_connected":
			status.ThreadsConnected = result.Value
		case "Threads_running":
			status.ThreadsRunning = result.Value
		case "Connections":
			status.Connections = result.Value
		case "Innodb_buffer_pool_read_requests":
			status.InnodbBufferPoolReadRequests = result.Value
		case "Innodb_buffer_pool_reads":
			status.InnodbBufferPoolReads = result.Value
		case "Innodb_row_lock_waits":
			status.InnodbRowLockWaits = result.Value
		case "Innodb_row_lock_current_waits":
			status.InnodbRowLockCurrentWaits = result.Value
		case "Innodb_row_lock_time":
			status.InnodbRowLockTime = result.Value
		case "Memory_used":
			status.MemoryUsed = result.Value
		case "Bytes_received":
			status.BytesReceived = result.Value
		case "Bytes_sent":
			status.BytesSent = result.Value
		}
	}
	return status, nil
}
