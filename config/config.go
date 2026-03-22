package config

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Config struct {
	AppPort         string
	MongoURI        string
	MongoDBName     string
	MongoCollection string
	UserCollection  string
	JwtSecret       string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading from environment")
	}

	return &Config{
		AppPort:         getEnv("APP_PORT", "8000"),
		MongoURI:        getEnv("MONGODB_URI", ""),
		MongoDBName:     getEnv("MONGO_DB_NAME", ""),
		MongoCollection: getEnv("MONGO_COLLECTION", "tasks"),
		UserCollection:  getEnv("USER_COLLECTION", "users"),
		JwtSecret:       getEnv("JWT_SECRET", ""),
	}
}

func ConnectMongo(cfg *Config) *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(options.Client().ApplyURI(cfg.MongoURI))

	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Ping to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("MongoDB ping failed: %v", err)
	}

	log.Println("Connected to MongoDB successfully")
	return client
}

func getEnv(key, fallback string) string {
	if ok := os.Getenv(key); ok != "" {
		return ok
	}
	return fallback
}
