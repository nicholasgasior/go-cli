package cli

import (
    "reflect"
)

type CLICmd struct {
    name string
    desc string
    flags map[string]*CLIFlag
    handler func(c *CLI) int
}
func (c *CLICmd) GetName() string {
    return c.name
}
func (c *CLICmd) GetDesc() string {
    return c.desc
}
func (c *CLICmd) GetFlagsUsage() string {
    s := ""
    for _, i := range c.GetFlags() {
        flag := c.GetFlag(i.String())
        s += " "
        if flag.GetNFlags() & CLIFlagRequired == 0 {
            s += "["
        }
        s += "--" + flag.GetName()
        if flag.GetNFlags() & CLIFlagTypeString > 0 {
            s += "=STRING"
        } else if flag.GetNFlags() & CLIFlagTypePathFile > 0 {
            s += "=FILEPATH"
        }
        if flag.GetNFlags() & CLIFlagRequired == 0 {
            s += "]"
        }
    }
    return s
}
func (c *CLICmd) AttachFlag(flag *CLIFlag) {
    n := reflect.ValueOf(flag).Elem().FieldByName("name").String()
    if c.flags == nil {
        c.flags = make(map[string]*CLIFlag);
    }
    c.flags[n] = flag
}
func (c *CLICmd) GetFlag(k string) *CLIFlag {
    return c.flags[k]
}
func (c *CLICmd) GetFlags() ([]reflect.Value) {
    return reflect.ValueOf(c.flags).MapKeys()
}
func (c *CLICmd) Run(cli *CLI) int {
    return c.handler(cli)
}
func NewCLICmd(n string, d string, f func(cli *CLI) int) *CLICmd {
    c := &CLICmd{ name: n, desc: d, handler: f }
    return c
}

