package repository

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/AleGaliev/runtimemetrics/internal/audit"
)

const (
	auditSaveFile = "auditSaveFile"
)

type AuditSaveFile struct {
	path string
	name string
}

func CreateAuditSaveFile(path string) (*AuditSaveFile, error) {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return nil, fmt.Errorf("could not open metrics file: %w", err)
	}
	file.Close()

	return &AuditSaveFile{path: path, name: auditSaveFile}, nil
}

func (a *AuditSaveFile) SendAudit(message audit.Audit) error {
	data, err := json.Marshal(message)
	file, err := os.OpenFile(a.path, os.O_WRONLY, 0666)
	if err != nil {
		return fmt.Errorf("could not open metrics file: %w", err)
	}

	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("could not write metrics to file: %w", err)
	}
	file.Close()

	return nil
}

func (a *AuditSaveFile) GetID() string {
	return a.name
}
