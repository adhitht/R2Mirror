package processor

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/adhitht/R2Mirror/internal/config"
	"github.com/adhitht/R2Mirror/internal/logger"
	"github.com/adhitht/R2Mirror/internal/storage"
	"github.com/fsnotify/fsnotify"
)

type ReleaseEntry struct {
	Version     string
	Key         string
	LastUpdated string
	Filename    string
}

type FileEntry struct {
	Filename    string
	LastUpdated string
}

type Processor struct {
	storage storage.R2Client
	logger  *logger.Logger
}

func New(storage *storage.R2Client, logger *logger.Logger) *Processor {
	return &Processor{
		storage: *storage,
		logger:  logger,
	}
}

func (p *Processor) ProcessReleases(cfg *config.Config) error {
	var releases []ReleaseEntry

	p.logger.Info("Processing Ubuntu releases", "count", len(cfg.Releases))

	for _, version := range cfg.Releases {
		if err := p.processVersion(cfg, version, &releases); err != nil {
			p.logger.Error("Failed to process version", "version", version, "error", err)
			continue
		}
	}

	if len(releases) == 0 {
		return fmt.Errorf("no releases were successfully processed")
	}

	lastUpdated := time.Now().Format("Monday, January 2, 2006 at 3:04 PM MST")
	fmt.Println("RELEASES: ", releases)
	if err := p.generateIndex(cfg, releases, lastUpdated); err != nil {
		return fmt.Errorf("failed to generate index: %w", err)
	}

	p.logger.Info("Processing complete", "successful", len(releases))
	return nil
}

func (p *Processor) WatchConfig(cfg *config.Config) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("failed to create watcher: %w", err)
	}
	defer watcher.Close()

	if err := watcher.Add(config.ConfigFile); err != nil {
		return fmt.Errorf("failed to watch config file: %w", err)
	}

	p.logger.Info("Watching config file for changes...")

	for {
		select {
		case event := <-watcher.Events:
			if event.Op&(fsnotify.Write|fsnotify.Create) > 0 {
				time.Sleep(config.DebounceDelay)
				p.logger.Info("Config changed, reloading...")

				newCfg, err := config.Load()
				if err != nil {
					p.logger.Error("Failed to reload config", "error", err)
					continue
				}

				if err := p.ProcessReleases(newCfg); err != nil {
					p.logger.Error("Failed to process updated config", "error", err)
				}
			}
		case err := <-watcher.Errors:
			p.logger.Error("Watcher error", "error", err)
		}
	}
}

func (p *Processor) processVersion(cfg *config.Config, version string, releases *[]ReleaseEntry) error {
	baseURL := fmt.Sprintf("https://releases.ubuntu.com/%s/", version)

	p.logger.Info("Fetching release directory", "version", version, "url", baseURL)

	// Get the HTML directory listing
	resp, err := http.Get(baseURL)
	if err != nil {
		return fmt.Errorf("failed to fetch release directory: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("directory fetch failed: HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read directory listing: %w", err)
	}

	// Regex to extract href links
	re := regexp.MustCompile(`href="([^"]+)"`)
	matches := re.FindAllStringSubmatch(string(body), -1)

	if len(matches) == 0 {
		return fmt.Errorf("no files found in release directory for %s", version)
	}

	var uploadedFiles []string
	for _, m := range matches {
		fileName := m[1]

		// Skip parent directory and subfolders
		if fileName == "../" || strings.HasSuffix(fileName, "/") {
			continue
		}

		fileURL := baseURL + fileName
		key := fmt.Sprintf("%s/%s", version, fileName) // version folder in R2

		p.logger.Info("Downloading file", "file", fileName, "url", fileURL)

		fResp, err := http.Get(fileURL)
		if err != nil {
			p.logger.Error("Download failed", "file", fileName, "error", err)
			continue
		}
		if fResp.StatusCode != http.StatusOK {
			p.logger.Error("Download failed (bad status)", "file", fileName, "status", fResp.StatusCode)
			fResp.Body.Close()
			continue
		}

		// Upload to R2
		err = p.storage.UploadFile(context.TODO(), cfg.Bucket, key, fResp.Body, "application/octet-stream")
		fResp.Body.Close()
		if err != nil {
			p.logger.Error("Upload failed", "file", fileName, "error", err)
			continue
		}

		uploadedFiles = append(uploadedFiles, fileName)
		p.logger.Info("Uploaded file", "key", key)
	}

	if len(uploadedFiles) == 0 {
		return fmt.Errorf("no files uploaded for version %s", version)
	}

	// Generate and upload index.html for the version
	if err := p.generateVersionIndex(cfg, version, uploadedFiles); err != nil {
		return fmt.Errorf("failed to generate version index: %w", err)
	}

	*releases = append(*releases, ReleaseEntry{
		Version:     version,
		Key:         fmt.Sprintf("%s/", version),
		LastUpdated: time.Now().Format("Monday, January 2, 2006 at 3:04 PM MST"),
	})

	p.logger.Info("Successfully processed release", "version", version, "files", len(uploadedFiles))
	return nil
}

func (p *Processor) generateVersionIndex(cfg *config.Config, version string, files []string) error {
	tmplBytes, err := os.ReadFile("templates/version_index.html")
	if err != nil {
		return err
	}

	tmpl, err := template.New("version_index").Parse(string(tmplBytes))
	if err != nil {
		return err
	}

	var fileEntries []FileEntry
	for _, f := range files {
		fileEntries = append(fileEntries, FileEntry{
			Filename:    f,
			LastUpdated: time.Now().Format("Monday, January 2, 2006 at 3:04 PM MST"),
		})
	}

	data := struct {
		Version string
		Files   []FileEntry
	}{
		Version: version,
		Files:   fileEntries,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return err
	}

	key := fmt.Sprintf("%s/index.html", version)
	return p.storage.UploadHTML(context.TODO(), cfg.Bucket, key, buf.String())
}
