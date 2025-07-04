module github.com/jurshsmith/vaultstream/seeder

go 1.24.1

require (
	github.com/jurshsmith/vaultstream/config v0.0.0
	github.com/jurshsmith/vaultstream/database v0.0.0
	github.com/jurshsmith/vaultstream/logger v0.0.0
	github.com/lib/pq v1.10.9
	go.uber.org/zap v1.27.0
)

replace github.com/jurshsmith/vaultstream/config => ../config

replace github.com/jurshsmith/vaultstream/database => ../database

replace github.com/jurshsmith/vaultstream/logger => ../logger

require (
	ariga.io/atlas v0.31.1-0.20250212144724-069be8033e83 // indirect
	entgo.io/ent v0.14.4 // indirect
	github.com/agext/levenshtein v1.2.1 // indirect
	github.com/apparentlymart/go-textseg/v13 v13.0.0 // indirect
	github.com/apparentlymart/go-textseg/v15 v15.0.0 // indirect
	github.com/bmatcuk/doublestar v1.3.4 // indirect
	github.com/go-openapi/inflect v0.19.0 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/hashicorp/hcl/v2 v2.13.0 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/mitchellh/go-wordwrap v0.0.0-20150314170334-ad45545899c7 // indirect
	github.com/zclconf/go-cty v1.14.4 // indirect
	github.com/zclconf/go-cty-yaml v1.1.0 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	golang.org/x/mod v0.23.0 // indirect
	golang.org/x/text v0.23.0 // indirect
)
