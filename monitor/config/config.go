package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/Bnei-Baruch/mdb/monitor/models"
	"github.com/Bnei-Baruch/mdb/monitor/plugins/inputs"
	"github.com/Bnei-Baruch/mdb/monitor/plugins/outputs"
	"github.com/Bnei-Baruch/mdb/monitor/plugins/serializers"
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
)

// Config specifies the monitoring configurations
// such as connection string, user/password,
// as well as all the plugins that the user has
// specified
type Config struct {
	Tags          map[string]string
	InputFilters  []string
	OutputFilters []string
	Agent         *AgentConfig
	Inputs        []*models.RunningInput
	Outputs       []*models.RunningOutput
	Aggregators   []*models.RunningAggregator
	Processors    models.RunningProcessors
}

// NewConfig - Config constructor
func NewConfig() *Config {
	c := &Config{
		// Agent defaults:
		Agent: &AgentConfig{
			Interval:      10 * time.Second,
			RoundInterval: true,
			FlushInterval: 10 * time.Second,
		},

		Tags:          make(map[string]string),
		Inputs:        make([]*models.RunningInput, 0),
		Outputs:       make([]*models.RunningOutput, 0),
		Processors:    make([]*models.RunningProcessor, 0),
		InputFilters:  make([]string, 0),
		OutputFilters: make([]string, 0),
	}
	return c
}

// AgentConfig specifies the monitoring agent configurations such as collection and flushing intervals
type AgentConfig struct {
	// Interval at which to gather information
	Interval time.Duration

	// RoundInterval rounds collection interval to 'interval'.
	//     ie, if Interval=10s then always collect on :00, :10, :20, etc.
	RoundInterval bool

	// By default or when set to "0s", precision will be set to the same
	// timestamp order as the collection interval, with the maximum being 1s.
	//   ie, when interval = "10s", precision will be "1s"
	//       when interval = "250ms", precision will be "1ms"
	// Precision will NOT be used for service inputs. It is up to each individual
	// service input to set the timestamp at the appropriate precision.
	Precision time.Duration

	// CollectionJitter is used to jitter the collection by a random amount.
	// Each plugin will sleep for a random time within jitter before collecting.
	// This can be used to avoid many plugins querying things like sysfs at the
	// same time, which can have a measurable effect on the system.
	CollectionJitter time.Duration

	// FlushInterval is the Interval at which to flush data
	FlushInterval time.Duration

	// FlushJitter Jitters the flush interval by a random amount.
	// This is primarily to avoid large write spikes for users running a large
	// number of telegraf instances.
	// ie, a jitter of 5s and interval 10s means flushes will happen every 10-15s
	FlushJitter time.Duration

	// MetricBatchSize is the maximum number of metrics that is wrote to an
	// output plugin in one call.
	MetricBatchSize int

	// MetricBufferLimit is the max number of metrics that each output plugin
	// will cache. The buffer is cleared when a successful write occurs. When
	// full, the oldest metrics will be overwritten. This number should be a
	// multiple of MetricBatchSize. Due to current implementation, this could
	// not be less than 2 times MetricBatchSize.
	MetricBufferLimit int

	// FlushBufferWhenFull tells Telegraf to flush the metric buffer whenever
	// it fills up, regardless of FlushInterval. Setting this option to true
	// does _not_ deactivate FlushInterval.
	FlushBufferWhenFull bool

	// TODO(cam): Remove UTC and parameter, they are no longer
	// valid for the agent config. Leaving them here for now for backwards-
	// compatibility
	UTC bool `toml:"utc"`

	// Debug is the option for running in debug mode
	Debug bool

	// Logfile specifies the file to send logs to
	Logfile string

	// Quiet is the option for running in quiet mode
	Quiet        bool
	Hostname     string
	OmitHostname bool
}

// LoadConfig loads the given config file and applies it to c
func (c *Config) LoadConfig(path string) error {
	log.Infof("Load monitor configuration file from : %s", path)
	viperInstance := viper.New()
	viperInstance.SetConfigFile(path)
	viperInstance.AddConfigPath(".") // path to look for the config file in
	viperInstance.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viperInstance.AutomaticEnv()
	if err := viperInstance.ReadInConfig(); err != nil { // Find and read the config file
		log.Errorf("couldn't load config: %s", err) // Handle errors reading the config file
		return err
	}

	// Read and unmarshal agent details
	if err := viperInstance.UnmarshalKey("agent", &c.Agent); err != nil {
		log.Errorf("couldn't read agent config: %s", err)
		return err
	}

	inputsConfigEntries := viper.GetStringMap("inputs")
	log.Infof("Read and parse %d inputs...", len(inputsConfigEntries))
	for inputName, inputConfigEntry := range inputsConfigEntries {
		log.Infof("Input name: %s, and value: %v\n", inputName, inputConfigEntry)
		inputConfigEntryMap := inputConfigEntry.([]interface{})[0].(map[string]interface{})
		creator, ok := inputs.Inputs[inputName]
		if !ok {
			log.Errorf("Undefined but requested input: %s", inputName)
			return fmt.Errorf("Undefined but requested input: %s, defined inputs : %s", inputName, inputs.Inputs)
		}
		input := creator()
		err := input.TryParseConfigurations(inputConfigEntryMap)
		if err != nil {
			log.Errorf("Couldn't decode and parse input configurations: %s", err)
		}
		pluginConfig := &models.InputConfig{Name: inputName}
		rp := models.NewRunningInput(input, pluginConfig)
		c.Inputs = append(c.Inputs, rp)
		log.Infof("Input with description: %s created", input.Description())
	}

	outputsConfigEntries := viper.GetStringMap("outputs")
	log.Infof("Read and parse %d outputs...", len(outputsConfigEntries))
	for outputName, outputConfigEntry := range outputsConfigEntries {
		log.Infof("Output name: %s, and value: %v\n", outputName, outputConfigEntry)
		outputConfigEntryMap := outputConfigEntry.([]interface{})[0].(map[string]interface{})
		creator, ok := outputs.Outputs[outputName]
		if !ok {
			log.Errorf("Undefined but requested output: %s", outputName)
			return fmt.Errorf("Undefined but requested output: %s", outputName)
		}
		output := creator()
		// If the output has a SetSerializer function, then this means it can write
		// arbitrary types of output, so build the serializer and set it.
		switch t := output.(type) {
		case serializers.SerializerOutput:
			serializerConfig := &serializers.Config{}
			if value, ok := outputConfigEntryMap["data_format"]; ok {
				serializerConfig.DataFormat = value.(string)
			}

			// serializerConfig.Prefix = viper.GetString("prefix")
			// serializerConfig.InfluxMaxLineBytes = viper.GetInt("influx_max_line_bytes")
			// serializerConfig.InfluxSortFields = viper.GetBool("influx_sort_fields")
			// serializerConfig.InfluxUintSupport = viper.GetBool("influx_uint_support")
			// serializerConfig.Template = viper.GetString("template")
			serializer, err := serializers.NewSerializer(serializerConfig)
			if err != nil {
				log.Errorf("Couldn't read serializer config and/or create it: %s", err)
			}

			t.SetSerializer(serializer)
			log.Infof("Set and use %s serializer...", serializerConfig.DataFormat)
		}

		err := output.TryParseConfigurations(outputConfigEntryMap)
		if err != nil {
			log.Errorf("Couldn't decode and parse output configurations: %s", err)
		}

		pluginConfig := &models.OutputConfig{Name: outputName}
		rp := models.NewRunningOutput(outputName, output, pluginConfig, 1000, 10000)
		c.Outputs = append(c.Outputs, rp)
		log.Infof("Output with description: %s created", output.Description())
	}

	log.Infof("Loading monitor configuration from file : %s successfully completed...", path)
	return nil
}
