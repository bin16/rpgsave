package main

import (
	"log"

	"github.com/BurntSushi/toml"
)

type Game struct {
	itemMap map[string]Item
	items   []Item

	save       *Save
	FilePath   string
	ItemPath   string
	ActorIndex int
}

var GAME *Game

func init() {
	GAME = &Game{ActorIndex: 1}
	if !isFileExist(configPath) {
		return
	}

	if _, err := toml.DecodeFile(configPath, GAME); err != nil {
		log.Fatalln("failed to load config", err.Error())
		return
	}
}
