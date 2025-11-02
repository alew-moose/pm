package main

import (
	"fmt"
	"log"
	"os"

	"github.com/alew-moose/pm/internal/sftp"
	"github.com/alew-moose/pm/internal/uploader"
)

func main() {
	log.SetFlags(0)

	if len(os.Args) != 3 {
		printUsage()
		os.Exit(1)
	}

	cmd, cmdConfigFile := os.Args[1], os.Args[2]
	if cmd != "create" && cmd != "update" {
		printUsage()
		os.Exit(1)
	}

	home := os.Getenv("HOME")
	if home == "" {
		log.Fatal("HOME is empty")
	}
	configFile := fmt.Sprintf("%s/.pm.json", home)

	sftpConfig, err := sftp.ConfigFromFile(configFile)
	if err != nil {
		log.Fatalf("failed to load config: %s", err)
	}

	sftpClient, err := sftp.NewClient(sftpConfig)
	if err != nil {
		log.Fatalf("new sftp client: %s", err)
	}

	switch cmd {
	case "create":
		upload(sftpClient, cmdConfigFile)
	case "update":
		download(sftpClient, cmdConfigFile)
	}

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

func printUsage() {
	usageStr := fmt.Sprintf(
		"Usage:\n"+
			"  %[1]s create <create-config-file>\n"+
			"  %[1]s update <update-config-file>\n",
		os.Args[0],
	)
	fmt.Fprintf(os.Stderr, usageStr)
}
