package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/JoseCarlosGarcia95/hidra/models"
	"github.com/JoseCarlosGarcia95/hidra/scenarios"
	_ "github.com/JoseCarlosGarcia95/hidra/scenarios/all"
)

type flagConfig struct {
	testFile string
}

// This mode is used for fast checking yaml
func runTestMode(cfg *flagConfig) {
	if cfg.testFile == "" {
		log.Fatal("testFile expected to be not null")
	}

	if _, err := os.Stat(cfg.testFile); os.IsNotExist(err) {
		log.Fatal("testFile does not exists")
	}

	log.Println("Running hidra in test mode")
	data, err := ioutil.ReadFile(cfg.testFile)

	if err != nil {
		log.Fatal(err)
	}

	slist, err := models.ReadScenariosYAML(data)
	if err != nil {
		log.Fatal(err)
	}

	for _, s := range slist.Scenarios {
		m := scenarios.RunScenario(s)

		if m.Error != nil {
			log.Fatal(m.Error)
		}

		scenarios.PrettyPrintScenarioMetrics(m)
	}

}

func runRunnerMode(cfg *flagConfig) {
	log.Println("Running hidra in agent mode")
}

func runApiMode(cfg *flagConfig) {
	log.Println("Running hidra in api mode")
}

func main() {
	// Start default configuration
	cfg := flagConfig{}

	// Initialize flags
	var runnerMode, apiMode bool
	flag.BoolVar(&apiMode, "api", false, "--api enable api mode in given hidra")
	flag.BoolVar(&runnerMode, "runner", false, "--runner enable runner mode in given hidra")
	flag.StringVar(&cfg.testFile, "testFile", "", "--testFile your-test-file-yaml")
	flag.Parse()

	if runnerMode {
		runRunnerMode(&cfg)
	} else if apiMode {
		runApiMode(&cfg)
	} else {
		runTestMode(&cfg)
	}
}
