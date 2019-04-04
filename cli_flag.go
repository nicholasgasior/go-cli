package cli

const (
    CLIFlagRequired     = 1
    CLIFlagTypeString   = 8
    CLIFlagTypePathFile = 16
    CLIFlagTypeBool     = 32
    CLIFlagMustExist    = 128
)
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

