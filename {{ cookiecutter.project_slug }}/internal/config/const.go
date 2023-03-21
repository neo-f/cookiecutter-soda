package config

const (
	URL_HEALTH         = "/healthz"
	URL_MONITOR        = "/sys/monitor"
	URL_VERSION        = "/sys/version"
	URL_OPENAPI        = "/openapi.json"
	URL_REDOC          = "/redoc"
	URL_SWAGGER        = "/swagger"
	URL_GATEWAY_CONFIG = "/v1/global-config/api-gateway"
	URL_METRICS        = "/metrics"
	URL_RAPIDOC        = "/rapidoc"
	URL_ELEMENTS       = "/"
)

const (
	LogTagURL           = "url"
	LogTagMethod        = "method"
	LogTagHeaders       = "headers"
	LogTagData          = "data"
	LogTagAuthorization = "authorization"
	LogTagBID           = "bid"
	LogTagStaffID       = "staff_id"
	LogTagTraceID       = "trace_id"
)
