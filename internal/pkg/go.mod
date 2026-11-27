module github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg

go 1.27.4

require (
	emperror.dev/errors v0.8.1
	github.com/EventStore/EventStore-Client-Go v1.0.2
	github.com/Masterminds/squirrel v1.5.4
	github.com/TylerBrock/colorjson v0.0.0-20200706003622-8a50f05110d2
	github.com/ahmetb/go-linq/v3 v3.2.0
	github.com/araddon/dateparse v0.0.0-20210429162001-6b43995a97de
	github.com/avast/retry-go v3.0.0+incompatible
	github.com/brianvoe/gofakeit/v6 v6.28.0
	github.com/caarlos0/env/v8 v8.0.0
	github.com/doug-martin/goqu/v9 v9.27.0
	github.com/elastic/go-elasticsearch/v8 v8.19.7
	github.com/glebarez/sqlite v1.11.0
	github.com/go-playground/validator v9.31.0+incompatible
	github.com/go-resty/resty/v2 v2.17.2
	github.com/go-testfixtures/testfixtures/v3 v3.19.0
	github.com/goccy/go-json v0.12.0
	github.com/goccy/go-reflect v1.2.0
	github.com/gofrs/uuid v4.4.0+incompatible
	github.com/golang-migrate/migrate/v4 v4.20.1
	github.com/grpc-ecosystem/go-grpc-middleware v1.4.0
	github.com/hibiken/asynq v0.26.0
	github.com/iancoleman/strcase v0.3.0
	github.com/jackc/pgconn v1.14.3
	github.com/jackc/pgx/v4 v4.18.3
	github.com/jmoiron/sqlx v1.4.0
	github.com/joho/godotenv v1.5.1
	github.com/kamva/mgm/v3 v3.5.0
	github.com/labstack/echo/v4 v4.16.6
	github.com/lib/pq v1.12.3
	github.com/mcuadros/go-defaults v1.2.0
	github.com/mehdihadeli/go-mediatr v1.4.0
	github.com/michaelklishin/rabbit-hole v1.5.0
	github.com/mitchellh/mapstructure v1.5.0
	github.com/moby/moby/api v1.57.0
	github.com/nolleh/caption_json_formatter v0.2.4
	github.com/onsi/ginkgo/v2 v2.34.0
	github.com/onsi/gomega v1.44.0
	github.com/orlangure/gnomock v0.32.0
	github.com/ory/dockertest/v3 v3.12.0
	github.com/phayes/freeport v0.0.0-20220201140144-74d24b5ae9f5
	github.com/pressly/goose/v3 v3.30.2
	github.com/prometheus/client_golang v1.25.0
	github.com/rabbitmq/amqp091-go v1.17.0
	github.com/redis/go-redis/extra/redisotel/v9 v9.24.0
	github.com/redis/go-redis/v9 v9.24.0
	github.com/samber/lo v1.53.0
	github.com/satori/go.uuid v1.2.0
	github.com/sirupsen/logrus v1.10.2
	github.com/spf13/viper v1.21.0
	github.com/stretchr/testify v1.12.1
	github.com/testcontainers/testcontainers-go v0.45.0
	github.com/testcontainers/testcontainers-go/modules/postgres v0.45.0
	github.com/ulule/limiter/v3 v3.11.2
	github.com/uptrace/bun v1.3.0
	github.com/uptrace/bun/dialect/pgdialect v1.3.0
	github.com/uptrace/bun/driver/pgdriver v1.3.0
	github.com/uptrace/opentelemetry-go-extra/otellogrus v0.4.0
	github.com/uptrace/opentelemetry-go-extra/otelzap v0.4.0
	go.mongodb.org/mongo-driver v1.17.10
	go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo v0.72.0
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.72.0
	go.opentelemetry.io/contrib/instrumentation/host v0.72.0
	go.opentelemetry.io/contrib/propagators/ot v1.47.0
	go.opentelemetry.io/otel v1.47.0
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v1.47.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.47.0
	go.opentelemetry.io/otel/exporters/prometheus v0.69.0
	go.opentelemetry.io/otel/exporters/stdout/stdoutmetric v1.47.0
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.47.0
	go.opentelemetry.io/otel/exporters/zipkin v1.47.0
	go.opentelemetry.io/otel/metric v1.47.0
	go.opentelemetry.io/otel/sdk v1.47.0
	go.opentelemetry.io/otel/sdk/metric v1.47.0
	go.opentelemetry.io/otel/trace v1.47.0
	go.uber.org/fx v1.24.0
	go.uber.org/zap v1.28.0
	google.golang.org/grpc v1.85.0
	gorm.io/driver/postgres v1.6.3
	gorm.io/gorm v1.31.2
	gorm.io/plugin/opentelemetry v0.1.16
)

require (
	dario.cat/mergo v1.0.2 // indirect
	filippo.io/edwards25519 v1.2.0 // indirect
	github.com/Azure/go-ansiterm v0.0.0-20261013165250-51075c6ea6bc // indirect
	github.com/ClickHouse/ch-go v0.76.0 // indirect
	github.com/ClickHouse/clickhouse-go/v2 v2.51.0 // indirect
	github.com/Masterminds/semver/v3 v3.5.0 // indirect
	github.com/Microsoft/go-winio v0.6.4 // indirect
	github.com/Nvveen/Gotty v0.0.0-20120604004816-cd527374f1e5 // indirect
	github.com/andybalholm/brotli v1.3.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/cenkalti/backoff/v5 v5.0.3 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/containerd/continuity v0.5.0 // indirect
	github.com/containerd/errdefs v1.0.0 // indirect
	github.com/containerd/errdefs/pkg v0.3.0 // indirect
	github.com/containerd/log v0.2.0 // indirect
	github.com/containerd/platforms v1.0.0-rc.5 // indirect
	github.com/cpuguy83/dockercfg v0.3.2 // indirect
	github.com/distribution/reference v0.6.0 // indirect
	github.com/docker/cli v29.10.1+incompatible // indirect
	github.com/docker/docker v28.5.2+incompatible // indirect
	github.com/docker/go-connections v0.8.2 // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/dustin/go-humanize v1.1.0 // indirect
	github.com/ebitengine/purego v0.11.3 // indirect
	github.com/elastic/elastic-transport-go/v8 v8.11.0 // indirect
	github.com/fatih/color v1.19.0 // indirect
	github.com/felixge/httpsnoop v1.1.0 // indirect
	github.com/fsnotify/fsnotify v1.10.1 // indirect
	github.com/glebarez/go-sqlite v1.23.0 // indirect
	github.com/go-faster/city v1.0.1 // indirect
	github.com/go-faster/errors v0.8.0 // indirect
	github.com/go-logr/logr v1.4.4 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-ole/go-ole v1.3.0 // indirect
	github.com/go-playground/locales v0.14.2 // indirect
	github.com/go-playground/universal-translator v0.18.2 // indirect
	github.com/go-sql-driver/mysql v1.10.2 // indirect
	github.com/go-task/slim-sprig/v3 v3.0.0 // indirect
	github.com/go-viper/mapstructure/v2 v2.5.0 // indirect
	github.com/goccy/go-yaml v1.19.2 // indirect
	github.com/golang/mock v1.6.0 // indirect
	github.com/golang/snappy v1.0.0 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/google/pprof v0.0.0-20261118010117-084b33a85f63 // indirect
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.34.0 // indirect
	github.com/hashicorp/go-version v1.9.0 // indirect
	github.com/hokaccha/go-prettyjson v0.0.0-20211117102719-0474bc63780f // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.3 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgtype v1.14.4 // indirect
	github.com/jackc/pgx/v5 v5.13.0 // indirect
	github.com/jackc/puddle v1.3.0 // indirect
	github.com/jackc/puddle/v2 v2.2.4 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/klauspost/compress v1.20.5 // indirect
	github.com/labstack/gommon v0.5.0 // indirect
	github.com/lann/builder v0.0.0-20180802200727-47ae307949d0 // indirect
	github.com/lann/ps v0.0.0-20150810152359-62de8c46ede0 // indirect
	github.com/leodido/go-urn v1.5.0 // indirect
	github.com/lufia/plan9stats v0.0.0-20261101121723-be84dd9509a0 // indirect
	github.com/magiconair/properties v1.18.13 // indirect
	github.com/mattn/go-colorable v0.1.16 // indirect
	github.com/mattn/go-isatty v0.0.24 // indirect
	github.com/mfridman/interpolate v0.2.0 // indirect
	github.com/moby/docker-image-spec v1.3.1 // indirect
	github.com/moby/go-archive v0.3.3 // indirect
	github.com/moby/moby/client v0.6.3 // indirect
	github.com/moby/patternmatcher v0.6.1 // indirect
	github.com/moby/sys/sequential v0.7.0 // indirect
	github.com/moby/sys/user v0.5.1 // indirect
	github.com/moby/sys/userns v0.2.1 // indirect
	github.com/moby/term v0.5.2 // indirect
	github.com/montanaflynn/stats v0.14.0 // indirect
	github.com/morikuni/aec v1.1.0 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/ncruces/go-strftime v1.1.0 // indirect
	github.com/nxadm/tail v1.4.8 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.1.1 // indirect
	github.com/opencontainers/runc v1.4.3 // indirect
	github.com/openzipkin/zipkin-go v0.4.3 // indirect
	github.com/pandatix/go-abnf v0.5.5 // indirect
	github.com/paulmach/orb v0.14.0 // indirect
	github.com/pelletier/go-toml/v2 v2.4.3 // indirect
	github.com/pierrec/lz4/v4 v4.1.34 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/power-devops/perfstat v0.0.0-20260916203055-22a1a467d9f0 // indirect
	github.com/prometheus/client_model v0.6.3 // indirect
	github.com/prometheus/common v0.74.0 // indirect
	github.com/prometheus/otlptranslator v1.0.0 // indirect
	github.com/prometheus/procfs v0.23.0 // indirect
	github.com/puzpuzpuz/xsync/v3 v3.5.1 // indirect
	github.com/redis/go-redis/extra/rediscmd/v9 v9.24.0 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	github.com/robfig/cron/v3 v3.0.1 // indirect
	github.com/sagikazarmark/locafero v0.12.0 // indirect
	github.com/santhosh-tekuri/jsonschema/v6 v6.0.3 // indirect
	github.com/segmentio/asm v1.2.1 // indirect
	github.com/shirou/gopsutil/v4 v4.26.10 // indirect
	github.com/shopspring/decimal v1.6.0 // indirect
	github.com/spf13/afero v1.15.0 // indirect
	github.com/spf13/cast v1.10.0 // indirect
	github.com/spf13/pflag v1.1.0 // indirect
	github.com/streadway/amqp v1.1.0 // indirect
	github.com/stretchr/objx v0.8.0 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	github.com/tklauser/go-sysconf v0.4.0 // indirect
	github.com/tklauser/numcpus v0.12.0 // indirect
	github.com/tmthrgd/go-hex v0.0.0-20190904060850-447a3041c3bc // indirect
	github.com/uptrace/opentelemetry-go-extra/otelutil v0.4.0 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	github.com/vmihailenco/msgpack/v5 v5.4.1 // indirect
	github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect
	github.com/xdg-go/scram v1.4.0 // indirect
	github.com/xdg-go/stringprep v1.0.5 // indirect
	github.com/youmark/pkcs8 v0.0.0-20240726163527-a2c0da244d78 // indirect
	github.com/yusufpapurcu/wmi v1.2.4 // indirect
	go.opentelemetry.io/auto/sdk v1.2.2 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.72.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.47.0 // indirect
	go.opentelemetry.io/otel/log v1.47.0 // indirect
	go.opentelemetry.io/proto/otlp v1.11.1 // indirect
	go.uber.org/atomic v1.12.0 // indirect
	go.uber.org/dig v1.19.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.yaml.in/yaml/v3 v3.0.6 // indirect
	golang.org/x/crypto v0.59.0 // indirect
	golang.org/x/mod v0.43.0 // indirect
	golang.org/x/net v0.63.0 // indirect
	golang.org/x/sync v0.25.0 // indirect
	golang.org/x/sys v0.50.0 // indirect
	golang.org/x/text v0.44.0 // indirect
	golang.org/x/time v0.16.0 // indirect
	golang.org/x/tools v0.53.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20261123191846-1e707ab9fbc2 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20261123191846-1e707ab9fbc2 // indirect
	google.golang.org/protobuf v1.36.12 // indirect
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	gorm.io/driver/clickhouse v0.7.0 // indirect
	gorm.io/driver/mysql v1.6.0 // indirect
	mellium.im/sasl v0.3.2 // indirect
	modernc.org/libc v1.82.13 // indirect
	modernc.org/mathutil v1.7.1 // indirect
	modernc.org/memory v1.13.1 // indirect
	modernc.org/sqlite v1.63.0 // indirect
)

// github.com/ory/dockertest/v3/docker/pkg/homedir still imports
// github.com/opencontainers/runc/libcontainer/user, which was removed in runc v1.3.0.
replace github.com/opencontainers/runc => github.com/opencontainers/runc v1.2.3
