package main

import (
	"fmt"
	"strconv"
)

func (g *Game) Load() error {
	if err := g.loadSaveFile(g.FilePath); err != nil {
		return fmt.Errorf("error: failed to open %s: %s", g.FilePath, err.Error())
	}
	fmt.Printf("[ load ] %s\n", g.FilePath)

	itempath := findItemsJSON(g.FilePath)
	g.ItemPath = itempath
	if !isFileExist(itempath) {
		fmt.Printf("warning: failed to find the %s\n", itempath)
		return nil
	}

	items, err := g.loadItems(itempath)
	if err != nil {
		fmt.Printf("warning: failed to read %s: %s\n", itempath, err.Error())
		return nil
	}

	g.itemMap = make(map[string]Item)
	for _, item := range items {
		k := strconv.Itoa(int(item.ID))
		g.itemMap[k] = item
	}
	fmt.Printf("[ load ] %s\n", "Items.json")

	return nil
}

func (g *Game) Unload() error {
	fmt.Printf("unload")

	if err := g.WriteSave(); err != nil {
		return err
	}

	if err := g.WriteConfig(); err != nil {
		return err
	}

	return nil
}
