package main

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
func (c *CLIFlag) GetName() string {
    return c.name
}
func (c *CLIFlag) GetDesc() string {
    return c.desc
}
func (c *CLIFlag) GetNFlags() int32 {
    return c.nflags
}
func NewCLIFlag(n string, d string, nf int32) *CLIFlag {
    f := &CLIFlag{ name: n, desc: d, nflags: nf }
    return f
}

