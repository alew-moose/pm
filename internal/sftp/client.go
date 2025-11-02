package sftp

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

type Client struct {
	client *sftp.Client
	config *Config
}

func NewClient(config *Config) (*Client, error) {
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("validate config: %s", err)
	}
	sshClient, err := sshConnect(config.Host, config.Port, config.User)
	if err != nil {
		return nil, fmt.Errorf("ssh connect: %s", err)
	}
	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		return nil, fmt.Errorf("new sftp client: %s", err)
	}
	client := &Client{
		client: sftpClient,
		config: config,
	}
	return client, nil
}

// TODO: version string->version
func (c *Client) PackageExists(name string) (bool, error) {
	path := fmt.Sprintf("%s/%s", c.config.Path, name)
	fileInfo, err := c.client.Stat(path)
	if err != nil && err.Error() == "file does not exist" {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("stat: %s", err)
	}
	// TODO: need to check fileInfo?
	_ = fileInfo
	return true, nil
}

// TODO: rename all name -> packageName ?
func (c *Client) UploadPackage(packageName string, archivePath string) error {
	log.Printf("uploading %s as package %s\n", archivePath, packageName)
	// TODO: refactor
	remotePath := fmt.Sprintf("%s/%s", c.config.Path, packageName)

	srcFile, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer func() {
		_ = srcFile.Close()
	}()

	dstFile, err := c.client.OpenFile(remotePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC)
	if err != nil {
		return fmt.Errorf("open remote file: %s", err)
	}
	defer func() {
		_ = dstFile.Close()
	}()

	bytes, err := io.Copy(dstFile, srcFile)
	if err != nil {
		return fmt.Errorf("copy: %s", err)
	}

	// TODO: log verbose
	_ = bytes

	// TODO: check close files?

	return nil
}

func (c *Client) DownloadPackage(packageName string) (string, error) {
	// TODO: refactor
	remotePath := fmt.Sprintf("%s/%s", c.config.Path, packageName)

	srcFile, err := c.client.OpenFile(remotePath, os.O_RDONLY)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = srcFile.Close()
	}()

	// TODO: move out of here
	tmpFilePattern := fmt.Sprintf("%s-*.tar.gz", packageName)
	dstFile, err := os.CreateTemp("", tmpFilePattern)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = dstFile.Close()
	}()

	log.Printf("downloading package %s to %s\n", packageName, dstFile.Name())

	bytes, err := io.Copy(dstFile, srcFile)
	if err != nil {
		return "", fmt.Errorf("copy: %s", err)
	}

	// TODO: log verbose
	_ = bytes

	// TODO: check close files?

	return dstFile.Name(), nil
}

func (c *Client) GetPackages() ([]os.FileInfo, error) {
	return c.client.ReadDir(c.config.Path)
}

func sshConnect(host, port, user string) (*ssh.Client, error) {
	socket := os.Getenv("SSH_AUTH_SOCK")
	if socket == "" {
		return nil, errors.New("SSH_AUTH_SOCK is empty")
	}
	conn, err := net.Dial("unix", socket)
	if err != nil {
		return nil, fmt.Errorf("open SSH_AUTH_SOCK: %s", err)
	}

	log.Printf("connecting to %s:%s as %s", host, port, user)

	agentClient := agent.NewClient(conn)
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeysCallback(agentClient.Signers),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", host+":"+port, config)
	if err != nil {
		return nil, fmt.Errorf("ssh dial: %s", err)
	}

	return client, nil
}

// TODO: https://sftptogo.com/blog/go-sftp/ get host key
