package main

import (
	"context"
	"database/sql"
	"fmt"
)

type plantStore struct {
	db *sql.DB
}

func (s *plantStore) AddPlant(ctx context.Context, plant fullPlant) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO plants (
			trefle_id, common_name, scientific_name, family, genus, status, year, image_url
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`,
		plant.ID,
		plant.CommonName,
		plant.ScientificName,
		plant.Family.Name,
		plant.Genus.Name,
		plant.Status,
		plant.Year,
		plant.ImageURL,
	)
	if err != nil {
		return fmt.Errorf("insert plant: %w", err)
	}

	return nil
}

func (s *plantStore) ReadAllPlants(ctx context.Context) ([]fullPlant, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT trefle_id, common_name, scientific_name, family, genus, status, year, image_url
		FROM plants
		ORDER BY trefle_id
	`)
	if err != nil {
		return nil, fmt.Errorf("query plants: %w", err)
	}
	defer rows.Close()

	plants := make([]fullPlant, 0)
	for rows.Next() {
		var plant fullPlant
		if err := rows.Scan(
			&plant.ID,
			&plant.CommonName,
			&plant.ScientificName,
			&plant.Family.Name,
			&plant.Genus.Name,
			&plant.Status,
			&plant.Year,
			&plant.ImageURL,
		); err != nil {
			return nil, fmt.Errorf("scan plant: %w", err)
		}
		plants = append(plants, plant)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read plants: %w", err)
	}

	return plants, nil
}
