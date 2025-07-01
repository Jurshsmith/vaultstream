CREATE TABLE signatures (
    id SERIAL PRIMARY KEY,
    record_id SERIAL NOT NULL,
    key_id INT NOT NULL,
    value TEXT NOT NULL,
    inserted_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT unique_signature_value UNIQUE (value),
    FOREIGN KEY (record_id) REFERENCES records(id)
);
