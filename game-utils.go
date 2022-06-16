package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func isFileExist(filename string) bool {
	if _, err := os.Stat(filename); err != nil {
		return false
	}

	return true
}

func findItemsJSON(savefile string) string {
	d := filepath.Dir(savefile)
	return filepath.Join(d, "..", "www/data", "Items.json")
}

func uLine(length int) string {
	return "+" + strings.Repeat("-", length-2) + "+"
}

func uPrint(kv ...string) string {
	var mk, mv int
	for i, u := range kv {
		if i%2 == 0 && len(u) > mk {
			mk = len(u)
		}

		if i%2 == 1 && len(u) > mv {
			mv = len(u)
		}
	}

	// | mk: mv |
	length := 2 + mk + 2 + mv + 2
	line := uLine(length)
	tpl := fmt.Sprintf("| %%%ds: %%%ds |", mk, mv)

	sl := []string{}

	sl = append(sl, line)
	for i := 0; i < len(kv)/2; i++ {
		s := strings.TrimSpace(kv[i*2+1])
		ss := s + strings.Repeat(" ", mv-len(s))
		sl = append(sl, fmt.Sprintf(tpl, kv[i*2], ss))
	}
	sl = append(sl, line)

	return strings.Join(sl, "\n")
}
