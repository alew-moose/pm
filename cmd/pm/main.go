package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/alew-moose/pm/internal/sftp"
	"github.com/alew-moose/pm/internal/uploader"
)

func main() {
	log.SetFlags(0)

	configFile, cmd, cmdConfigFile := parseArgs()

	sftpClient := newSftpClient(configFile)

	switch cmd {
	case "create":
		upload(sftpClient, cmdConfigFile)
	case "update":
		download(sftpClient, cmdConfigFile)
	}

}

func newSftpClient(configFile string) *sftp.Client {
	sftpConfig, err := sftp.ConfigFromFile(configFile)
	if err != nil {
		log.Fatalf("failed to load config: %s", err)
	}

	sftpClient, err := sftp.NewClient(sftpConfig)
	if err != nil {
		log.Fatalf("new sftp client: %s", err)
	}

	return sftpClient
}

func upload(sftpClient *sftp.Client, cmdConfigFile string) {
	uploaderConfig, err := uploader.ConfigFromFile(cmdConfigFile)
	if err != nil {
		log.Fatalf("failed to parse uploader config: %s", err)
	}

	uploader, err := uploader.New(uploaderConfig, sftpClient)
	if err != nil {
		log.Fatalf("failed to create new uploader: %s", err)
	}

	if err := uploader.Upload(); err != nil {
		log.Fatalf("failed to upload: %s", err)
	}
}

func download(sftpClient *sftp.Client, configFile string) {
}

func parseArgs() (string, string, string) {
	defaultConfigFile := fmt.Sprintf("%s/.pm.json", os.Getenv("HOME"))

	var configFile string
	flag.StringVar(&configFile, "f", defaultConfigFile, "config file")

	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		printUsage()
		os.Exit(1)
	}

	if args[0] != "create" && args[0] != "update" {
		printUsage()
		os.Exit(1)
	}

	return configFile, args[0], args[1]
}

func printUsage() {
	usageStr := fmt.Sprintf(
		"Usage:\n"+
			"  %[1]s [-f <config-file>] create <create-config-file>\n"+
			"  %[1]s [-f <config-file>] update <update-config-file>\n",
		os.Args[0],
	)
	fmt.Fprintf(os.Stderr, usageStr)
}
