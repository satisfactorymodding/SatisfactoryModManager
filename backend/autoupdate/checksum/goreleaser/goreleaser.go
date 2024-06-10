package goreleaser

import (
	"encoding/hex"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"text/template"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend/autoupdate/updater"
)

type checksumSource struct {
	checksumArtifactTemplate *template.Template
	split                    bool
}

func MakeGoreleaserChecksumSource(checksumArtifactFormat string, split bool) updater.ChecksumSource {
	return &checksumSource{
		checksumArtifactTemplate: template.Must(template.New("checksumArtifact").Parse(checksumArtifactFormat)),
		split:                    split,
	}
}

func (g *checksumSource) GetChecksumForFile(source updater.Source, version string, filename string) ([]byte, error) {
	var checksumFilenameBuilder strings.Builder
	err := g.checksumArtifactTemplate.Execute(&checksumFilenameBuilder, map[string]string{"ArtifactName": filename, "Version": version})
	if err != nil {
		return nil, fmt.Errorf("failed to build checksum filename: %w", err)
	}
	checksumFilename := checksumFilenameBuilder.String()
	chesumFile, _, err := source.GetFile(version, checksumFilename)
	if err != nil {
		return nil, fmt.Errorf("failed to get checksum file: %w", err)
	}
	defer chesumFile.Close()
	checksum, err := io.ReadAll(chesumFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read checksum: %w", err)
	}

	if g.split {
		// Checksum file will only contain one hex string
		sum, err := hex.DecodeString(strings.TrimSpace(string(checksum)))
		if err != nil {
			return nil, fmt.Errorf("failed to decode checksum: %w", err)
		}
		return sum, nil
	}

	// Checksum file will contain multiple lines of {artifact} {hash}
	checksums := parseChecksumFile(checksum)
	if sum, ok := checksums[filename]; ok {
		return sum, nil
	}
	return nil, fmt.Errorf("failed to find checksum for file")
}

func parseChecksumFile(checksumFile []byte) map[string][]byte {
	checksums := make(map[string][]byte)
	lines := strings.Split(string(checksumFile), "\n")
	for _, line := range lines {
		if len(line) == 0 {
			// Skip empty lines
			continue
		}
		parts := strings.Split(line, "  ")
		if len(parts) != 2 {
			slog.Debug("invalid checksum entry", slog.String("entry", line))
			continue
		}
		hexSum := parts[0]
		filename := parts[1]
		sum, err := hex.DecodeString(hexSum)
		if err != nil {
			slog.Debug("failed to decode checksum", slog.String("checksum", hexSum), slog.String("filename", filename), slog.Any("error", err))
			continue
		}
		checksums[parts[1]] = sum
	}
	return checksums
}
