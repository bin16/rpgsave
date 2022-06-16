package main

import (
	"fmt"
	"sort"

	"github.com/tidwall/gjson"
)

func (s *Save) Actor(index int) *Save {
	s.actorIndex = index
	return s
}

func (s *Save) Gold() int64 {
	return int64(gjson.Get(s.json, "party._gold").Float())
}

func (s *Save) SetGold(num int64) {
	s.setData("party._gold", num)
}

func (s *Save) HP() float64 {
	path := fmt.Sprintf("actors._data.%0d._hp", s.actorIndex)
	return gjson.Get(s.json, path).Float()
}

func (s *Save) MP() float64 {
	return gjson.Get(s.json, fmt.Sprintf("actors._data.%0d._mp", s.actorIndex)).Float()
}

func (s *Save) Exp() int64 {
	return gjson.Get(s.json, fmt.Sprintf("actors._data.%0d._exp.1", s.actorIndex)).Int()
}

func (s *Save) SetExp(exp int64) {
	s.setData(fmt.Sprintf("actors._data.%0d._exp.1", s.actorIndex), exp)
}

func (s *Save) AddExp(exp int64) {
	s.setData(fmt.Sprintf("actors._data.%0d._exp.1", s.actorIndex), s.Exp()+exp)
}

const (
	MaxHP = 0
	MaxMP = 1
	ATK   = 2
	DEF   = 3
	MAT   = 4
	MDF   = 5
	AGI   = 6
	LUK   = 7
)

func (s *Save) Extra(name int) int64 {
	path := fmt.Sprintf("actors._data.%0d._paramPlus.%d", s.actorIndex, name)
	return gjson.Get(s.json, path).Int()
}

func (s *Save) SetExtra(name int, num int64) {
	path := fmt.Sprintf("actors._data.%0d._paramPlus.%d", s.actorIndex, name)
	s.setData(path, num)
}

func (s *Save) Name() string {
	path := fmt.Sprintf("actors._data.%0d._name", s.actorIndex)
	return gjson.Get(s.json, path).String()
}

func (s *Save) SetName(name string) {
	path := fmt.Sprintf("actors._data.%0d._name", s.actorIndex)
	s.setData(path, name)
}

func (d *Save) Print() {
	name := d.Name()
	fmt.Printf("NAME: %s\n", name)

	gold := d.Gold()
	fmt.Printf("Gold: %d\n", gold)

	exp := d.Exp()
	fmt.Printf("Exp: %d\n", exp)

	hp := d.HP()
	mp := d.MP()
	fmt.Printf("HP: %.0f, MP: %.0f\n", hp, mp)

	mhp := d.Extra(MaxHP)
	mmp := d.Extra(MaxMP)
	fmt.Printf("MHP: %.0f, MMP: %.0f\n", mhp, mmp)

	atk := d.Extra(ATK)
	def := d.Extra(DEF)
	fmt.Printf("ATK: %.0f, DEF: %.0f\n", atk, def)

	mat := d.Extra(MAT)
	mdf := d.Extra(MDF)
	fmt.Printf("MAT: %.0f, MDF: %.0f\n", mat, mdf)

	agi := d.Extra(AGI)
	luk := d.Extra(LUK)
	fmt.Printf("AGI: %.0f, LUK: %.0f\n", agi, luk)
}

func (s *Save) Items() []Item {
	items := []Item{}

	itemMap := gjson.Get(s.json, "party._items").Map()
	for k, v := range itemMap {
		info := GAME.itemMap[k]
		u := Item{
			ID:          info.ID,
			Name:        info.Name,
			Description: info.Description,
			count:       v.Int(),
		}

		items = append(items, u)
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].ID < items[j].ID
	})

	return items
}

func (s *Save) SetItem(id, num int64) {
	path := fmt.Sprintf("party._items.%0d", id)
	s.setData(path, num)
}

func (s *Save) Item(id int64) *Item {
	info, ok := GAME.itemMap[fmt.Sprintf("%d", id)]
	if !ok {
		return nil
	}

	path := fmt.Sprintf("party._items.%0d", id)
	num := gjson.Get(s.json, path).Int()
	return &Item{
		ID:          id,
		Name:        info.Name,
		Description: info.Description,
		count:       num,
	}
}
