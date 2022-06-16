package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/bin16/rpgsave/lzstring"
)

func (g *Game) loadSaveFile(filename string) error {
	if !isFileExist(g.FilePath) {
		return fmt.Errorf("file not found")
	}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	raw := string(data)
	str := lzstring.Decode(raw)

	d := &Save{}
	if err := json.Unmarshal([]byte(str), d); err != nil {
		return err
	}
	d.json = str
	d.actorIndex = g.ActorIndex
	g.save = d

	return nil
}

func (g *Game) WriteSave() error {
	dat := lzstring.Encode(g.save.json)
	dst, err := os.Create(g.FilePath)
	if err != nil {
		return err
	}

	if _, err := dst.WriteString(dat); err != nil {
		return err
	}

	return nil
}
