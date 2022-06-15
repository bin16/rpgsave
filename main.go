package main

import (
	"fmt"
	"log"
	"os"

	"github.com/tidwall/sjson"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type Object map[string]any

type Save struct {
	actorIndex int
	json       string
}

func (s *Save) setData(k string, value any) {
	json, err := sjson.Set(s.json, k, value)
	if err != nil {
		log.Println(err)
	}

	s.json = json
}
