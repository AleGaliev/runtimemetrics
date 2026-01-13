package repository

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/AleGaliev/runtimemetrics/internal/audit"
)

const (
	auditSenderName = "auditSender"
)

//generate:reset
type AuditSender struct {
	client http.Client
	url    url.URL
	name   string
	mu     sync.Mutex
}

func NewAuditSender(addrAuditServer string) (*AuditSender, error) {
	if addrAuditServer == "" {
		return nil, errors.New("address audit server is empty")
	}
	auditServer, err := url.Parse(addrAuditServer)
	if err != nil {
		return nil, fmt.Errorf("could not parse url: %w", err)
	}
	return &AuditSender{
		client: http.Client{
			Timeout: 2 * time.Second,
		},
		url: url.URL{
			Scheme: auditServer.Scheme,
			Host:   auditServer.Host,
			Path:   auditServer.Path,
		},
		name: auditSenderName,
	}, nil
}

func (a *AuditSender) SendAudit(message audit.Audit) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	jsonData, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("could not marshal json: %w", err)
	}
	resp, err := a.client.Post(a.url.String(), "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("could not send audit: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non 200 status code: %d", resp.StatusCode)
	}
	return nil
}

func (a *AuditSender) GetID() string {
	return a.name
}
