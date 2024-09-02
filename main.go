package main

import (
	"database/sql"
	"time"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db *sql.DB

type Product struct {
	Id          int64
	Name        string
	Description string
	Image       string
	Category    string
	Price       float32
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	connectionStr := os.Getenv("DATABASE_URL")

	db, err = sql.Open("postgres", connectionStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	fmt.Println("Successfully connected!")

	// Configure CORS
	config := cors.DefaultConfig()
    config.AllowAllOrigins = true
    config.AllowMethods = []string{"POST", "GET", "PUT", "OPTIONS"}
    config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "Accept", "User-Agent", "Cache-Control", "Pragma"}
    config.ExposeHeaders = []string{"Content-Length"}
    config.AllowCredentials = true
    config.MaxAge = 12 * time.Hour


	// Setup routes
	router := gin.Default()
	router.Use(cors.New(config))
	router.GET("/products", GetProducts)
	router.GET("/products/:productId", GetSingleProduct)
	router.POST("/products", CreateProduct)
	router.PUT("/products/:productId", UpdateProduct)
	router.DELETE("/products/:productId", DeleteProduct)


	// Run the router
	router.Run()

}

func GetProducts(c *gin.Context) {
	query := `SELECT * FROM products`
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal("(GetProducts) db.Query", err)
	}
	defer rows.Close()

	products := []Product{}
	for rows.Next() {
		var product Product
		err := rows.Scan(&product.Id, &product.Name, &product.Description, &product.Image, &product.Category, &product.Price)
		if err != nil {
			log.Fatal("(GetProducts) rows.Scan", err)
		}
		products = append(products, product)
	}

	c.JSON(http.StatusOK, products)
}

func GetSingleProduct(c *gin.Context) {
	productId := c.Param("productId")
	productId = strings.ReplaceAll(productId, "/", "")
	productIdInt, err := strconv.Atoi(productId)
	if err != nil {
		log.Fatal("(GetSingleProduct) strconv.Atoi", err)
	}

	var product Product
	query := `SELECT * FROM products WHERE id = $1`
	err = db.QueryRow(query, productIdInt).Scan(&product.Id, &product.Name, &product.Description, &product.Image, &product.Category, &product.Price)
	if err != nil {
		log.Fatal("(GetSingleProduct) db.Exec", err)
	}

	c.JSON(http.StatusOK, product)
}

func CreateProduct(c *gin.Context) {
	var newBike Product
	err := c.BindJSON(&newBike)
	if err != nil {
		log.Fatal("(CreateProduct) c.BindJSON", err)
	}

	query := `INSERT INTO products (name, description, image, category, price) VALUES ($1, $2, $3, $4, $5)`
	res, err := db.Exec(query, newBike.Name, newBike.Description, newBike.Image, newBike.Category, newBike.Price)
	if err != nil {
		log.Fatal("(CreateProduct) db.Exec", err)
	}
	newBike.Id, err = res.LastInsertId()
	if err != nil {
		log.Fatal("(CreateProduct) res.LastInsertId", err)
	}

	c.JSON(http.StatusOK, newBike)
}

func UpdateProduct(c *gin.Context) {
	var updates Product
	err := c.BindJSON(&updates)
	if err != nil {
		log.Fatal("(UpdateProduct) c.BindJSON", err)
	}

	productId := c.Param("productId")
	productId = strings.ReplaceAll(productId, "/", "")
	productIdInt, err := strconv.Atoi(productId)
	if err != nil {
		log.Fatal("(UpdateProduct) strconv.Atoi", err)
	}

	query := `UPDATE products SET name = $1, description = $2, image = $3, category = $4, price = $5 WHERE id = $6`
	_, err = db.Exec(query, updates.Name, updates.Description, updates.Image, updates.Category, updates.Price, productIdInt)
	if err != nil {
		log.Fatal("(UpdateProduct) db.Exec", err)
	}

	c.Status(http.StatusOK)
}

func DeleteProduct(c *gin.Context) {
	productId := c.Param("productId")

	productId = strings.ReplaceAll(productId, "/", "")
	productIdInt, err := strconv.Atoi(productId)
	if err != nil {
		log.Fatal("(DeleteProduct) strconv.Atoi", err)
	}
	query := `DELETE FROM products WHERE id = $1`
	_, err = db.Exec(query, productIdInt)
	if err != nil {
		log.Fatal("(DeleteProduct) db.Exec", err)
	}

	c.Status(http.StatusOK)
}
