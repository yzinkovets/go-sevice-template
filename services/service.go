package services

import (
	"context"
	"fmt"
	"time"

	"go-service-template/config"
	"go-service-template/db"
)

type SomeService struct {
	db *db.DbConnection
}

func NewSomeService(cfg *config.MainConfig, db *db.DbConnection) (*SomeService, error) {
	s := &SomeService{
		db: db,
	}
	return s, nil
}

// Search returns iot data by mac, hostname, dhcpFingerprint
// if [force] is true, then search in external service, otherwise search in cache first
func (s *SomeService) Call(gwUUID string, mac string) (string, error) {
	return "OK", nil
}

func (s *SomeService) dbCall() error {
	query := `
	CREATE TABLE IF NOT EXISTS public.service (
		id serial4 PRIMARY KEY,
		name text not null,
		mac macaddr not null,
		CONSTRAINT service_mac UNIQUE (mac)
	);
	`
	db, err := s.db.Conn()
	if err != nil {
		return fmt.Errorf("can't get connection: %w", err)
	}

	if _, err := db.Query(context.Background(), query); err != nil {
		return fmt.Errorf("can't create service table: %w", err)
	}

	return nil
}

func (s *SomeService) dbSelect(mac string) (string, error) {
	query := `select id, name from public.fingerprints where mac = $1`

	name := ""

	db, err := s.db.Conn()
	if err != nil {
		return "", fmt.Errorf("can't get connection: %w", err)
	}

	r := db.QueryRow(context.Background(), query, mac)
	if err := r.Scan(&name); err != nil {
		return "", fmt.Errorf("can't execute query: %w", err)
	}

	return name, nil
}

func (s *SomeService) dbInsert(mac, name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.db.InsertTimeoutSec)*time.Second)
	defer cancel()

	query := "INSERT INTO public.service (mac, name) VALUES ($1, $2)"

	db, err := s.db.Conn()
	if err != nil {
		return fmt.Errorf("can't get connection: %w", err)
	}

	_, err = db.Exec(ctx, query, mac, name)

	return err
}
