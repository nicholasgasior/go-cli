package cli

const (
	// CLIFlagRequired sets flag to be required.
	CLIFlagRequired = 1
	// CLIFlagTypeString sets flag to be string.
	CLIFlagTypeString = 8
	// CLIFlagTypePathFile sets flag to be path to a file.
	CLIFlagTypePathFile = 16
	// CLIFlagTypeBool sets flag to be boolean.
	CLIFlagTypeBool = 32
	// CLIFlagTypeInt sets flag to be integer.
	CLIFlagTypeInt = 64
	// CLIFlagTypeFloat sets flag to be float.
	CLIFlagTypeFloat = 128
	// CLIFlagTypeAlphanumeric sets flag to be alphanumeric.
	CLIFlagTypeAlphanumeric = 256
	// CLIFlagMustExist sets flag must point to an existing file (required CLIFlagTypePathFile to be added as well).
	CLIFlagMustExist = 512
	// CLIFlagAllowMany allows flag to have more than one value separated by comma by default.
	// For example: CLIFlagAllowMany with CLIFlagTypeInt allows values like: 123 or 123,455,666 or 12,222
	// CLIFlagAllowMany works only with CLIFlagTypeInt, CLIFlagTypeFloat and CLIFlagTypeAlphanumeric.
	CLIFlagAllowMany = 1024
	// CLIFlagManySeparatorColon works with CLIFlagAllowMany and sets colon to be the value separator, instead of colon.
	CLIFlagManySeparatorColon = 2048
	// CLIFlagManySeparatorSemiColon works with CLIFlagAllowMany and sets semi-colon to be the value separator.
	CLIFlagManySeparatorSemiColon = 4096
	// CLIFlagAllowDots can be used only with CLIFlagTypeAlphanumeric and additionally allows flag to have dots.
	CLIFlagAllowDots = 8192
	// CLIFlagAllowUnderscore can be used only with CLIFlagTypeAlphanumeric and additionally allows flag to have underscore chars.
	CLIFlagAllowUnderscore = 16384
)

// CLIFlag represends flag and has a name, description and configuration which
// is represented as integer value, eg. it can be value of
// CLIFlagRequired|CLIFlagTypePathFile|CLIFlagMustExist.
type CLIFlag struct {
	name   string
	desc   string
	nflags int32
}

// GetName returns flag name.
func (c *CLIFlag) GetName() string {
	return c.name
}

// GetDesc return flag description.
func (c *CLIFlag) GetDesc() string {
	return c.desc
}

// GetNFlags return flag configuration.
func (c *CLIFlag) GetNFlags() int32 {
	return c.nflags
}

// NewCLIFlag creates instance of CLIFlag with name n, description d and config
// of nf and returns it.
func NewCLIFlag(n string, d string, nf int32) *CLIFlag {
	f := &CLIFlag{name: n, desc: d, nflags: nf}
	return f
}
