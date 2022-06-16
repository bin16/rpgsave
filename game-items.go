package main

import (
	"encoding/json"
	"os"
	"sort"
	"strconv"
)

type Item struct {
	ID          int64
	Name        string
	Description string

	count int64
}

func (g *Game) loadItems(filename string) ([]Item, error) {
	dl := []Item{}

	jsonFile, err := os.Open(filename)
	if err != nil {
		return dl, err
	}

	dec := json.NewDecoder(jsonFile)
	if err := dec.Decode(&dl); err != nil {
		return dl, err
	}

	items := []Item{}
	for _, d := range dl {
		if d.Name != "" {
			items = append(items, d)
		}
	}

	sort.Slice(g.items, func(i, j int) bool {
		return dl[i].ID < dl[j].ID
	})
	g.items = items

	for i, d := range g.items {
		k := strconv.Itoa(int(d.ID))
		if u, ok := g.itemMap[k]; ok {
			g.items[i].count = u.count
		}
	}

	return dl, nil
}

func (g *Game) Items(showAll bool) []Item {
	if showAll {
		return g.items
	}

	return GAME.save.Items()
}
