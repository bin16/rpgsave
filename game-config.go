package main

import (
	"os"

	"github.com/BurntSushi/toml"
)

func (g *Game) WriteConfig() error {
	dst, err := os.Create(configPath)
	if err != nil {
		return err
	}

	enc := toml.NewEncoder(dst)
	return enc.Encode(g)
}
