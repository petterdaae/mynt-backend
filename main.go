package main

import (
	"mynt/cmd"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
)

func main() {
	cmd.Execute()
}
