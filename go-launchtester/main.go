package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"runtime"
	"time"

	"drixevel.dev/launchtester/internal/config"
	"github.com/iancoleman/orderedmap"
	"github.com/shirou/gopsutil/v3/process"
	"github.com/urfave/cli/v2"
)

const (
	dirPrefix    = "D:/sourcemod/servers"     //Windows Only
	dirPrefixWSL = "/mnt/d/sourcemod/servers" //Linux Only
)

type Config struct {
	ServerDirectories []string              `json:"serverDirectories"`
	ServerConnects    orderedmap.OrderedMap `json:"serverConnects"`
}

func getPath(path string) string {
	//Change the prefix of the path depending on if it's Windows or Linux.
	var prefix string
	if runtime.GOOS == "linux" {
		prefix = dirPrefixWSL
	} else {
		prefix = dirPrefix
	}

	//Format a valid path.
	fPath := fmt.Sprintf("%s/%s/start.bat", prefix, path)
	return fPath
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func KillProcess(name string) error {
	processes, err := process.Processes()
	if err != nil {
		return err
	}
	for _, p := range processes {
		n, err := p.Name()
		if err != nil {
			return err
		}
		if n == name {
			if err = p.Kill(); err != nil {
				return err
			}
		}
	}
	return fmt.Errorf("Test environment not running.")
}

func startDevEnvironment(cfg *Config, game string, project string) {
	//Generate a path for the input.
	gamePath := getPath(game)

	//Start the test server for the game.
	fmt.Println("Starting the test server...")
	serverBinary := exec.Command(gamePath)
	serverBinary.Start()

	time.Sleep(15 * time.Second)

	//Start the game and connect to the test server.
	ip, ok := cfg.ServerConnects.Get(game)

	if !ok {
		log.Fatalf("error while finding ip for %s", game)
	}

	//Starting the game itself.
	fmt.Println("Starting game...")
	var url = fmt.Sprintf("steam://connect/%s", ip)
	gameBinary := exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	gameBinary.Start()

	//Starts VSCode just cuz.
	exec.Command("code", path.Join("D:", project)).Start()
}

func stopDevEnvironment() {
	KillProcess("srcds.exe")
	KillProcess("hl2.exe")
}

func main() {
	var cfg Config
	var game string
	var project string

	//Parse the config.json file.
	if err := config.ParseConfig(&cfg); err != nil {
		log.Fatal(err)
	}

	//Choose a game to setup the work environment for.
	//fmt.Printf("Choose a game you would like to launch the game and test server for:\nAvailable Options: %v\nInput: ", cfg.ServerDirectories)

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "start",
				Aliases: []string{"s"},
				Usage:   "Starts the game and test server",
				Action: func(c *cli.Context) error {
					startDevEnvironment(&cfg, game, project)
					return nil
				},
			},
			{
				Name:    "stop",
				Aliases: []string{"x"},
				Usage:   "Stops the game and test server",
				Action: func(c *cli.Context) error {
					stopDevEnvironment()
					return nil
				},
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "game",
				Aliases:     []string{"g"},
				Usage:       "Test server to run",
				Destination: &game,
			},
			&cli.StringFlag{
				Name:        "project",
				Aliases:     []string{"p"},
				Usage:       "Project to run in VSCode",
				Destination: &project,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
