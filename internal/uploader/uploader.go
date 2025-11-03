package uploader

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/alew-moose/pm/internal/downloader"
	"github.com/alew-moose/pm/internal/sftp"
)

type PackageUploader struct {
	config     *Config
	sftpClient *sftp.Client
	downloader *downloader.PackageDownloader
}

func NewPackageUploader(config *Config, sftpClient *sftp.Client) (*PackageUploader, error) {
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %s", err)
	}
	pu := &PackageUploader{
		config:     config,
		sftpClient: sftpClient,
	}
	if len(config.Dependencies) > 0 {
		downloaderConfig := &downloader.Config{
			Packages: config.Dependencies,
		}
		pd, err := downloader.NewPackageDownloader(downloaderConfig, sftpClient)
		if err != nil {
			return nil, fmt.Errorf("create new downloader: %s", err)
		}
		pu.downloader = pd
	}
	return pu, nil
}

func (u *PackageUploader) Upload() error {
	packageName := u.config.FileName()
	packageExists, err := u.sftpClient.PackageExists(packageName)
	if err != nil {
		return fmt.Errorf("check if package exists: %s", err)
	}
	if packageExists {
		return fmt.Errorf("package %s already exists", packageName)
	}

	if len(u.config.Dependencies) > 0 {
		log.Println("downloading dependencies")
		if err := u.downloader.Download(); err != nil {
			return fmt.Errorf("download dependencies: %s", err)
		}
	}

	paths, err := u.getPaths()
	if err != nil {
		return fmt.Errorf("get paths: %s", err)
	}

	archivePath, err := u.createArchive(paths)
	if err != nil {
		return fmt.Errorf("create archive: %s", err)
	}
	defer func() {
		err := os.Remove(archivePath)
		if err != nil {
			log.Printf("remove %q: %s\n", archivePath, err)
		}
	}()

	if err := u.sftpClient.UploadPackage(u.config.FileName(), archivePath); err != nil {
		return fmt.Errorf("sftp client upload package: %s", err)
	}

	return nil
}

func (u *PackageUploader) getPaths() ([]string, error) {
	seen := make(map[string]struct{})
	var paths []string
	for _, target := range u.config.Targets {
		log.Printf("find files for target %q excluding %q\n", target.Path, target.Exclude)
		files, err := filepath.Glob(target.Path)
		for _, file := range files {
			log.Printf("found file %q\n", file)
		}
		if err != nil {
			return nil, fmt.Errorf("glob %q: %s", target.Path, err)
		}
		files, err = filterPaths(files, target.Exclude)
		if err != nil {
			return nil, fmt.Errorf("filter paths: %s", err)
		}
		for _, file := range files {
			if _, ok := seen[file]; ok {
				log.Printf("duplicate file %q, skipping\n", file)
				continue
			}
			seen[file] = struct{}{}
			paths = append(paths, file)
		}
	}
	return paths, nil
}

func filterPaths(paths []string, exclude string) ([]string, error) {
	if exclude == "" {
		return paths, nil
	}
	filtered := make([]string, 0, len(paths))
	excludeRe, err := excludeStrToRegexp(exclude)
	if err != nil {
		return nil, fmt.Errorf("exclude str to regexp: %s", err)
	}
	for _, path := range paths {
		if excludeRe.MatchString(path) {
			log.Printf("excluded %q by %q\n", path, excludeRe)
		} else {
			filtered = append(filtered, path)
		}
	}
	return filtered, nil
}

func excludeStrToRegexp(exclude string) (*regexp.Regexp, error) {
	exclude = regexp.QuoteMeta(exclude)
	exclude = strings.ReplaceAll(exclude, `\*`, ".*")
	exclude = strings.ReplaceAll(exclude, `\?`, ".?")
	exclude = "^" + exclude + "$"
	return regexp.Compile(exclude)
}

func (u *PackageUploader) createArchive(paths []string) (string, error) {
	tmpFilePattern := fmt.Sprintf("%s-%s-*.tar.gz", u.config.Name, u.config.Version)
	f, err := os.CreateTemp("", tmpFilePattern)
	if err != nil {
		return "", fmt.Errorf("create temp file: %s", err)
	}
	log.Printf("created archive %q\n", f.Name())
	defer func() {
		_ = f.Close()
	}()

	gzw := gzip.NewWriter(f)
	defer func() {
		_ = gzw.Close()
	}()
	tw := tar.NewWriter(gzw)
	defer func() {
		_ = tw.Close()
	}()

	for _, path := range paths {
		log.Printf("adding file %q\n", path)
		if err := u.addFile(tw, path); err != nil {
			return "", fmt.Errorf("upload %q: %s", path, err)
		}
	}

	if err := tw.Close(); err != nil {
		return "", fmt.Errorf("close tar writer: %s", err)
	}
	if err := gzw.Close(); err != nil {
		return "", fmt.Errorf("close gzip writer: %s", err)
	}
	if err := f.Close(); err != nil {
		return "", fmt.Errorf("close temp file: %s", err)
	}

	return f.Name(), nil
}

func (u *PackageUploader) addFile(tw *tar.Writer, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("stat: %s", err)
	}

	header, err := tar.FileInfoHeader(fileInfo, fileInfo.Name())
	if err != nil {
		return fmt.Errorf("file info header: %s", err)
	}
	header.Name = path

	err = tw.WriteHeader(header)
	if err != nil {
		return fmt.Errorf("write header: %s", err)
	}

	_, err = io.Copy(tw, file)
	if err != nil {
		return fmt.Errorf("copy: %s", err)
	}

	if err := file.Close(); err != nil {
		return fmt.Errorf("close %q: %s", file.Name(), err)
	}

	return nil
}
