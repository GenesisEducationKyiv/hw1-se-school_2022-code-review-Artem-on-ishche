package config

import (
	"os"
)

var (
	CustomersServerURL  = os.Getenv("CUSTOMERS_SERVICE_URL")
	CustomersServerPort = os.Getenv("CUSTOMERS_SERVICE_PORT")
	CreateRoute         = os.Getenv("CREATE_CUSTOMERS_ROUTE")

	DtmCoordinatorAddress = os.Getenv("DTM_COORDINATOR")
	MySqlDsn              = os.Getenv("MYSQL_DSN")
)
