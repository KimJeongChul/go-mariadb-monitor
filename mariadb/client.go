package mariadb

import (
	cerror "github.com/KimJeongChul/go-mariadb-monitor/error"
	"github.com/KimJeongChul/go-mariadb-monitor/logger"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MariaDBClient struct {
	UserName      string
	Password      string
	ServerAddress string
	DbName        string
	db            *gorm.DB
}

func (mc *MariaDBClient) Connect() *cerror.CError {
	dsn := mc.UserName + ":" + mc.Password + "@tcp(" + mc.ServerAddress + ")/" + mc.DbName + "?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	// MariaDB Connect
	mc.db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.LogE("mariadb", "Connect", "error:", err)
		return cerror.NewCError(cerror.MARIA_DB_ERR, "database open error:"+err.Error())
	}
	return nil
}
