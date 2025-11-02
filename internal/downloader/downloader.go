package downloader

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

	"github.com/alew-moose/pm/internal/sftp"
	"github.com/alew-moose/pm/internal/version"
)

type PackageDownloader struct {
	config     *Config
	sftpClient *sftp.Client
}

// TODO: move from here + tests?
type PackageVersion struct {
	Name    string
	Version version.Version
}

func (pv PackageVersion) Validate() error {
	if !packageNameRe.MatchString(pv.Name) {
		return fmt.Errorf("invalid package name %q", pv.Name)
	}
	if err := pv.Version.Validate(); err != nil {
		return err
	}
	return nil
}

// TODO: remove?
func (pv PackageVersion) String() string {
	return fmt.Sprintf("%s-%s", pv.Name, pv.Version)
}

var packageVersionRe = regexp.MustCompile(`^(.+)-(.+)$`)

// TODO: tests
func PackageVersionFromString(s string) (PackageVersion, error) {
	var packageVersion PackageVersion
	var err error

	matches := packageVersionRe.FindStringSubmatch(s)
	if len(matches) != 3 {
		return packageVersion, fmt.Errorf("invalid package version %q", s)
	}

	packageVersion.Name = matches[1]

	packageVersion.Version, err = version.VersionFromString(matches[2])
	if err != nil {
		return packageVersion, fmt.Errorf("invalid version: %q", err)
	}

	if err := packageVersion.Validate(); err != nil {
		return packageVersion, fmt.Errorf("invalid package version %q", s)
	}

	return packageVersion, nil
}

// TODO: rename?
func New(config *Config, sftpClient *sftp.Client) (*PackageDownloader, error) {
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %s", err)
	}
	pd := &PackageDownloader{
		config:     config,
		sftpClient: sftpClient,
	}
	return pd, nil
}

func (d *PackageDownloader) Download() error {
	packages, err := d.findPackages()
	if err != nil {
		return fmt.Errorf("find packages: %s", err)
	}

	for _, p := range packages {
		// TODO: log verbose
		archivePath, err := d.sftpClient.DownloadPackage(p)
		if err != nil {
			return fmt.Errorf("download package: %s", err)
		}

		// TODO: log verbose
		if err := d.extractArchive(archivePath); err != nil {
			return fmt.Errorf("extract archive: %s", err)
		}
	}

	return nil
}

// TODO: refactor
func (d *PackageDownloader) findPackages() ([]string, error) {
	files, err := d.sftpClient.GetPackages()
	if err != nil {
		return nil, fmt.Errorf("get packages: %s", err)
	}

	// TODO: rename all found*
	found := make(map[PackageVersionSpec]PackageVersion)
	for _, file := range files {
		pv, err := PackageVersionFromString(file.Name())
		if err != nil {
			log.Printf("invalid package name %q: %s, skipping\n", file.Name(), err)
			continue
		}
		for _, pvs := range d.config.Packages {
			if !pvs.Match(pv) {
				continue
			}
			if foundPV, ok := found[pvs]; !ok || pvs.VersionSpec.Version.GreaterThan(foundPV.Version) {
				found[pvs] = pv
			}
		}
	}

	var notFound []PackageVersionSpec
	foundPackages := make(map[PackageVersion][]PackageVersionSpec, len(found))
	for _, pvs := range d.config.Packages {
		if pv, ok := found[pvs]; ok {
			foundPackages[pv] = append(foundPackages[pv], pvs)
		} else {
			notFound = append(notFound, pvs)
		}
	}

	if len(notFound) > 0 {
		return nil, fmt.Errorf("packages not found: %s", stringersSliceToString(notFound))
	}

	packages := make([]string, 0, len(foundPackages))
	for pv, pvss := range foundPackages {
		if len(pvss) > 1 {
			log.Printf("package %s satisfies several specs (%s), but will be downloaded and extracted only once\n", pv, stringersSliceToString(pvss))
		}
		packages = append(packages, pv.String())
	}

	return packages, nil
}

func (d *PackageDownloader) extractArchive(archivePath string) error {
	archiveFile, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer func() {
		_ = archiveFile.Close()
	}()

	gzr, err := gzip.NewReader(archiveFile)
	if err != nil {
		return fmt.Errorf("gzip reader: %s", err)
	}
	defer func() {
		_ = gzr.Close()
	}()
	tr := tar.NewReader(gzr)

	for {
		// TODO: warn when overwriting files?
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err == tar.ErrInsecurePath {
			log.Printf("insecure path: %q, skipping", header.Name)
		}
		if err != nil {
			return fmt.Errorf("tar: %s", err)
		}

		dir := filepath.Dir(header.Name)
		// TODO: perm?
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("mkdir: %s", err)
		}

		f, err := os.OpenFile(header.Name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, header.FileInfo().Mode())
		if err != nil {
			return err
		}

		bytes, err := io.Copy(f, tr)
		if err != nil {
			return fmt.Errorf("copy: %s", err)
		}

		// TODO log verbose
		_ = bytes
	}

	return nil
}

func stringersSliceToString[S fmt.Stringer](stringers []S) string {
	strs := make([]string, 0, len(stringers))
	for _, stringer := range stringers {
		strs = append(strs, stringer.String())
	}
	return strings.Join(strs, ", ")
}
