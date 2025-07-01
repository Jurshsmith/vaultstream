module github.com/jurshsmith/vaultstream/nats

go 1.24.1

require (
	github.com/jurshsmith/vaultstream/config v0.0.0
	github.com/nats-io/nats.go v1.40.1
)

require (
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/nats-io/nkeys v0.4.9 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	golang.org/x/crypto v0.33.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/text v0.23.0 // indirect
)

replace github.com/jurshsmith/vaultstream/config => ../config
