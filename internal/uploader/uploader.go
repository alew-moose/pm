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

	"github.com/alew-moose/pm/internal/sftp"
)

type PackageUploader struct {
	config     *Config
	sftpClient *sftp.Client
}

func New(config *Config, sftpClient *sftp.Client) (*PackageUploader, error) {
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %s", err)
	}
	pu := &PackageUploader{
		config:     config,
		sftpClient: sftpClient,
	}
	return pu, nil
}

func (u *PackageUploader) Upload() error {
	// TODO: create dir unless exists
	packageName := u.config.FileName()
	packageExists, err := u.sftpClient.PackageExists(packageName)
	if err != nil {
		return fmt.Errorf("check if package exists: %s", err)
	}
	if packageExists {
		return fmt.Errorf("package %s already exists", packageName)
	}

	paths, err := u.getPaths()
	if err != nil {
		return fmt.Errorf("get paths: %s", err)
	}

	// XXX
	for _, p := range paths {
		fmt.Println(p)
	}

	archivePath, err := u.createArchive(paths)
	if err != nil {
		return fmt.Errorf("create archive: %s", err)
	}
	defer func() {
		err := os.Remove(archivePath)
		// TODO: Remove error check ?
		if err != nil {
			log.Printf("failed to remove %q: %s", archivePath, err)
		}
	}()

	// XXX log verbose
	log.Printf("archive: %s\n", archivePath)

	if err := u.sftpClient.UploadPackage(u.config.FileName(), archivePath); err != nil {
		return fmt.Errorf("sftp client upload package: %s", err)
	}

	return nil
}

func (u *PackageUploader) getPaths() ([]string, error) {
	seen := make(map[string]struct{})
	var paths []string
	for _, target := range u.config.Targets {
		files, err := filepath.Glob(target.Path)
		if err != nil {
			return nil, fmt.Errorf("glob %q: %s", target.Path, err)
		}
		files, err = filterPaths(files, target.Exclude)
		if err != nil {
			return nil, fmt.Errorf("filter paths: %s", err)
		}
		for _, file := range files {
			if _, ok := seen[file]; ok {
				// TODO: verbose
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
		if !excludeRe.MatchString(path) {
			filtered = append(filtered, path)
		}
	}
	return filtered, nil
}

// TODO: regexp.Replace -> strings.Replace
var globStarRe = regexp.MustCompile(`\\\*`)
var globQuestionRe = regexp.MustCompile(`\\\?`)

func excludeStrToRegexp(exclude string) (*regexp.Regexp, error) {
	exclude = regexp.QuoteMeta(exclude)
	exclude = globStarRe.ReplaceAllString(exclude, ".*")
	exclude = globQuestionRe.ReplaceAllString(exclude, ".?")
	exclude = "^" + exclude + "$"
	return regexp.Compile(exclude)
}

func (u *PackageUploader) createArchive(paths []string) (string, error) {
	tmpFilePattern := fmt.Sprintf("%s-%s-*.tar.gz", u.config.Name, u.config.Version)
	f, err := os.CreateTemp("", tmpFilePattern)
	if err != nil {
		return "", fmt.Errorf("create temp file: %s", err)
	}
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
		return fmt.Errorf("writer header: %s", err)
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
