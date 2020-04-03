package cli

import (
	"reflect"
	"sort"
)

// CLICmd represent a command which has a name (used in args when calling app),
// description, a handler and flags attached to it.
type CLICmd struct {
	name    string
	desc    string
	flags   map[string]*CLIFlag
	handler func(c *CLI) int
}

// GetName returns CLICmd name.
func (c *CLICmd) GetName() string {
	return c.name
}

// GetDesc returns CLICmd description.
func (c *CLICmd) GetDesc() string {
	return c.desc
}

// GetFlagsUsage eturns flags usage.
func (c *CLICmd) GetFlagsUsage() string {
	s := ""
	for i := 0; i < 2; i++ {
		for _, n := range c.GetSortedFlags() {
			flag := c.GetFlag(n)
			if (i == 0 && flag.IsRequired()) || (i == 1 && !flag.IsRequired()) {
				s += " " + flag.GetUsage("--", true)
			}
		}
	}
	return s
}

// AttachFlag attaches instance of CLIFlag to CLICmd.
func (c *CLICmd) AttachFlag(flag *CLIFlag) {
	n := reflect.ValueOf(flag).Elem().FieldByName("name").String()
	if c.flags == nil {
		c.flags = make(map[string]*CLIFlag)
	}
	c.flags[n] = flag
}

// AddFlag adds a flag with name n, description d and config of nf.
// It creates CLIFlag instance and attaches it.
func (c *CLICmd) AddFlag(n string, d string, nf int32) {
	flg := NewCLIFlag(n, d, nf)
	c.AttachFlag(flg)
}

// GetFlag returns instance of CLIFlag of flag k.
func (c *CLICmd) GetFlag(k string) *CLIFlag {
	return c.flags[k]
}

func (c *CLICmd) GetSortedFlags() []string {
	fs := reflect.ValueOf(c.flags).MapKeys()
	sfs := make([]string, len(fs))
	for i, f := range fs {
		sfs[i] = f.String()
	}
	sort.Strings(sfs)
	return sfs
}

// GetFlags returns list of flag names.
func (c *CLICmd) GetFlags() []reflect.Value {
	return reflect.ValueOf(c.flags).MapKeys()
}

// Run calls command handler.
func (c *CLICmd) Run(cli *CLI) int {
	return c.handler(cli)
}

// NewCLICmd creates CLICmd instance with name n, description d and handler f
// and returns it.
func NewCLICmd(n string, d string, f func(cli *CLI) int) *CLICmd {
	c := &CLICmd{name: n, desc: d, handler: f}
	return c
}
