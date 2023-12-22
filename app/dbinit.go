package app

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"gopkg.in/gorp.v1"
)

type Configuration struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Dbname   string `json:"dbname"`
}

type Order struct {
	Id           int64  `db:"order_id"`
	CustomerName string `db:"customer_name"`
	Status       string `db:"status"`
	CreatedAt    int64  `db:"created_at"`
	UpdatedAt    int64  `db:"updated_at"`
}

type OrderProduct struct {
	Id           int64  `db:"order_product_id"`
	OrderId      int64  `db:"order_id"`
	ProductId    int64  `db:"item_id"`
	ProductCode  string `db:"item_code"`
	CustomerName string `db:"customer_name"`
	CreatedAt    int64  `db:"created_at"`
	UpdatedAt    int64  `db:"updated_at"`
}

type Product struct {
	Id          int64  `db:"product_id" json:"item_id"`
	ProductCode string `db:"item_code"`
	Quantity    int64  `db:"quantity"`
	Description string `db:"description"`
	OrderId     int64  `db:"order_id"`
	CreatedAt   int64  `db:"created_at"`
	UpdatedAt   int64  `db:"updated_at"`
}

func InitDB() {
	dbmap := connectToDB()
	defer dbmap.Db.Close()

	err := dbmap.CreateTablesIfNotExists()
	checkErr(err, "Create tables failed")
}

func connectToDB() *gorp.DbMap {
	config := loadDBConfigs("../../configs/dbconf.json")

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.Dbname)

	db, err := sql.Open("postgres", psqlInfo)
	checkErr(err, "sql.Open failed")
	//defer db.Close()

	err = db.Ping()
	checkErr(err, "Connection to DB or ping failed")

	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}

	dbmap.AddTableWithName(Order{}, "orders").SetKeys(true, "Id")
	dbmap.AddTableWithName(Product{}, "products").SetKeys(true, "Id")
	dbmap.AddTableWithName(OrderProduct{}, "order_products").SetKeys(true, "Id")

	return dbmap
}

func loadDBConfigs(filepath string) Configuration {

	configFile, err := os.Open(filepath)
	defer configFile.Close()
	checkErr(err, "Error reading DB configs from JSON file")
	jsonParser := json.NewDecoder(configFile)
	config := Configuration{}
	jsonParser.Decode(&config)
	return config
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
