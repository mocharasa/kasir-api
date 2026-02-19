package main

import (
	"encoding/json"
	"fmt"
	"kasir-api/database"
	"kasir-api/handlers"
	"kasir-api/middlewares"
	"kasir-api/repositories"
	"kasir-api/services"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Port    string `mapstructure:"PORT"`
	DBConn  string `mapstructure:"DB_CONN"`
	API_KEY string `mapstructure:"API_KEY"`
}

func main() {
	//1. setup viper
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}
	config := Config{
		Port:    viper.GetString("PORT"),
		DBConn:  viper.GetString("DB_CONN"),
		API_KEY: viper.GetString("API_KEY"),
	}

	//2. setup database
	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	//3. dependency injection [harus diatas HandleFunc]
	apiKeyMiddleware := middlewares.APIKey(config.API_KEY)
	productRepo := repositories.NewProductRepository(db)
	categoryRepo := repositories.NewCategoryRepository(db)
	transactionRepo := repositories.NewTransactionRepository(db)
	reportRepo := repositories.NewReportRepository(db)

	productService := services.NewProductService(productRepo)
	categoryService := services.NewCategoryService(categoryRepo)
	transactionService := services.NewTransactionService(transactionRepo)
	reportService := services.NewReportService(reportRepo)

	productHandler := handlers.NewProductHandler(productService)
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	transactionHandler := handlers.NewTransactionHandler(transactionService)
	reportHandler := handlers.NewReportHandler(reportService)

	//4. setup routes
	//health check
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		healthCheck(w, r)
	})
	http.HandleFunc("/api/product", middlewares.CORS(middlewares.Logger(productHandler.HandleProducts)))
	http.HandleFunc("/api/product/", middlewares.CORS(middlewares.Logger(apiKeyMiddleware(productHandler.HandleProductByID))))
	http.HandleFunc("/api/categories", middlewares.CORS(middlewares.Logger(categoryHandler.HandleCategories)))
	http.HandleFunc("/api/categories/", middlewares.CORS(middlewares.Logger(apiKeyMiddleware(categoryHandler.HandleCategoryByID))))
	http.HandleFunc("/api/checkout", middlewares.CORS(middlewares.Logger(apiKeyMiddleware(transactionHandler.HandleCheckout)))) // POST
	http.HandleFunc("/api/report/hari-ini", reportHandler.HandleReport)                                                         // GET
	http.HandleFunc("/api/report", reportHandler.HandleReport)                                                                  // GET

	// 5. Definisikan Handler untuk Root (Opsional, agar muncul saat web dibuka)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "OK",
			"routes": map[string]string{
				"product":  "/api/product",
				"category": "/api/categories",
				"checkout": "/api/checkout",
				"report":   "/api/report",
			},
		})
	})

	//6. run server
	addr := "0.0.0.0:" + config.Port
	fmt.Println("Server running di", addr)
	err = http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Println("gagal running server", err)
	}

}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "OK",
		"message": "API running",
	})
}
