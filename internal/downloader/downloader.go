package downloader

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/alew-moose/pm/internal/pkg"
	"github.com/alew-moose/pm/internal/sftp"
)

type PackageDownloader struct {
	config     *Config
	sftpClient *sftp.Client
}

func NewPackageDownloader(config *Config, sftpClient *sftp.Client) (*PackageDownloader, error) {
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
	log.Printf("download packages: %s\n", stringersSliceToString(d.config.Packages))
	packages, err := d.findPackages()
	if err != nil {
		return fmt.Errorf("find packages: %s", err)
	}

	for _, p := range packages {
		archivePath, err := d.sftpClient.DownloadPackage(p)
		if err != nil {
			return fmt.Errorf("download package: %s", err)
		}

		log.Printf("extracting %s\n", archivePath)
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
	found := make(map[pkg.PackageVersionSpec]pkg.PackageVersion)
	for _, file := range files {
		pv, err := pkg.PackageVersionFromString(file.Name())
		if err != nil {
			log.Printf("invalid package name %q: %s, skipping\n", file.Name(), err)
			continue
		}
		for _, pvs := range d.config.Packages {
			if !pvs.Match(pv) {
				continue
			}
			if foundPV, ok := found[pvs]; !ok || pvs.VersionSpec.Version.GreaterThan(foundPV.Version) {
				log.Printf("found packages for %s: %s\n", pvs, pv)
				found[pvs] = pv
			}
		}
	}

	var notFound []pkg.PackageVersionSpec
	foundPackages := make(map[pkg.PackageVersion][]pkg.PackageVersionSpec, len(found))
	packages := make([]string, 0, len(foundPackages))
	for _, pvs := range d.config.Packages {
		if pv, ok := found[pvs]; ok {
			if len(foundPackages[pv]) == 0 {
				packages = append(packages, pv.String())
			}
			foundPackages[pv] = append(foundPackages[pv], pvs)
		} else {
			notFound = append(notFound, pvs)
		}
	}

	if len(notFound) > 0 {
		return nil, fmt.Errorf("packages not found: %s", stringersSliceToString(notFound))
	}

	for pv, pvss := range foundPackages {
		if len(pvss) > 1 {
			log.Printf("package %s satisfies several specs (%s), but will be downloaded and extracted only once\n", pv, stringersSliceToString(pvss))
		}
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

	createdDirs := make(map[string]struct{})
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err == tar.ErrInsecurePath {
			log.Printf("insecure path: %q, skipping\n", header.Name)
			continue
		}
		if err != nil {
			return fmt.Errorf("tar: %s", err)
		}

		dir := filepath.Dir(header.Name)
		if _, ok := createdDirs[dir]; !ok {
			log.Printf("creating dir %q\n", dir)
			if err := os.MkdirAll(dir, 0755); err != nil {
				return fmt.Errorf("mkdir: %s", err)
			}
			createdDirs[dir] = struct{}{}
		}

		log.Printf("extracting file %q\n", header.Name)

		if _, err := os.Stat(header.Name); !errors.Is(err, os.ErrNotExist) {
			log.Printf("%q already exists, overwriting\n", header.Name)
		}

		f, err := os.OpenFile(header.Name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, header.FileInfo().Mode())
		if err != nil {
			return err
		}
		defer func() {
			_ = f.Close()
		}()

		bytes, err := io.Copy(f, tr)
		if err != nil {
			return fmt.Errorf("copy: %s", err)
		}

		if err := f.Close(); err != nil {
			return fmt.Errorf("close file: %s", err)
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
