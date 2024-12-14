package config

// Env is an environment string
type Env string

const (
	// Local is a local development environment
	Local Env = "local"
	// Demo is the demo environment
	Demo Env = "demo"
	// Stage
	Stage Env = "stg"
	// Preprod
	Preprod Env = "preprod"
	// QA is the QA environment
	QA Env = "qa"
	// Production is the production environment
	Production Env = "prod"
	// Test
	Test Env = "test"
)

func (e Env) String() string {
	switch e {
	case Demo:
		return "demo"
	case Stage:
		return "stg"
	case QA:
		return "qa"
	case Production:
		return "prod"
	case Preprod:
		return "preprod"
	case Test:
		return "test"
	default:
		return "local"
	}
}

// UnmarshalConfig unmarshals a config value string to the associated interface
// type
func (e *Env) UnmarshalConfig(value string) {
	switch value {
	case "demo":
		*e = Demo
	case "qa":
		*e = QA
	case "production", "prod":
		*e = Production
	case "preprod":
		*e = Preprod
	case "stage", "stg", "staging":
		*e = Stage
	case "test":
		*e = Test
	default:
		*e = Local
	}
}

// init registers the defaults
func init() {
	RegisterDefault("dd_env", Local)
	RegisterDefault("system_port", 9102)
	RegisterDefault("export_runtime_metrics", true)
	RegisterDefault("trace_sample_rate", 0.1)
	RegisterDefault("trace_max_batch_count", 256)
	RegisterDefault("bind_address", "0.0.0.0")
}

// Base is the base configuration for all services
type Base struct {
	Env                          Env     `config:"dd_env"`
	ServiceName                  string  `config:"dd_service"`
	ServiceVersion               string  `config:"dd_version"`
	TraceAddress                 string  `config:"trace_address"`
	MetricsAddress               string  `config:"metrics_address"`
	LogLevel                     string  `config:"log_level"`
	SystemPort                   int     `config:"system_port"`
	TraceMaxBatchCount           int     `config:"trace_max_batch_count"`
	TraceSampleRate              float64 `config:"trace_sample_rate"`
	ProfilerBlockProfileRate     int     `config:"profiler_block_profile_rate"`
	ProfilerMutexProfileFraction int     `config:"profiler_mutex_profile_fraction"`
	ExportRuntimeMetrics         bool    `config:"export_runtime_metrics"`
	ProfilerEnabled              bool    `config:"profiler_enabled"`
	DisableLogSampling           bool    `config:"disable_log_sampling"`
	DisableStackTraces           bool    `config:"disable_stack_traces"`
	BindAddress                  string  `config:"bind_address"`
}
