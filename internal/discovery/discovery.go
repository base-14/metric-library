package discovery

import (
	"os"
	"path/filepath"
	"strings"
)

type MetadataFile struct {
	Path          string
	ComponentName string
	ComponentType string
}

type MetadataDiscovery struct {
	componentDirs []string
}

func NewMetadataDiscovery() *MetadataDiscovery {
	return &MetadataDiscovery{
		componentDirs: []string{
			"receiver",
			"processor",
			"exporter",
			"extension",
			"connector",
		},
	}
}

func (d *MetadataDiscovery) FindMetadataFiles(repoPath string) ([]MetadataFile, error) {
	var files []MetadataFile

	for _, componentDir := range d.componentDirs {
		dirPath := filepath.Join(repoPath, componentDir)

		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			continue
		}

		entries, err := os.ReadDir(dirPath)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}

			metadataPath := filepath.Join(dirPath, entry.Name(), "metadata.yaml")
			if _, err := os.Stat(metadataPath); err == nil {
				files = append(files, MetadataFile{
					Path:          metadataPath,
					ComponentName: entry.Name(),
					ComponentType: componentDir,
				})
			}
		}
	}

	return files, nil
}

func (d *MetadataDiscovery) ComponentTypeFromPath(path string) string {
	parts := strings.Split(filepath.ToSlash(path), "/")
	for i, part := range parts {
		for _, cd := range d.componentDirs {
			if part == cd && i+1 < len(parts) {
				return cd
			}
		}
	}
	return ""
}
