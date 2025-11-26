package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/AleGaliev/runtimemetrics/internal/audit"
)

const (
	auditSaveFile = "auditSaveFile"
)

type AuditSaveFile struct {
	path string
	name string
	mu   sync.Mutex
}

func CreateAuditSaveFile(path string) (*AuditSaveFile, error) {
	if path == "" {
		return nil, fmt.Errorf("path cannot be empty")
	}
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o666)
	if err != nil {
		return nil, fmt.Errorf("could not open metrics file: %w", err)
	}
	file.Close()

	return &AuditSaveFile{path: path, name: auditSaveFile}, nil
}

func (a *AuditSaveFile) SendAudit(message audit.Audit) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}
	data = append(data, '\n')
	file, err := os.OpenFile(a.path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0o666)
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
