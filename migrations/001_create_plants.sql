-- +goose Up
CREATE TABLE plants (
	trefle_id INTEGER PRIMARY KEY,
	common_name TEXT NOT NULL,
	scientific_name TEXT NOT NULL,
	family TEXT NOT NULL,
	genus TEXT NOT NULL,
	status TEXT NOT NULL,
	year INTEGER NOT NULL,
	image_url TEXT NOT NULL
);

-- +goose Down
DROP TABLE plants;
