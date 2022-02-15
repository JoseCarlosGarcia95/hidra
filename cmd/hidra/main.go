// Represent essential entrypoint for hidra
package main

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/namsral/flag"

	"github.com/hidracloud/hidra/exporter"
	"github.com/hidracloud/hidra/models"
	"github.com/hidracloud/hidra/scenarios"
	_ "github.com/hidracloud/hidra/scenarios/all"
	"github.com/hidracloud/hidra/utils"
	"github.com/joho/godotenv"
)

type flagConfig struct {
	// Which configuration file to use
	testFile string

	// Which conf path to use
	confPath string

	// maxExecutor is the maximum number of executor to run in parallel
	maxExecutor int

	// port is the port to listen on
	port int
}

// This mode is used for fast checking yaml
func runTestMode(cfg *flagConfig, wg *sync.WaitGroup) {
	if cfg.testFile == "" {
		log.Fatal("testFile expected to be not null")
	}

	var configFiles = []string{cfg.testFile}

	// Check if test file ends with .yaml or .yml
	if !strings.Contains(cfg.testFile, ".yaml") && !strings.Contains(cfg.testFile, ".yml") {
		configFiles, _ = utils.AutoDiscoverYML(cfg.testFile)
	}

	for _, configFile := range configFiles {
		if _, err := os.Stat(configFile); os.IsNotExist(err) {
			log.Fatal("testFile does not exists")
		}

		log.Println("Running hidra in test mode")
		data, err := ioutil.ReadFile(configFile)

		if err != nil {
			log.Fatal(err)
		}

		slist, err := models.ReadScenariosYAML(data)
		if err != nil {
			log.Fatal(err)
		}

		if len(slist.Tags) > 0 {
			log.Println("Tags:")
			for key, val := range slist.Tags {
				log.Println(" ", key, "=", val)
			}
		}

		m := scenarios.RunScenario(slist.Scenario, slist.Name, slist.Description)

		if m.Error != nil {
			log.Fatal(m.Error)
		}

		scenarios.PrettyPrintScenarioResults(m, slist.Name, slist.Description)
	}
	wg.Done()

}

func runExporter(cfg *flagConfig, wg *sync.WaitGroup) {
	exporter.Run(wg, cfg.confPath, cfg.maxExecutor, cfg.port)
}

func runRCA(cfg *flagConfig, wg *sync.WaitGroup) {
	//rca.Run()
}

func main() {
	godotenv.Load()

	// Start default configuration
	cfg := flagConfig{}

	// Initialize flags
	var testMode, exporter, rca bool

	// Operating mode
	flag.BoolVar(&testMode, "test", false, "-test enable test mode in given hidra")
	flag.BoolVar(&exporter, "exporter", false, "-exporter enable exporter mode in given hidra")
	flag.BoolVar(&rca, "rca", false, "-rca enable rca mode in given hidra")

	// Test mode
	flag.StringVar(&cfg.testFile, "file", "", "-file your_test_file_yaml")
	flag.StringVar(&cfg.confPath, "conf", "", "-conf your_conf_path")

	// Exporter mode
	flag.IntVar(&cfg.maxExecutor, "maxExecutor", 1, "-maxExecutor your_max_executor")
	flag.IntVar(&cfg.port, "port", 19090, "-port your_port")
	flag.Parse()

	var wg sync.WaitGroup

	if testMode {
		wg.Add(1)
		go runTestMode(&cfg, &wg)
	}

	if exporter {
		wg.Add(1)
		go runExporter(&cfg, &wg)
	}

	if rca {
		wg.Add(1)
		go runRCA(&cfg, &wg)
	}

	wg.Wait()
}
