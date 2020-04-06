package cli

import (
	"fmt"
	"os"
	"path"
	"reflect"
	"sort"
	"text/tabwriter"
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

// PrintHelp prints command usage information to stdout file.
func (c *CLICmd) PrintHelp(cli *CLI) {
	fmt.Fprintf(cli.GetStdout(), "\nUsage:  "+path.Base(os.Args[0])+" "+c.GetName()+" [OPTIONS]\n\n")
	fmt.Fprintf(cli.GetStdout(), c.GetDesc()+"\n")

	w := new(tabwriter.Writer)
	w.Init(cli.GetStdout(), 8, 8, 0, '\t', 0)

	var s [2]string
	i := 1
	for _, n := range c.GetSortedFlags() {
		flag := c.GetFlag(n)
		if flag.IsRequired() {
			i = 0
		} else {
			i = 1
		}
		s[i] += flag.GetHelpLine()
	}

	if s[0] != "" {
		fmt.Fprintf(w, "\nRequired flags: \n")
		fmt.Fprintf(w, s[0])
		w.Flush()
	}
	if s[1] != "" {
		fmt.Fprintf(w, "\nOptional flags: \n")
		fmt.Fprintf(w, s[1])
		w.Flush()
	}

}

// AttachFlag attaches instance of CLIFlag to CLICmd.
func (c *CLICmd) AttachFlag(flag *CLIFlag) {
	n := reflect.ValueOf(flag).Elem().FieldByName("name").String()
	if c.flags == nil {
		c.flags = make(map[string]*CLIFlag)
	}
	c.flags[n] = flag
}

// AddFlag adds a flag to a command.
// It creates CLIFlag instance and attaches it.
func (c *CLICmd) AddFlag(n string, a string, hv string, d string, nf int32) {
	flg := NewCLIFlag(n, a, hv, d, nf)
	c.AttachFlag(flg)
}

// GetFlag returns instance of CLIFlag of flag k.
func (c *CLICmd) GetFlag(k string) *CLIFlag {
	return c.flags[k]
}

// GetSortedFlags returns sorted list of flag names.
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
