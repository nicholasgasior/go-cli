package cli

import (
	"reflect"
    "sort"
)

// CLICmd represent a command which has a name (used in args when calling app),
// description, a handler and flags attached to it.
type CLICmd struct {
	name                string
	desc                string
	flags               map[string]*CLIFlag
	args                map[string]*CLIFlag
	argsOrder           map[int]string
	nextArgIdx        int
	excludedGlobalFlags map[string]bool
	excludedGlobalArgs  map[string]bool
	handler             func(c *CLI) int
}

// GetName returns CLICmd name.
func (c *CLICmd) GetName() string {
	return c.name
}

// GetDesc returns CLICmd description.
func (c *CLICmd) GetDesc() string {
	return c.desc
}

// AttachFlag attaches instance of CLIFlag to CLICmd.
func (c *CLICmd) AttachFlag(flag *CLIFlag) {
	n := reflect.ValueOf(flag).Elem().FieldByName("name").String()
	if c.flags == nil {
		c.flags = make(map[string]*CLIFlag)
	}
	c.flags[n] = flag
}

// AttachArg attaches instance of CLIFlag to CLICmd as argument.
func (c *CLICmd) AttachArg(flag *CLIFlag) {
	n := reflect.ValueOf(flag).Elem().FieldByName("name").String()
	if c.args == nil {
		c.args = make(map[string]*CLIFlag)
	}
    if c.argsOrder == nil {
        c.argsOrder = make(map[int]string)
    }
	c.args[n] = flag
	c.argsOrder[c.nextArgIdx] = n
	c.nextArgIdx++
}

// AddFlag adds a flag with name n, description d and config of nf.
// It creates CLIFlag instance and attaches it.
func (c *CLICmd) AddFlag(n string, d string, dv string, nf int32) {
	f := NewCLIFlag(n, d, dv, nf)
	c.AttachFlag(f)
}

// AddArg adds a flag as argument and attaches it.
func (c *CLICmd) AddArg(n string, d string, nf int32) {
	f := NewCLIFlag(n, d, "", nf)
	c.AttachArg(f)
}

func (c *CLICmd) ExcludeGlobalFlag(f string) {
	if c.excludedGlobalFlags == nil {
		c.excludedGlobalFlags = make(map[string]bool)
	}
	c.excludedGlobalFlags[f] = true
}

func (c *CLICmd) ExcludeGlobalArg(f string) {
	if c.excludedGlobalArgs == nil {
		c.excludedGlobalArgs = make(map[string]bool)
	}
	c.excludedGlobalArgs[f] = true
}

// GetFlag returns instance of CLIFlag of flag k.
func (c *CLICmd) GetFlag(k string) *CLIFlag {
	return c.flags[k]
}

// Getarg returns instance of CLIFlag of arg k.
func (c *CLICmd) GetArgByName(k string) *CLIFlag {
	return c.args[k]
}

func (c *CLICmd) GetArgByIdx(i int) *CLIFlag {
	return c.args[c.argsOrder[i]]
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

// GetArgs returns list of arg names.
func (c *CLICmd) GetArgsOrder() map[int]string {
	return c.argsOrder
}

func (c *CLICmd) GetExcludedGlobalFlags() map[string]bool {
    return c.excludedGlobalFlags
}

func (c *CLICmd) GetExcludedGlobalArgs() map[string]bool {
    return c.excludedGlobalArgs
}

// GetFlagsUsage eturns flags usage.
func (c *CLICmd) GetFlagsUsage() string {
	s := ""
    for i:=0; i<2; i++ {
        for _, n := range c.GetSortedFlags() {
            f := c.GetFlag(n)
            if (f.IsRequired() && i == 0) || (!f.IsRequired() && i == 1) {
                s += " " + f.GetUsage("--", true)
            }
        }
    }
	return s
}

func (c *CLICmd) GetArgsUsage(r bool) string {
    s := ""
    for j:=0; j<len(reflect.ValueOf(c.args).MapKeys()); j++ {
        f := c.GetArgByIdx(j)
        if (r && f.IsRequired()) || (!r && !f.IsRequired()) {
            s += " " + f.GetUsage("", false)
        }
    }
    return s
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
