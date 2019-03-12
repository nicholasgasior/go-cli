package main

import (
    "reflect"
    "os"
    "fmt"
    "path"
)

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
    os.Exit(1)
}
func (c *CLI) Run() {
    if len(os.Args[1:]) < 1 {
        c.PrintUsage()
    }
    for _, cmd := range c.GetCmds() {
        if cmd.String() == os.Args[1] {
            _ = c.GetCmd(cmd.String()).GetFlags()
            // invalid flags should print an error and usage
            // for flags a flagset should be created
            os.Exit(c.GetCmd(cmd.String()).Run())
        }
    }
    c.PrintUsage()
}
func NewCLI(n string, d string, a string) *CLI {
    c := &CLI{ name: n, desc: d, author: a }
    return c
}

