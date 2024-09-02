# Connect An Application to Neon Example (Go)

Learn how to use Go with Neon by exploring a demo Go API built with [Gin](https://github.com/gin-gonic/gin) plus [PQ driver](https://github.com/lib/pq) and how the API is used to create and manage E-bike products. E-Bikes are hands down one of the best ways to experience an area up close while covering a lot of ground easily.

**Overview**

To connect a Neon PostgreSQL database with a sample E-Bike Go API application built using the Gin routing framework and PQ driver, complete the following steps:

1. Set up your environment.
2. Set up your Postgres database using the Neon console, Neon CLI or Psql.
3. Build and run your Go API.
4. With a Neon account, database, and a Go API application, you can perform basic CRUD operations by making your API calls using any of the REST API clients including [cURL](https://curl.se/download.html) and [POSTMAN](https://www.postman.com/). Also, an application such as a [Vue](https://vuejs.org/) or [React](https://react.dev/) app can connect to the GetProducts endpoint to retrieve and render all the products in a database table.

**Prerequisites**

1. [Go](https://go.dev/doc/install) 1.23 or higher.
2. [A Neon account](https://console.neon.tech/signup)
3. [VS Code](https://code.visualstudio.com/download) (optional)
4. [POSTMAN](https://www.postman.com/) (optional)
5. [cURL](https://curl.se/download.html) (optional)
6. [Vue](https://vuejs.org/guide/quick-start) (optional)
7. [Vuetify](https://vuetifyjs.com/en/getting-started/installation/#installation) (optional)
8. [Vite](https://vitejs.dev/guide/#scaffolding-your-first-vite-project) (optional)

**Create a Neon Postges database**

1. If you do not have a Neon account, click [here](https://console.neon.tech/signup) to sign up for an account.
2. [Log in](https://console.neon.tech/login) to your Neon account.
3. On the Console page, click Create project.
4. On the Create project page, the highest postgres version is selected by default. Name the project **goebike** and name the database **productdb**. Select the desired **Region**. 
5. Click **Create project** to create a Neon project with a database.
![Image description](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/5c22f9zlc2n7g1v3yw1r.png)

6. Your Neon project and database will be created in just a few.
7. After your Neon database is created, click the **SQL Editor** menu and run the following three commands to drop an existing sample table, create a new sample table and insert some data:

```
DROP TABLE IF EXISTS products;
CREATE TABLE products (
  id            SERIAL PRIMARY KEY,
  name          VARCHAR(128) NOT NULL,
  description   VARCHAR(255) NOT NULL,
  image         VARCHAR(128) NOT NULL,
  category      VARCHAR(128) NOT NULL,
  price         DECIMAL(5,2) NOT NULL
);

INSERT INTO products
  (name, description, image, category, price)
VALUES
  ('ELECTRA X2', 'Lorem ipsum dolor sit amet, consectetur adipiscing elit', 'https://s3-us-west-2.amazonaws.com/dev-or-devrl-s3-bucket/sample-apps/ebikes/electrax2.jpg', 'Mountain', 56.99),
  ('ELECTRA X3', 'Lorem ipsum dolor sit amet, consectetur adipiscing elit', 'https://s3-us-west-2.amazonaws.com/dev-or-devrl-s3-bucket/sample-apps/ebikes/electrax3.jpg', 'Mountain', 63.99),
  ('ELECTRA X1', 'Lorem ipsum dolor sit amet, consectetur adipiscing elit', 'https://s3-us-west-2.amazonaws.com/dev-or-devrl-s3-bucket/sample-apps/ebikes/electrax1.jpg', 'Mountain', 34.98);
```

8.On the menu bar, click **Quickstart**, select your preferred connection method to get the corresponding connection string. This document uses Postgres client as an example. Copy the connection string for Postgres, you will need it in the next section.

**Run the sample Go API to connect to Neon**

This section illustrates how to run the sample Go API application code and connect to Neon.

**Step 1: Clone the sample app repository**

Run the following commands in your terminal window to clone the sample code repository:

```
git clone https://github.com/mikoaro/go-gin-getting-started.git
cd go-gin-getting-started.git
```
Open the project in VS Code or any other editor of your choice.

**Step 2: Run the code and check the result**

1. Run the following command to copy **.env.example** and rename it to **.env**:

```
cp .env.example .env
```
2. Copy and paste the corresponding connection string into the .env file. The example result is as follows:

```
DATABASE_URL='{}'
```
Be sure to replace the placeholders {} with the connection parameters for Postgres obtained from the Quickstart menu.

3. Save the **.env** file. 

4. Open a VS Code terminal and run the project using the following commands:

```
go mod tidy
go run .
```
The following output will be printed in the terminal.

![Image description](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/tv5x5r55t4gv8iwu8abj.png)

**Sample code snippets and testing with POSTMAN and Vue**

You can refer to the following sample code snippets to complete your own application development.

For complete sample code and how to run it, check out the [mikoaro/go-gin-getting-started](https://github.com/mikoaro/go-gin-getting-started) repository.

**Connect to Neon and Setup Routing**

```
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
```

**Query data**

```
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
```

Here’s an example of sending a GET request to the specified route to get all the e-bikes using Postman:
![Image description](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/6yzs9y6k6vst6ojo5jn1.png)

Here’s an example of sending a GET request to the specified route to retrieve an e-bike by its ID:

![Image description](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/uug5wejj8zc3vuzoy8ci.png)

**GoeBike Vue App**


![Image description](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/au0toty7dofuphakzsdy.png)

**Coming Soon**

If you are searching for fun things to do, then look no further than GoeBike rentals and tours. It is a fun and easy outdoor activity that the entire family can enjoy.

Electric bikes are one of the very best ways to explore an area. You can easily cover lots of ground, while still getting an up close and personal look at all the scenery and attractions.