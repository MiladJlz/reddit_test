package main

//import (
//	"github.com/joho/godotenv"
//	"log"
//	"os"
//)
//
//type Config struct {
//	Host     string
//	Port     string
//	Password string
//	User     string
//	DBName   string
//	SSLMode  string
//}
//
//func DBConfig() Config {
//	err := godotenv.Load(".env")
//	if err != nil {
//		log.Fatal("Error loading .env file")
//	}
//
//	return Config{
//		Host:     os.Getenv("DB_HOST"),
//		Port:     os.Getenv("DB_PORT"),
//		Password: os.Getenv("DB_PASS"),
//		User:     os.Getenv("DB_USER"),
//		DBName:   os.Getenv("DB_NAME"),
//		SSLMode:  os.Getenv("SSLMODE")}
//
//}
//func InitD
