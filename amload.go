package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/urfave/cli"
)

// Panic if there is an error
func check(err error) {
	if err != nil {
		panic(err)
	}
}

type Labels struct {
	Alertid   string `json:"alertid"`
	Service   string `json:"service"`
	Plan      string `json:"plan"`
	Region    string `json:"region"`
	Env       string `json:"env"`
	Host      string `json:"host`
	IP        string `json:"ip"`
	AlertName string `json:"alertname"`
}

type Annotation struct {
	Checkname   string `json:"checkname"`
	Severity    string `json:"severity"`
	Runbook     string `json:"runbook"`
	Description string `json:"description"`
	Servicekey  string `json:"serviceke"`
	Target      string `json:"target"`
	AlertValue  string `json:"alertValu"`
	AlertWarn   string `json:"alertWarn"`
	AlertError  string `json:"alertError"`
}

type Alert struct {
	Labels     `json:"labels"`
	Annotation `json:"annotation"`
	StartsAt   string `json:"startsAt"`
}

type Alerts []Alert

func main() {
	var (
		alerts int
	)

	// The Go random number generator source is deterministic, so we need to seed
	// it to avoid getting the same output each time
	rand.Seed(time.Now().UTC().UnixNano())

	// Configure our command line app
	app := cli.NewApp()
	app.Name = "Alertmanager alert  Data Generator"
	app.Usage = "generate a stream of test data for vegeta. Type 'amload help' for details"

	// Add -users flag, which defaults to 5
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:        "alerts",
			Value:       5,
			Usage:       "Number of alerts to simulate",
			Destination: &alerts,
		},
	}

	// Our app's main action
	app.Action = func(c *cli.Context) error {

		// Combine verb and URL to a target for Vegeta
		verb := c.Args().Get(0)
		url := c.Args().Get(1)
		target := fmt.Sprintf("%s %s", verb, url)

		if len(target) > 1 {

			for i := 1; i < alerts; i++ {

				// Generate request data

				labels := &Labels{
					Alertid:   "DSH00016",
					Service:   "dash",
					Plan:      "default",
					Region:    "dallas",
					Env:       "ys1",
					Host:      "node" + strconv.Itoa(i),
					IP:        "172.30.36." + strconv.Itoa(i),
					AlertName: "Memory used is high",
				}

				annotations := &Annotation{
					Checkname:   "Memory used is high",
					Severity:    "CRITICAL",
					Description: "Memory used is high. ",
					Servicekey:  "xxxxxxxxxxxx",
					Target:      "ys1.dallas.dash.default.node0_172_30_36_0.memory.memory-used",
					Runbook:     "https://ibm.biz/BdHsQb",
					AlertValue:  "1.0",
					AlertWarn:   "0,5",
					AlertError:  "1.0",
				}

				alert := &Alert{
					Labels:     *labels,
					Annotation: *annotations,
					StartsAt:   time.Now().Format(time.RFC3339),
				}

				alerts := Alerts{
					*alert,
				}

				// Convert the map to JSON
				body, err := json.Marshal(alerts)
				check(err)

				// Create a tmp directory to write our JSON files
				err = os.MkdirAll("tmp", 0755)
				check(err)

				// Use the user's name as the filename
				filename := fmt.Sprintf("tmp/%s.json", "alert"+strconv.Itoa(i))

				// Write the JSON to the file
				err = ioutil.WriteFile(filename, body, 0644)
				check(err)

				// Get the absolute path to the file
				filePath, err := filepath.Abs(filename)
				check(err)

				// Print the attack target
				fmt.Println(target)

				// Print '@' followed by the absolute path to our JSON file, followed by
				// two newlines, which is the delimiter Vegeta uses
				fmt.Printf("@%s\n\n", filePath)
			}
		} else {
			// Return an error if we're missing the required command line arguments
			return cli.NewExitError("You must specify the target in format 'VERB url'", 1)
		}
		return nil
	}

	app.Run(os.Args)
}
