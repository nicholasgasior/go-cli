package main

import (
    "reflect"
    "os"
    "fmt"
    "path"
)


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


type CLICmd struct {
    name string
    desc string
    flags map[string]*CLIFlag
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
func (c *CLICmd) AddFlag(flag *CLIFlag) {
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
func NewCLICmd(n string, d string) *CLICmd {
    c := &CLICmd{ name: n, desc: d }
    return c
}


type CLI struct {
    name string
    desc string
    author string
    cmds map[string]*CLICmd
}
func (c *CLI) GetName() string {
    return c.name
}
func (c *CLI) GetDesc() string {
    return c.desc
}
func (c *CLI) GetAuthor() string {
    return c.author
}
func (c *CLI) AddCmd(cmd *CLICmd) {
    n := reflect.ValueOf(cmd).Elem().FieldByName("name").String()
    if c.cmds == nil { 
        c.cmds = make(map[string]*CLICmd); 
    }
    c.cmds[n] = cmd
}
func (c *CLI) GetCmd(k string) *CLICmd {
    return c.cmds[k]
}
func (c *CLI) GetCmds() ([]reflect.Value) {
    return reflect.ValueOf(c.cmds).MapKeys()
}
func (c *CLI) PrintUsage() {
    fmt.Fprintf(os.Stdout, c.name + " by " + c.author + "\n" + c.desc + "\n\n")
    fmt.Fprintf(os.Stdout, "Available commands:\n")
    for _, i := range c.GetCmds() {
        cmd := c.GetCmd(i.String())
        fmt.Fprintf(os.Stdout, path.Base(os.Args[0]) + " " + cmd.GetName() + cmd.GetFlagsUsage() + "\n")
    }
}
func (c *CLI) Run() {
    if len(os.Args[1:]) < 1 {
        c.PrintUsage()
    }
}
func NewCLI(n string, d string, a string) *CLI {
    c := &CLI{ name: n, desc: d, author: a }
    return c
}


func main() {
    cli := NewCLI("Example API", "Simple API application", "Nicholas Gasior <nameless@laatu.org>")

    cmdServe := NewCLICmd("serve", "Runs http server")
    cmdCreateSchema := NewCLICmd("create_schema", "Creates schema in the database")

    flagConfig := NewCLIFlag("config", "Config file", CLIFlagTypePathFile | CLIFlagMustExist | CLIFlagRequired)
    flagDrop := NewCLIFlag("drop", "Drop tables", CLIFlagTypeBool)

    cmdServe.AddFlag(flagConfig)
    cmdCreateSchema.AddFlag(flagConfig)
    cmdCreateSchema.AddFlag(flagDrop)

    cli.AddCmd(cmdServe)
    cli.AddCmd(cmdCreateSchema)
    cli.Run()
}

