package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Payment struct {
	ID         string `json:"id"`
	FromUserID string `json:"from_user_id"`
	ToUserID   string `json:"to_user_id"`
	Amount     int    `json:"amount"`
}

var payments = []Payment{
	{ID: "1", FromUserID: "1", ToUserID: "2", Amount: 100},
	{ID: "2", FromUserID: "2", ToUserID: "3", Amount: 200},
	{ID: "3", FromUserID: "3", ToUserID: "4", Amount: 300},
}

func main() {
	dsn := "postgres://squadbot:p4ssw0rd@localhost:5433/squadbot?sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	rows, err := db.Query("SELECT id, name FROM users")
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var (
			id   string
			name string
		)

		err := rows.Scan(&id, &name)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(id, name)
	}

	router := gin.Default()
	router.GET("/payments", getPayments)
	router.GET("/payments/:id", getPaymentById)
	router.POST("/payments", postPayments)

	router.Run("localhost:8080")
}

func getPayments(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, payments)
}

func postPayments(c *gin.Context) {
	var newPayment Payment

	if err := c.BindJSON(&newPayment); err != nil {
		return
	}

	payments = append(payments, newPayment)
	c.IndentedJSON(http.StatusCreated, newPayment)
}

func getPaymentById(c *gin.Context) {
	id := c.Param("id")

	for _, payment := range payments {
		if payment.ID == id {
			c.IndentedJSON(http.StatusOK, payment)
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "payment not found"})
}
