module github.com/sks/kihocche

go 1.23.5
require (
	github.com/arran4/golang-ical v0.3.2
	github.com/drone/go-scm v1.39.1
	github.com/google/uuid v1.6.0
	github.com/hashicorp/go-retryablehttp v0.7.7
	github.com/spf13/cobra v1.9.1
	gocloud.dev v0.41.0
	golang.org/x/sync v0.12.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/golang/groupcache v0.0.0-20241129210726-2c02b8208cf8 // indirect
	github.com/googleapis/gax-go/v2 v2.14.1 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
	go.opencensus.io v0.24.0 // indirect
	golang.org/x/crypto v0.36.0 // indirect
	golang.org/x/net v0.38.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	golang.org/x/xerrors v0.0.0-20240903120638-7835f813f4da // indirect
	google.golang.org/api v0.228.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250324211829-b45e905df463 // indirect
	google.golang.org/grpc v1.71.0 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
)

replace (
	github.com/drone/go-scm => github.com/appcd-dev/go-scm v0.0.0-20241009172542-a16030046ecd
	golang.org/x/crypto => golang.org/x/crypto v0.32.0
	golang.org/x/net => golang.org/x/net v0.34.0
)
