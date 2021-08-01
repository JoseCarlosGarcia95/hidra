// Represent essential entrypoint for hidra
package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/JoseCarlosGarcia95/hidra/agent"
	"github.com/JoseCarlosGarcia95/hidra/api"
	"github.com/JoseCarlosGarcia95/hidra/models"
	"github.com/JoseCarlosGarcia95/hidra/scenarios"
	_ "github.com/JoseCarlosGarcia95/hidra/scenarios/all"
	"github.com/joho/godotenv"
)

type flagConfig struct {
	testFile    string
	listenAddr  string
	configFile  string
	agentSecret string
	apiEndpoint string
	dataDir     string
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

func runAgentMode(cfg *flagConfig) {
	log.Println("Running hidra in agent mode")
	agent.StartAgent(cfg.apiEndpoint, cfg.agentSecret, cfg.dataDir)
}

func runApiMode(cfg *flagConfig) {
	log.Println("Running hidra in api mode")
	api.StartApi(cfg.listenAddr)
}

func main() {
	godotenv.Load()

	// Start default configuration
	cfg := flagConfig{}

	// Initialize flags
	var agentMode, apiMode, testMode bool
	flag.BoolVar(&apiMode, "api", false, "-api enable api mode in given hidra")
	flag.BoolVar(&agentMode, "agent", false, "-agent enable agent mode in given hidra")
	flag.BoolVar(&testMode, "test", false, "-test enable test mode in given hidra")
	flag.StringVar(&cfg.configFile, "config", "", "-config your configuration")
	flag.StringVar(&cfg.testFile, "file", "", "-file your-test-file-yaml")
	flag.StringVar(&cfg.listenAddr, "listen-addr", ":8080", "-listen-addr listen address")
	flag.StringVar(&cfg.agentSecret, "agent-secret", "", "-agent-secret for registering this agent")
	flag.StringVar(&cfg.apiEndpoint, "api-url", "", "-api-url where is api url?")
	flag.StringVar(&cfg.dataDir, "data-dir", "/tmp", "-data-dir where you want to store agent data?")

	flag.Parse()

	if agentMode {
		runAgentMode(&cfg)
	} else if apiMode {
		runApiMode(&cfg)
	} else if testMode {
		runTestMode(&cfg)
	}
}
