package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/golang/glog"
	"gopkg.in/yaml.v2"
)

var (
	config      Config
	dayDuration = time.Hour * 24
)

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
	_, err = exec.Command(`rm`, `-rf`, fmt.Sprintf(`~/.config/JetBrains/%s*/eval`, jetName)).Output()
	if err != nil {
		glog.Error(err)
		return
	}
	glog.Infof(`Успешно убил ключик для %s`, jetName)
	return
}

// ResetOther - Убиваю other
func ResetOther(jetName string) (err error) {
	_, err = exec.Command(`rm`, `-rf`, fmt.Sprintf(`~/.config/JetBrains/%s*/options/other.xml`, jetName)).Output()
	if err != nil {
		glog.Error(err)
		return
	}
	glog.Infof(`Успешно убил other для %s`, jetName)
	return
}

// ResetJetBrains - Убиваю jetbrains
func ResetJetBrains(jetName string) (err error) {
	_, err = exec.Command(`rm`, `-rf`, fmt.Sprintf(`~/.java/.userPrefs/jetbrains/%s`, strings.ToLower(jetName))).Output()
	if err != nil {
		glog.Error(err)
		return
	}
	glog.Infof(`Успешно убил jetbrains для %s`, jetName)
	return
}

// Resetter - Главный обнулятор
func Resetter() (err error) {
	if (time.Now().Sub(config.LastReset)).Hours() < float64(24*config.ResetCooldown) {
		glog.Info(`Пока еще рано обновлять лицензию`)
		return
	}

	for _, jetName := range config.JetList {
		err = ResetEval(jetName)
		if err != nil {
			return
		}
		err = ResetOther(jetName)
		if err != nil {
			return
		}
		err = ResetJetBrains(jetName)
		if err != nil {
			return
		}
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
