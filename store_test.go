package main

import (
	"context"
	"database/sql"
	"testing"
)

func TestPlantStoreAddAndReadAllPlants(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	db.SetMaxOpenConns(1)

	_, err = db.Exec(`
		CREATE TABLE plants (
			trefle_id INTEGER PRIMARY KEY,
			common_name TEXT NOT NULL,
			scientific_name TEXT NOT NULL,
			family TEXT NOT NULL,
			genus TEXT NOT NULL,
			status TEXT NOT NULL,
			year INTEGER NOT NULL,
			image_url TEXT NOT NULL
		)
	`)
	if err != nil {
		t.Fatal(err)
	}

	plants := []fullPlant{
		{
			ID:             42,
			CommonName:     "Example plant",
			ScientificName: "Planta exemplaris",
			Family:         plantTaxon{Name: "Exampleaceae"},
			Genus:          plantTaxon{Name: "Planta"},
			Status:         "accepted",
			Year:           2024,
			ImageURL:       "https://example.com/plant.jpg",
		},
		{
			ID:             11,
			CommonName:     "Second plant",
			ScientificName: "Planta secunda",
			Family:         plantTaxon{Name: "Otheraceae"},
			Genus:          plantTaxon{Name: "Othera"},
			Status:         "accepted",
			Year:           2025,
			ImageURL:       "https://example.com/second-plant.jpg",
		},
	}
	store := &plantStore{db: db}
	for _, plant := range plants {
		if err := store.AddPlant(context.Background(), plant); err != nil {
			t.Fatal(err)
		}
	}

	got, err := store.ReadAllPlants(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != len(plants) {
		t.Fatalf("got %d plants, want %d", len(got), len(plants))
	}
	want := []fullPlant{plants[1], plants[0]}
	for index := range want {
		if got[index] != want[index] {
			t.Errorf("plant %d = %#v, want %#v", index, got[index], want[index])
		}
	}
}
