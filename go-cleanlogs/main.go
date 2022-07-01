/*
	Cleans log folders on a system manually through the folder names based on the Source game.
	Thanks to TheXeon for helping write this.
*/

package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"
)

const (
	dirPrefix    = "D:/sourcemod/servers"     //Windows Only
	dirPrefixWSL = "/mnt/d/sourcemod/servers" //Linux Only
	dirPostfix   = "addons/sourcemod/logs/"
)

//Games based on the base directory and the game directory itself.
var gameList = []string{"tf", "left4dead2", "csgo", "cstrike"}

func cleanLogs(path string) error {
	//Change the prefix of the path depending on if it's Windows or Linux.
	var prefix string
	if runtime.GOOS == "linux" {
		prefix = dirPrefixWSL
	} else {
		prefix = dirPrefix
	}

	//Format a valid path.
	fPath := fmt.Sprintf("%s/%s/%s/%s", prefix, path, path, dirPostfix)

	//Remove the directory entirely.
	if err := os.RemoveAll(fPath); err != nil {
		return err
	}

	//Remake the directory again.
	return os.Mkdir(fPath, 0750)
}

func main() {
	//Starting program, print it.
	fmt.Println("---\nCleaning up log folders...\n---")

	//Loop through the list of Source game servers to clean logs.
	for _, game := range gameList {
		if err := cleanLogs(game); err != nil {
			log.Fatal(err)
		}

		//Log as successful for each individual game.
		fmt.Println(game + " log file has been cleaned.")
	}

	//Success, print it.
	fmt.Println("---\nAll log folders have been cleaned successfully.\n---")

	//Keep the window open for 2 seconds.
	time.Sleep(time.Duration(2) * time.Second)
}
