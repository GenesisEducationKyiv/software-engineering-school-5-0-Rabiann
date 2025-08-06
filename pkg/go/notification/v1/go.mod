module github.com/Rabiann/weather-mailer/pkg/go/notification

go 1.24.3

require (
	github.com/IBM/sarama v1.45.2
	github.com/Rabiann/weather-mailer/lib/go/config v0.0.0-00010101000000-000000000000
	github.com/Rabiann/weather-mailer/lib/go/logger v0.0.0-00010101000000-000000000000
)

require (
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	github.com/sendgrid/rest v2.6.9+incompatible // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/eapache/go-resiliency v1.7.0 // indirect
	github.com/eapache/go-xerial-snappy v0.0.0-20230731223053-c322873962e3 // indirect
	github.com/eapache/queue v1.1.0 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/hashicorp/errwrap v1.0.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-uuid v1.0.3 // indirect
	github.com/jcmturner/aescts/v2 v2.0.0 // indirect
	github.com/jcmturner/dnsutils/v2 v2.0.0 // indirect
	github.com/jcmturner/gofork v1.7.6 // indirect
	github.com/jcmturner/gokrb5/v8 v8.4.4 // indirect
	github.com/jcmturner/rpc/v2 v2.0.3 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/pierrec/lz4/v4 v4.1.22 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20201227073835-cf1acfcdf475 // indirect
	github.com/sendgrid/sendgrid-go v3.16.1+incompatible
	github.com/stretchr/testify v1.10.0
	golang.org/x/crypto v0.38.0 // indirect
	golang.org/x/net v0.40.0 // indirect
)

replace (
	github.com/Rabiann/weather-mailer/lib/go/config => ../../../../lib/go/config
	github.com/Rabiann/weather-mailer/lib/go/logger => ../../../../lib/go/logger
	github.com/Rabiann/weather-mailer/protos/user => ../../../../protos/user/v1
	github.com/Rabiann/weather-mailer/protos/weather => ../../../../protos/weather/v1
)
