package config

import (
	"io/ioutil"
	"log"
	"time"

	"gopkg.in/yaml.v2"
)

// ProcChanSize: 1000
// Sched:
//   TimeQuatum: 50
// CPU:
//   ClockSpeed: 10
// Memory:
//   PageSize: 32
//   TotalRam: 4096
//
// Note: struct fields must be public in order for unmarshal to
// correctly populate the data.
type Config struct {
	ProcChanSize      int     `yaml:"ProcChanSize"`
	MinimumFreeFrames int     `yaml:"MinimumFreeFrames"`
	Sched             *Sched  `yaml:"Sched"`
	CPU               *CPU    `yaml:"CPU"`
	Memory            *Memory `yaml:"Memory"`
}

// Sched : Scheduler configurations
type Sched struct {
	TimeQuantum int `yaml:"TimeQuantum"`
}

// CPU : CPU configuration
type CPU struct {
	ClockSpeed1 time.Duration `yaml:"ClockSpeed1"`
	ClockSpeed2 time.Duration `yaml:"ClockSpeed2"`
}

// Memory : Memory configuration
type Memory struct {
	PageSize int `yaml:"PageSize"`
	TotalRam int `yaml:"TotalRam"`
}

// ReadConfig : read config file and serialize
func ReadConfig(file string) Config {

	conf := Config{}

	yamlFile, err := ioutil.ReadFile(file)
	if err != nil {
		log.Printf("[ERROR] yamlFile.Get err   #%v ", err)
	}

	err = yaml.Unmarshal([]byte(yamlFile), &conf)
	if err != nil {
		log.Fatalf("[ERROR] %v", err)
	}

	if conf.ProcChanSize <= 0 {
		log.Fatal("[ERROR] Process channel size must be above zero")
	}

	if conf.MinimumFreeFrames <= 0 {
		log.Fatal("[ERROR] Minimum Free Frames must be above zero")
	}

	if conf.Sched.TimeQuantum <= 0 {
		log.Fatal("[ERROR] Time Quantum must be above zero")
	}

	if conf.CPU.ClockSpeed1 <= 0 {
		log.Fatal("[ERROR] ClockSpeed must be above zero")
	}

	if conf.Memory.PageSize <= 0 {
		log.Fatal("[ERROR] Page Size must be above zero")
	}

	if conf.Memory.TotalRam <= 0 || conf.Memory.TotalRam < conf.Memory.PageSize {
		log.Fatal("[ERROR] Total RAM must be above zero")
	}

	return conf

}
