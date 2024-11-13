package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const monitoring = 2
const delay = 5

func main() {
	showIntroduction()

	for {
		showMenu()

		selection := command()

		switch selection {
		case 1:
			startCheck()
		case 2:
			startMonitoring()
		case 3:
			printLogs()
		case 4:
			fmt.Println("Leaving the program...")
			os.Exit(3)
		default:
			fmt.Println("Type a valid command")
			os.Exit(-1)
		}
	}
}

func showIntroduction() {
	fmt.Println("This program checks the status of the sites and keeps a historical record of the periods of activity and inactivity of each site.")
}

func showMenu() {
	fmt.Println("1 - Type to start check")
	fmt.Println("2 - Start monitoring")
	fmt.Println("3 - View logs")
	fmt.Println("4 - Exit")
}

func command() int {
	var command int
	fmt.Scan(&command)
	fmt.Println("you chose the", command, "command")

	return command
}

// check website status
func checkStatus(site string) {
	resp, err := http.Get(site)

	if err != nil {
		fmt.Println("error", err)
	}

	if resp.StatusCode == 200 {
		fmt.Println("Site:", site, "has been loaded successfully")
		log(site, true)
	} else {
		fmt.Println("Site:", site, "has problems. Response Status:", resp.StatusCode)
		log(site, false)
	}
}

// starts checking with a url that the CLI User enters
func startCheck() string {
	fmt.Println("Please, type url: ")
	var site string
	fmt.Scan(&site)

	checkStatus(site)

	return site
}

func startMonitoring() {
	fmt.Println("Start monitoring...")

	sites := readWebsite()

	for i := 0; i < monitoring; i++ {
		for i, site := range sites {
			fmt.Println("site test:", i, ":", site)

			checkStatus(site)
		}
		time.Sleep(delay * time.Second)
		fmt.Println("")
	}
}

func readWebsite() []string {

	var sites []string

	file, err := os.Open("sites.txt")
	if err != nil {
		fmt.Println("error", err)
	}

	scanner := bufio.NewReader(file)

	for {
		line, err := scanner.ReadString('\n')
		line = strings.TrimSpace(line)
		sites = append(sites, line)

		if err == io.EOF {
			break
		}
	}

	file.Close()

	fmt.Println(sites)

	return sites
}

func log(site string, status bool) {
	file, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		fmt.Println(err)
	}

	file.WriteString(time.Now().Format("Mon, 2 Jan 2006 15:04") + " - " + site + " | online " + strconv.FormatBool(status) + "\n")

	file.Close()
}

func printLogs() {
	fmt.Println("Displaying logs...")
	file, err := ioutil.ReadFile("log.txt")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(file))
}
