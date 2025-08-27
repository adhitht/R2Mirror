package processor

import (
	"bytes"
	"context"
	"html/template"
	"os"

	"github.com/adhitht/R2Mirror/internal/config"
)

func (p *Processor) generateIndex(cfg *config.Config, releases []ReleaseEntry, lastUpdated string) error {
	tmplBytes, err := os.ReadFile("templates/index.html")
	if err != nil {
		return err
	}

	tmpl, err := template.New("index").Parse(string(tmplBytes))
	if err != nil {
		return err
	}

	data := struct {
		Releases    []ReleaseEntry
		LastUpdated string
	}{
		Releases:    releases,
		LastUpdated: lastUpdated,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return err
	}

	return p.storage.UploadHTML(context.TODO(), cfg.Bucket, "index.html", buf.String())
}
