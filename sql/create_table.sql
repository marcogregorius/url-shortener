CREATE TABLE tb_shortlinks (
	id TEXT PRIMARY KEY,
	source_url TEXT,
	visited integer,
	last_visited_at TIMESTAMPTZ,
	created_at TIMESTAMPTZ
)
