module github.com/Rabiann/weather-mailer/pkg/go/weather

go 1.24.3

require github.com/stretchr/testify v1.10.0

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	golang.org/x/net v0.40.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250528174236-200df99c418a // indirect
	google.golang.org/protobuf v1.36.6 // indirect
)

require (
	github.com/Rabiann/weather-mailer/lib/go/config v0.0.0-00010101000000-000000000000
	github.com/Rabiann/weather-mailer/lib/go/logger v0.0.0-00010101000000-000000000000
	github.com/Rabiann/weather-mailer/protos/weather v0.0.0-00010101000000-000000000000
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/redis/go-redis/v9 v9.12.0
	github.com/stretchr/objx v0.5.2 // indirect
	google.golang.org/grpc v1.74.2
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/Rabiann/weather-mailer/lib/go/config => ../../../../lib/go/config
	github.com/Rabiann/weather-mailer/lib/go/logger => ../../../../lib/go/logger
	github.com/Rabiann/weather-mailer/protos/user => ../../../../protos/user/v1
	github.com/Rabiann/weather-mailer/protos/weather => ../../../../protos/weather/v1
)
