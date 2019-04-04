package cli

const (
    // Flag is required.
    CLIFlagRequired     = 1
    // Flag is string.
    CLIFlagTypeString   = 8
    // Flag is path to a file.
    CLIFlagTypePathFile = 16
    // Flag is a boolean type.
    CLIFlagTypeBool     = 32
    // If flag is path to a file, this flag tells that the file must exist.
    CLIFlagMustExist    = 128
)

// CLIFlag represends flag and has a name, description and configuration which
// is represented as integer value, eg. it can be value of
// CLIFlagRequired|CLIFlagTypePathFile|CLIFlagMustExist.
type CLIFlag struct {
    name string
    desc string
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
    f := &CLIFlag{ name: n, desc: d, nflags: nf }
    return f
}

