package main

import (
	"context"
	"log"
	"os"
	"x/internal/processor"
)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		log.Println("failed to get working directory: ", err)
	}

	p := processor.New(
		processor.WithEnv(os.Environ()),
		processor.WithDirectory(wd+"/sample"),
	)
	p.Generate(context.Background())
}
