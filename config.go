package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/bin16/rpgsave/lzstring"
)

var DATA *App

const (
	configPath = "rpgsave.toml"
)

func init() {
	// if _, err := os.Stat(configPath); os.IsNotExist(err) {
	// 	DATA = initApp(configPath)
	// 	return
	// }

	// d, err := loadApp(configPath)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// DATA = d
	// DATA.Load()
	// DATA.save.Actor(DATA.Actor.Index)
}

type AppConfig struct {
	Backup bool
}

type AppBackup struct {
	FilePath  string
	UpdatedAt time.Time
}

type AppSave struct {
	FilePath  string
	UpdatedAt time.Time
}

type AppActor struct {
	Name  string
	Index int
}

type App struct {
	Backup AppBackup
	Save   AppSave
	Actor  AppActor
	Config AppConfig

	save *Save
}

func (a *App) Load() {
	log.Printf("load")
	if err := a.ReadSave(); err != nil {
		log.Fatalln(err)
	}

	if err := a.DebugJSON(); err != nil {
		log.Fatalln(err)
	}
}

func (a *App) Unload() {
	log.Printf("unload")
	if err := a.WriteJSON(); err != nil {
		log.Fatalln(err)
	}

	if err := a.WriteSave(); err != nil {
		log.Fatalln(err)
	}

	if err := a.WriteConfig(); err != nil {
		log.Fatalln(err)
	}
}

func (a *App) DebugJSON() error {
	data, err := ioutil.ReadFile(a.Save.FilePath)
	if err != nil {
		return err
	}

	raw := string(data)

	str := lzstring.Decode(raw)

	p := strings.Replace(a.Save.FilePath, path.Ext(a.Save.FilePath), ".debug.json", -1)
	dst, err := os.Create(p)
	if err != nil {
		return err
	}

	if _, err := dst.WriteString(str); err != nil {
		return err
	}

	return nil
}

func (a *App) ReadSave() error {
	data, err := ioutil.ReadFile(a.Save.FilePath)
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
	a.save = d

	return nil
}

func (a *App) WriteSave() error {
	dat := lzstring.Encode(a.save.json)
	dst, err := os.Create(a.Save.FilePath)
	if err != nil {
		return err
	}

	if _, err := dst.WriteString(dat); err != nil {
		return err
	}

	return nil
}

func (a *App) WriteJSON() error {
	p := strings.Replace(a.Save.FilePath, path.Ext(a.Save.FilePath), ".json", -1)
	dst, err := os.Create(p)
	if err != nil {
		return err
	}

	if _, err := dst.WriteString(a.save.json); err != nil {
		return err
	}

	return nil
}

func (a *App) WriteConfig() error {
	dst, err := os.Create(configPath)
	if err != nil {
		return err
	}

	enc := toml.NewEncoder(dst)
	return enc.Encode(a)
}

func loadApp(filename string) (*App, error) {
	d := &App{}
	if _, err := toml.DecodeFile(filename, d); err != nil {
		return d, err
	}

	return d, nil
}

func initApp(filename string) *App {
	return &App{
		Save: AppSave{
			FilePath: filename,
		},
		Config: AppConfig{
			Backup: false,
		},
		Actor: AppActor{
			Name:  "",
			Index: 1,
		},
	}
}
