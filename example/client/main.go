package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	. "github.com/kkdai/coapmq"
	"github.com/spf13/cobra"
)

func toggleLogging(enable bool) {
	if enable {
		log.SetOutput(os.Stdout)
	} else {
		log.SetOutput(ioutil.Discard)
	}
}

func printConsole() {

	fmt.Println("Command: C:Create S:Subscription P:Publish R:RemoveTopic:>")
	fmt.Printf(":>")
}

func main() {

	var serverAddr string
	var verbose bool

	rootCmd := &cobra.Command{
		Use:   "coapmq_client",
		Short: "Client to connect to coapmq broker",
		Run: func(ccmd *cobra.Command, args []string) {

			toggleLogging(verbose)

			fmt.Println("Connect to coapmq server:", serverAddr)
			client := NewClient(serverAddr)
			if client == nil {
				fmt.Fatalln("Cannot connect to server, please check your setting.")
			}
			quit := false
			scanner := bufio.NewScanner(os.Stdin)
			printConsole()
			for !quit {

				var topic, msg string

				if !scanner.Scan() {
					break
				}
				line := scanner.Text()
				parts := strings.Split(line, " ")
				cmd := parts[0]
				if len(parts) > 1 {
					topic = parts[1]
				}
				if len(parts) > 2 {
					msg = parts[2]
				}

				fmt.Println("cmd:", cmd, " topic:", topic, " msg:", msg)
				var err error
				switch cmd {
				case "C", "c": //CREATE TOPIC
					err = client.CreateTopic(topic)
					fmt.Println("CreateTopic topic:", topic, " ret=", err)
				case "S", "s": //SUBSCRIPTION
					ch, err := client.Subscription(topic)
					fmt.Println("Subscription topic:", topic, " ret=", err)
					go func() {
						fmt.Println("Got pub from topic:", topic, " pub:", <-ch)
					}()
				case "P", "p": //PUBLISH
					err = client.Publish(topic, msg)
				case "R", "r": //REMOVE
					err = client.RemoveTopic(topic)
				case "Q", "q":
					quit = true
				case "V", "v":
					verbose = !verbose
					toggleLogging(verbose)
				default:
					fmt.Println("Command not support.")
				}
				printConsole()
			}
		},
	}
	rootCmd.Flags().StringVarP(&serverAddr, "server", "s", "localhost:5683", "coapmq server address")
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose")
	rootCmd.Execute()
}
