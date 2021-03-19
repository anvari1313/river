package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCMD = &cobra.Command{
	Use:   "river",
	Short: "CLI for publishing messages to RabbitMQ",
	Long:  `River is a simple CLI to generate a flow of messages from different data stores to a channel of RabbitMQ`,
}

var (
	rabbitURI string
	rabbitEx  string

	// GitCommit is SHA1 ref of current build
	GitCommit string

	// GitRef is branch name of current build
	GitRef string

	// GitRef is branch name of current build
	GitTag string

	// BuildDate is the timestamp of build
	BuildDate string

	asciiArt = `
 _______ _________          _______  _______ 
(  ____ )\__   __/|\     /|(  ____ \(  ____ )
| (    )|   ) (   | )   ( || (    \/| (    )|
| (____)|   | |   | |   | || (__    | (____)|
|     __)   | |   ( (   ) )|  __)   |     __)
| (\ (      | |    \ \_/ / | (      | (\ (   
| ) \ \_____) (___  \   /  | (____/\| ) \ \__
|/   \__/\_______/   \_/   (_______/|/   \__/
`
)

func init() {
	cobra.OnInitialize(configure)

	rootCMD.PersistentFlags().StringVar(&rabbitURI, "rabbit-uri", "amqp://guest:guest@localhost:5672", "RabbitMQ connection URI")
	rootCMD.PersistentFlags().StringVar(&rabbitEx, "rabbit-ex", "exchange1", "RabbitMQ exchange name")

	rootCMD.AddCommand(streamCMD)
}

func configure() {
	fmt.Print(asciiArt)
	fmt.Printf("Build info: %s [%s/%s] @ %s\n\n", GitTag, GitRef, GitCommit, BuildDate)
}

// Execute executes the root command.
func Execute() error {
	return rootCMD.Execute()
}
