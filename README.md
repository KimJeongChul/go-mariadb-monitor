# go-mariadb-monitor
MariaDB montior displays information about the MariaDB db stats in real time.

![image](https://user-images.githubusercontent.com/10591350/103218912-ab530e80-495f-11eb-891b-6dc6c9028ebc.png)

![image](https://user-images.githubusercontent.com/10591350/103203305-f22d0e00-4937-11eb-8ce7-9f94be866f95.png)

### Start
```shell
go build
./go-mariadb-monitor
```

### Config file
config.json
```json
{
    "port": {HTTP_SERVER_PORT},
    "period": {MONITORING_PERIOD},
    "mariaDBServerAddr": {MARIA_SERVER_ADDR},
    "mariaDBServerPort": {MARIA_SERVER_PORT},
    "mariaDBUser": {MARIA_DB_USER},
    "mariaDBPassword": {MARIA_DB_PASSWORD},
    "mariaDBName": {MARIA_DB_NAME}
}
```