module github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice

go 1.20

// https://go.dev/doc/tutorial/call-module-code
replace github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg => ../../pkg/

require (
	emperror.dev/errors v0.8.1
	github.com/EventStore/EventStore-Client-Go v1.0.2
	github.com/brianvoe/gofakeit/v6 v6.21.0
	github.com/elastic/go-elasticsearch/v8 v8.8.0
	github.com/gavv/httpexpect/v2 v2.15.0
	github.com/go-playground/validator v9.31.0+incompatible
	github.com/goccy/go-json v0.10.2
	github.com/labstack/echo/v4 v4.10.2
	github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg v0.0.0-00010101000000-000000000000
	github.com/mehdihadeli/go-mediatr v1.1.10
	github.com/michaelklishin/rabbit-hole v1.5.0
	github.com/satori/go.uuid v1.2.0
	github.com/spf13/cobra v1.1.3
	github.com/spf13/viper v1.12.0
	github.com/stretchr/testify v1.8.3
	github.com/swaggo/echo-swagger v1.3.3
	github.com/swaggo/swag v1.16.1
	go.mongodb.org/mongo-driver v1.11.6
	go.opentelemetry.io/otel v1.16.0
	go.opentelemetry.io/otel/metric v1.16.0
	go.opentelemetry.io/otel/trace v1.16.0
	go.uber.org/fx v1.20.0
	google.golang.org/grpc v1.55.0
	google.golang.org/protobuf v1.30.0
)

require (
	github.com/Azure/go-ansiterm v0.0.0-20210617225240-d185dfc1b5a1 // indirect
	github.com/KyleBanks/depth v1.2.1 // indirect
	github.com/Microsoft/go-winio v0.6.0 // indirect
	github.com/PuerkitoBio/purell v1.1.1 // indirect
	github.com/PuerkitoBio/urlesc v0.0.0-20170810143723-de5bf2ad4578 // indirect
	github.com/TylerBrock/colorjson v0.0.0-20200706003622-8a50f05110d2 // indirect
	github.com/ahmetb/go-linq/v3 v3.2.0 // indirect
	github.com/ajg/form v1.5.1 // indirect
	github.com/andybalholm/brotli v1.0.5 // indirect
	github.com/araddon/dateparse v0.0.0-20210429162001-6b43995a97de // indirect
	github.com/avast/retry-go v3.0.0+incompatible // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/caarlos0/env/v8 v8.0.0 // indirect
	github.com/cenkalti/backoff/v4 v4.2.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/containerd/containerd v1.6.19 // indirect
	github.com/cpuguy83/dockercfg v0.3.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/docker/distribution v2.8.1+incompatible // indirect
	github.com/docker/docker v23.0.5+incompatible // indirect
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/elastic/elastic-transport-go/v8 v8.2.0 // indirect
	github.com/fatih/color v1.13.0 // indirect
	github.com/fatih/structs v1.1.0 // indirect
	github.com/fsnotify/fsnotify v1.5.4 // indirect
	github.com/go-logr/logr v1.2.4 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/jsonreference v0.19.6 // indirect
	github.com/go-openapi/spec v0.20.4 // indirect
	github.com/go-openapi/swag v0.21.1 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-redis/redis/v8 v8.11.5 // indirect
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/goccy/go-reflect v1.2.0 // indirect
	github.com/gofrs/uuid v4.4.0+incompatible // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	github.com/golang-migrate/migrate/v4 v4.15.2 // indirect
	github.com/golang/mock v1.6.0 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.4.0 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/iancoleman/strcase v0.2.0 // indirect
	github.com/imkira/go-interpol v1.1.0 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.3.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/kamva/mgm/v3 v3.5.0 // indirect
	github.com/klauspost/compress v1.16.5 // indirect
	github.com/labstack/gommon v0.4.0 // indirect
	github.com/leodido/go-urn v1.2.4 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/mcuadros/go-defaults v1.2.0 // indirect
	github.com/mitchellh/go-wordwrap v1.0.1 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/moby/patternmatcher v0.5.0 // indirect
	github.com/moby/sys/sequential v0.5.0 // indirect
	github.com/moby/term v0.0.0-20221128092401-c43b287e0e0f // indirect
	github.com/montanaflynn/stats v0.0.0-20171201202039-1bf9dbcd8cbe // indirect
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/nolleh/caption_json_formatter v0.2.2 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.1.0-rc2 // indirect
	github.com/opencontainers/runc v1.1.5 // indirect
	github.com/openzipkin/zipkin-go v0.4.1 // indirect
	github.com/pelletier/go-toml v1.9.5 // indirect
	github.com/pelletier/go-toml/v2 v2.0.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_golang v1.15.1 // indirect
	github.com/prometheus/client_model v0.4.0 // indirect
	github.com/prometheus/common v0.42.0 // indirect
	github.com/prometheus/procfs v0.9.0 // indirect
	github.com/rabbitmq/amqp091-go v1.8.1 // indirect
	github.com/sanity-io/litter v1.5.5 // indirect
	github.com/sergi/go-diff v1.0.0 // indirect
	github.com/sirupsen/logrus v1.9.2 // indirect
	github.com/solsw/collate v0.1.0 // indirect
	github.com/solsw/generichelper v0.4.0 // indirect
	github.com/solsw/go2linq/v2 v2.5.3 // indirect
	github.com/spf13/afero v1.8.2 // indirect
	github.com/spf13/cast v1.5.0 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/streadway/amqp v1.1.0 // indirect
	github.com/stretchr/objx v0.5.0 // indirect
	github.com/subosito/gotenv v1.3.0 // indirect
	github.com/swaggo/files v0.0.0-20220610200504-28940afbdbfe // indirect
	github.com/testcontainers/testcontainers-go v0.19.0 // indirect
	github.com/tmthrgd/go-hex v0.0.0-20190904060850-447a3041c3bc // indirect
	github.com/uptrace/bun v1.1.14 // indirect
	github.com/uptrace/bun/driver/pgdriver v1.1.14 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.34.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	github.com/vmihailenco/msgpack/v5 v5.3.5 // indirect
	github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.1 // indirect
	github.com/xdg-go/stringprep v1.0.3 // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20190905194746-02993c407bfb // indirect
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415 // indirect
	github.com/xeipuuv/gojsonschema v1.2.0 // indirect
	github.com/yalp/jsonpath v0.0.0-20180802001716-5cc68e5049a0 // indirect
	github.com/youmark/pkcs8 v0.0.0-20181117223130-1be2e3e5546d // indirect
	github.com/yudai/gojsondiff v1.0.0 // indirect
	github.com/yudai/golcs v0.0.0-20170316035057-ecda9a501e82 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.42.0 // indirect
	go.opentelemetry.io/contrib/propagators/ot v1.17.0 // indirect
	go.opentelemetry.io/otel/exporters/jaeger v1.16.0 // indirect
	go.opentelemetry.io/otel/exporters/prometheus v0.39.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.16.0 // indirect
	go.opentelemetry.io/otel/exporters/zipkin v1.16.0 // indirect
	go.opentelemetry.io/otel/sdk v1.16.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v0.39.0 // indirect
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/dig v1.17.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.24.0 // indirect
	golang.org/x/crypto v0.9.0 // indirect
	golang.org/x/exp v0.0.0-20230315142452-642cacee5cc0 // indirect
	golang.org/x/mod v0.9.0 // indirect
	golang.org/x/net v0.10.0 // indirect
	golang.org/x/sync v0.1.0 // indirect
	golang.org/x/sys v0.8.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	golang.org/x/time v0.3.0 // indirect
	golang.org/x/tools v0.7.0 // indirect
	google.golang.org/genproto v0.0.0-20230306155012-7f2fa6fef1f4 // indirect
	gopkg.in/ini.v1 v1.66.4 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	gorm.io/driver/postgres v1.5.2 // indirect
	gorm.io/gorm v1.25.1 // indirect
	mellium.im/sasl v0.3.1 // indirect
	moul.io/http2curl/v2 v2.3.0 // indirect
)
