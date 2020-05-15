package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/allaboutapps/integresql-client-go"
)

func main() {
	hash := os.Getenv("INTEGRESQL_CLIENT_TEMPLATE_HASH")
	if len(hash) == 0 {
		log.Fatalln("No template hash provided, please set INTEGRESQL_CLIENT_TEMPLATE_HASH")
	}

	c, err := integresql.DefaultClientFromEnv()
	if err != nil {
		log.Fatalf("Failed to create IntegreSQL client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	test, err := c.GetTestDatabase(ctx, "meepmeep")
	if err != nil {
		log.Fatalf("Failed to retrieve test database: %v", err)
	}

	s, err := json.MarshalIndent(test, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal test database: %v", err)
	}

	fmt.Printf("%s\n", s)
}
