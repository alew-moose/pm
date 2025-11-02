package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/alew-moose/pm/internal/downloader"
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

	sftpClient, err := newSftpClient()
	if err != nil {
		log.Fatalf("new sftp client: %s", err)
	}

	switch cmd {
	case "create":
		if err := upload(sftpClient, cmdConfigFile); err != nil {
			log.Fatalf("failed to upload: %s", err)
		}
	case "update":
		if err := download(sftpClient, cmdConfigFile); err != nil {
			log.Fatalf("failed to download: %s", err)
		}
	}
}

func newSftpClient() (*sftp.Client, error) {
	home := os.Getenv("HOME")
	if home == "" {
		return nil, errors.New("HOME is empty")
	}
	configFile := fmt.Sprintf("%s/.pm.json", home)

	sftpConfig, err := sftp.ConfigFromFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("load config: %s", err)
	}

	sftpClient, err := sftp.NewClient(sftpConfig)
	if err != nil {
		return nil, err
	}

	return sftpClient, nil
}

func upload(sftpClient *sftp.Client, cmdConfigFile string) error {
	config, err := uploader.ConfigFromFile(cmdConfigFile)
	if err != nil {
		return fmt.Errorf("parse uploader config: %s", err)
	}

	uploader, err := uploader.New(config, sftpClient)
	if err != nil {
		return fmt.Errorf("create new uploader: %s", err)
	}

	if err := uploader.Upload(); err != nil {
		return fmt.Errorf("upload: %s", err)
	}

	return nil
}

func download(sftpClient *sftp.Client, cmdConfigFile string) error {
	config, err := downloader.ConfigFromFile(cmdConfigFile)
	if err != nil {
		return fmt.Errorf("parse downloader config: %s", err)
	}

	fmt.Printf("config: %#v\n", config)

	// downloader, err := downloader.New(config, sftpClient)
	// if err != nil {
	// 	return fmt.Errorf("create new downloader: %s", err)
	// }

	// if err := downloader.Download(); err != nil {
	// 	return fmt.Errorf("download: %s", err)
	// }

	return nil
}

func printUsage() {
	usageStr := fmt.Sprintf(
		"Usage:\n"+
			"\t%[1]s create <create-config-file>\n"+
			"\t%[1]s update <update-config-file>\n",
		os.Args[0],
	)
	fmt.Fprintln(os.Stderr, usageStr)
}
