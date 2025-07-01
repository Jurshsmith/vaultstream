CREATE TABLE records (
    -- Actual record fields
    id SERIAL PRIMARY KEY,
    inserted_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
