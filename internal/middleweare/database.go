package middleweare

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/megajandrox/go-finance-api/pkg/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// InitializeDatabase initializes the database connection and runs migrations
func InitializeDatabase() *gorm.DB {
	// Cargar variables de entorno desde el archivo .env
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Obtener variables de entorno
	dbName := os.Getenv("POSTGRES_DB")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbHost := os.Getenv("POSTGRES_HOST")
	dbPort := os.Getenv("POSTGRES_PORT")

	// Formatear el DSN
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbHost, dbUser, dbPassword, dbName, dbPort)

	// Conectar a la base de datos
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Migrar el esquema
	db.AutoMigrate(&models.Asset{}, &models.Position{})

	fmt.Println("Database connected and migrated successfully")
	return db
}
