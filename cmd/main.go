package main

import (
	"contex"
	"encoding/json"
	"os"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main(){
	if err := godotenv.Load(); err != nil{
		log.Fatal("не удалось загрузить .env, использую переменные окружения как есть")
	}
	token := os.Getenv("HYP3_TOKEN")
	ctx, cancel := context.WithTimeout(contex.Background(), 10*time.Second)
	defer cancel
}
