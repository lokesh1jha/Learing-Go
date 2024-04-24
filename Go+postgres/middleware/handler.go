package middleware

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/lokesh1jha/stocksAPI/models"
)

type Response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

func CreateConnection() *sql.DB {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	db, err := sql.Open("mysql", os.Getenv("POSTGRES_URL"))
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected to DB!")

	return db
}

func CreateStock(w http.ResponseWriter, r *http.Request) {
	var stock models.Stock
	err := json.NewDecoder(r.Body).Decode(&stock)
	if err != nil {
		log.Fatal("Unable to decode JSON")
	}

	insertID := insertStock(stock)
	res := Response{
		ID:      insertID,
		Message: "Stock created successfully",
	}
	json.NewEncoder(w).Encode(res)
}

func GetStock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert string to int. %v", err)
	}

	stock, err := getStock(int64(id))

	if err != nil {
		log.Fatalf("Unable to get stock. %v", err)
	}

	json.NewEncoder(w).Encode(stock)
}

func GetAllStock(w http.ResponseWriter, r *http.Request) {
	stocks, err := getAllStocks()

	if err != nil {
		log.Fatalf("Unable to get all stocks. %v", err)
	}

	json.NewEncoder(w).Encode(stocks)

}

func UpdateStock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert string to int. %v", err)
	}

	var stock models.Stock

	err = json.NewDecoder(r.Body).Decode(&stock)

	if err != nil {
		log.Fatalf("Unable to decode Request Body. %v", err)
	}

	updatedRows := updateStock(int64(id), stock)

	// Create a message
	msg := fmt.Sprintf("Stock updated successfully. Total rows/record affected %v", updatedRows)

	res := Response{
		ID:      int64(id),
		Message: msg,
	}
	json.NewEncoder(w).Encode(res)

}

func DeleteStock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert string to int. %v", err)
	}

	deletedRows := deleteStock(int64(id))

	msg := fmt.Sprintf("Stock deleted successfully. Total rows/record affected %v", deletedRows)

	res := Response{
		ID:      int64(id),
		Message: msg,
	}

	json.NewEncoder(w).Encode(res)
}

func insertStock(stock models.Stock) int64 {
	db := CreateConnection()
	defer db.Close() // at the end of function close connection

	sqlStatement := `INSERT INTO stock(name, price, company) VALUES ($1, $2, $3) RETURNING stockid`
	var id int64

	err := db.QueryRow(sqlStatement, stock.Name, stock.Price, stock.Company).Scan(&id)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	fmt.Printf("Inserted a single record %v", id)
	return id
}

func getStock(id int64) (models.Stock, error) {
	db := CreateConnection()
	defer db.Close()

	var stock models.Stock

	sqlStatement := `SELECT * FROM stock WHERE stockid=$1`
	row := db.QueryRow(sqlStatement, id)

	err := row.Scan(&stock.StockID, &stock.Name, &stock.Price, &stock.Company)

	switch err {
	case sql.ErrNoRows:
		fmt.Printf("Unable to get stock. %v", err)
		return stock, nil
	case nil:
		return stock, nil
	default:
		log.Fatalf("Unable to scan the row. %v", err)
	}
	return stock, err

}

func getAllStocks() ([]models.Stock, error) {
	db := CreateConnection()
	defer db.Close()

	var stocks []models.Stock
	sqlStatement := `SELECT * FROM stock`
	rows, err := db.Query(sqlStatement)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var stock models.Stock
		err := rows.Scan(&stock.StockID, &stock.Name, &stock.Price, &stock.Company)

		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
		}
		stocks = append(stocks, stock)
	}
	return stocks, err
}

func updateStock(id int64, stock models.Stock) int64 {

	db := CreateConnection()
	defer db.Close()

	sqlStatement := `UPDATE stock SET name=$2, price=$3, company=$4 WHERE stockid=$1`

	res, err := db.Exec(sqlStatement, id, stock.Name, stock.Price, stock.Company)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}

	fmt.Printf("Total rows/record affected %v", rowsAffected)

	return rowsAffected
}

func deleteStock(id int64) int64 {

	db := CreateConnection()
	defer db.Close()

	sqlStatement := `DELETE FROM stock WHERE stockid=$1`

	res, err := db.Exec(sqlStatement, id)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}

	fmt.Printf("Total rows/record affected %v", rowsAffected)

	return rowsAffected
}
