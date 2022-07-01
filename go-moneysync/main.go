package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/iancoleman/orderedmap"
)

type Config struct {
	Totals []float64             `json:"totals"`
	Costs  orderedmap.OrderedMap `json:"costs"`
}

func getConfig() *Config {
	file, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatal(err)
	}
	var config Config
	err = json.Unmarshal(file, &config)
	if err != nil {
		log.Fatal(err)
	}
	return &config
}

func main() {
	//Calculate the total dollars made this year so far.
	var total float64
	referenceValues := []float64{150.00, 100.00, 50.00, 25.00, 20.00, 15.00, 10.00, 5.00, 2.00}

	config := getConfig()

	//Add up the total.
	for _, i := range config.Totals {
		total += i
	}

	//Show total values.
	fmt.Println("-- Total Values --")
	fmt.Printf("Total So Far: %.2f\n", total)
	for _, k := range config.Costs.Keys() {
		val, _ := config.Costs.Get(k)
		floatVal := val.(float64)
		fmt.Printf("Total %s: %.2f\n", k, total*floatVal)
	}

	//Reference Values
	fmt.Println("\n-- Calculations --")
	for _, i := range referenceValues {
		fmt.Printf("%.2f = ", i)
		for keyIndex, k := range config.Costs.Keys() {
			val, _ := config.Costs.Get(k)
			floatVal := val.(float64)

			line := fmt.Sprintf("%s(%.2f) | ", k, i*floatVal)

			if len(config.Costs.Keys())-1 == keyIndex {
				line = strings.TrimSuffix(line, " | ")
			}

			fmt.Printf(line)
		}
		fmt.Println()
	}

	//Wait for a key to be pressed to exit the program.
	fmt.Println("\nPress any key to exit.")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
