package goda

// LoggerFunc is the signature for a custom logger.
type LoggerFunc func(config *Config, node *Node, level LogLevel, format string, args ...interface{}) int

// CloneNodeFunc is the signature for a custom node cloning callback.
type CloneNodeFunc func(oldNode *Node, owner *Node, childIndex int) *Node

// DirtiedFunc is called when a node becomes dirty.
type DirtiedFunc func(node *Node)

// MeasureFunc is the signature for a custom measure function.
type MeasureFunc func(node *Node, width float32, widthMode MeasureMode, height float32, heightMode MeasureMode) Size

// BaselineFunc is the signature for a custom baseline function.
type BaselineFunc func(node *Node, width, height float32) float32

// Config holds global layout configuration.
type Config struct {
	logger               LoggerFunc
	cloneNodeCallback    CloneNodeFunc
	useWebDefaults       bool
	pointScaleFactor     float32
	context              interface{}
	version              uint32
	errata               Errata
	experimentalFeatures uint64
}

func NewConfig(logger LoggerFunc) *Config {
	return &Config{
		logger:           logger,
		errata:           ErrataMinSizeUndefinedInsteadOfAuto,
		pointScaleFactor: 1.0,
	}
}

var defaultConfig *Config

func GetDefaultConfig() *Config {
	if defaultConfig == nil {
		defaultConfig = NewConfig(DefaultLogger)
	}
	return defaultConfig
}

func (c *Config) SetUseWebDefaults(use bool) { c.useWebDefaults = use }
func (c *Config) UseWebDefaults() bool       { return c.useWebDefaults }

func (c *Config) SetPointScaleFactor(factor float32) {
	if c.pointScaleFactor != factor {
		c.pointScaleFactor = factor
		c.version++
	}
}
func (c *Config) GetPointScaleFactor() float32 { return c.pointScaleFactor }

func (c *Config) SetErrata(errata Errata) {
	if c.errata != errata {
		c.errata = errata
		c.version++
	}
}
func (c *Config) GetErrata() Errata { return c.errata }

func (c *Config) AddErrata(errata Errata) {
	if !c.HasErrata(errata) {
		c.errata |= errata
		c.version++
	}
}
func (c *Config) RemoveErrata(errata Errata) {
	if c.HasErrata(errata) {
		c.errata &^= errata
		c.version++
	}
}
func (c *Config) HasErrata(errata Errata) bool { return (c.errata & errata) != ErrataNone }

func (c *Config) SetExperimentalFeatureEnabled(feature ExperimentalFeature, enabled bool) {
	mask := uint64(1) << uint(feature)
	currently := (c.experimentalFeatures & mask) != 0
	if currently != enabled {
		if enabled {
			c.experimentalFeatures |= mask
		} else {
			c.experimentalFeatures &^= mask
		}
		c.version++
	}
}
func (c *Config) IsExperimentalFeatureEnabled(feature ExperimentalFeature) bool {
	return (c.experimentalFeatures & (uint64(1) << uint(feature))) != 0
}
func (c *Config) GetEnabledExperiments() uint64 { return c.experimentalFeatures }

func (c *Config) SetContext(ctx interface{}) { c.context = ctx }
func (c *Config) GetContext() interface{}    { return c.context }

func (c *Config) GetVersion() uint32 { return c.version }

func (c *Config) SetLogger(logger LoggerFunc) {
	if logger != nil {
		c.logger = logger
	} else {
		c.logger = DefaultLogger
	}
}

func (c *Config) Log(node *Node, level LogLevel, format string, args ...interface{}) {
	c.logger(c, node, level, format, args...)
}

func (c *Config) SetCloneNodeCallback(callback CloneNodeFunc) { c.cloneNodeCallback = callback }

func (c *Config) CloneNode(node *Node, owner *Node, childIndex int) *Node {
	if c.cloneNodeCallback != nil {
		clone := c.cloneNodeCallback(node, owner, childIndex)
		if clone != nil {
			return clone
		}
	}
	return node.Clone()
}

func configUpdateInvalidatesLayout(oldCfg, newCfg *Config) bool {
	return oldCfg.GetErrata() != newCfg.GetErrata() ||
		oldCfg.GetEnabledExperiments() != newCfg.GetEnabledExperiments() ||
		oldCfg.GetPointScaleFactor() != newCfg.GetPointScaleFactor() ||
		oldCfg.UseWebDefaults() != newCfg.UseWebDefaults()
}

// DefaultLogger is a no-op logger that suppresses error/fatal messages.
var DefaultLogger LoggerFunc = func(config *Config, node *Node, level LogLevel, format string, args ...interface{}) int {
	if level == LogLevelError || level == LogLevelFatal {
		return 0
	}
	return 0
}

func (c *Config) SetUseWebDefaultsBool(b bool) { c.useWebDefaults = b }
