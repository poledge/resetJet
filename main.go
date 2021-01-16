package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/golang/glog"
	"gopkg.in/yaml.v2"
)

var (
	config      Config
	dayDuration = time.Hour * 24
	homeDir     string
)

func init() {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	homeDir = usr.HomeDir
}

// Config - Конфиг
type Config struct {
	JetList       []string  `yaml:"jet_list"`
	ResetCooldown int       `yaml:"reset_cooldown"`
	LastReset     time.Time `yaml:"last_reset"`
}

// ParseConfig - Разбираю конфиг
func ParseConfig() (err error) {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	data, err := ioutil.ReadFile(path.Join(dir, `config.yaml`))
	if err != nil {
		glog.Error(err)
		return
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		glog.Error(err)
		return
	}
	return
}

// SaveConfig - Сохраняю конфиг
func SaveConfig() (err error) {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// Обновляю дату ресета
	config.LastReset = time.Now()

	var data []byte
	data, err = yaml.Marshal(config)
	if err != nil {
		glog.Error(err)
		return
	}

	ioutil.WriteFile(path.Join(dir, `config.yaml`), data, 0)
	if err != nil {
		glog.Error(err)
		return
	}
	return
}

// ResetEval - Убиваю ключик активации
func ResetEval(jetName string) (err error) {
	files, err := filepath.Glob(fmt.Sprintf(`%s/.config/JetBrains/%s*`, homeDir, jetName))
	if err != nil {
		return
	}
	for _, file := range files {
		err = os.RemoveAll(fmt.Sprintf(`%s/eval`, file))
		if err != nil {
			glog.Error(err)
		}
	}
	if err == nil {
		glog.Infof(`Успешно убил ключик для %s`, jetName)
	}
	return
}

// ResetOther - Убиваю other
func ResetOther(jetName string) (err error) {
	files, err := filepath.Glob(fmt.Sprintf(`%s/.config/JetBrains/%s*`, homeDir, jetName))
	if err != nil {
		return
	}
	for _, file := range files {
		err = os.Remove(fmt.Sprintf(`%s/options/other.xml`, file))
		if err != nil {
			glog.Error(err)
		}
	}
	if err == nil {
		glog.Infof(`Успешно убил other для %s`, jetName)
	}
	return
}

// ResetJetBrains - Убиваю jetbrains
func ResetJetBrains(jetName string) (err error) {
	err = os.RemoveAll(fmt.Sprintf(`%s/.java/.userPrefs/jetbrains/%s`, homeDir, strings.ToLower(jetName)))
	if err != nil {
		glog.Error(err)
	}
	if err == nil {
		glog.Infof(`Успешно убил jetbrains для %s`, jetName)
	}
	return
}

// Resetter - Главный обнулятор
func Resetter() (err error) {
	if (time.Now().Sub(config.LastReset)).Hours() < float64(24*config.ResetCooldown) {
		glog.Info(`Пока еще рано обновлять лицензию`)
		return
	}

	for _, jetName := range config.JetList {
		ResetEval(jetName)
		ResetOther(jetName)
		ResetJetBrains(jetName)
	}

	err = SaveConfig()
	if err != nil {
		return
	}

	return
}

func main() {
	var err error
	err = ParseConfig()
	if err != nil {
		return
	}

	err = Resetter()
	if err != nil {
		return
	}

	return
}
