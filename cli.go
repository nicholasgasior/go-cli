package cli

import (
    "reflect"
    "os"
    "fmt"
    "path"
    "flag"
)

type CLI struct {
    name string
    desc string
    author string
    cmds map[string]*CLICmd
    parsedFlags map[string]string
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
func (c *CLI) AttachCmd(cmd *CLICmd) {
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
func (c *CLI) AddCmd(n string, d string, f func(cli *CLI) int) *CLICmd {
    cmd := NewCLICmd(n, d, f)
    c.AttachCmd(cmd)
    return cmd
}
func (c *CLI) getFlagSetPtrs(cmd *CLICmd) map[string]interface{} {
    flagSet := flag.NewFlagSet("flagset", flag.ExitOnError)
    flagSetPtrs := make(map[string]interface{})
    flags := cmd.GetFlags()
    for _, i := range flags {
        flagName := i.String()
        flag := cmd.GetFlag(flagName)
        if flag.GetNFlags() & CLIFlagTypeString > 0 || flag.GetNFlags() & CLIFlagTypePathFile > 0 {
            flagSetPtrs[flagName] = flagSet.String(flagName, "", flag.GetDesc())
        } else if flag.GetNFlags() & CLIFlagTypeBool > 0 {
            flagSetPtrs[flagName] = flagSet.Bool(flagName, false, flag.GetDesc())
        }
    }
    flagSet.Parse(os.Args[2:])
    return flagSetPtrs
}
func (c *CLI) parseFlags(cmd *CLICmd) {
    if c.parsedFlags == nil {
        c.parsedFlags = make(map[string]string)
    }

    flags := cmd.GetFlags()
    flagSetPtrs := c.getFlagSetPtrs(cmd)
    for _, i := range flags {
        flagName := i.String()
        flag := cmd.GetFlag(flagName)
        if flag.GetNFlags() & CLIFlagRequired > 0 && (flag.GetNFlags() & CLIFlagTypeString > 0 || flag.GetNFlags() & CLIFlagTypePathFile > 0) {
            if *(flagSetPtrs[flagName]).(*string) == "" {
                fmt.Fprintf(os.Stderr, "ERROR: Flag --"+flagName+" is missing!\n\n")
                c.PrintUsage()
            }
            if flag.GetNFlags() & CLIFlagTypePathFile > 0 && flag.GetNFlags() & CLIFlagMustExist > 0 {
                filePath := *(flagSetPtrs[flagName]).(*string)
                if _, err := os.Stat(filePath); os.IsNotExist(err) {
                    fmt.Fprintf(os.Stderr, "ERROR: File "+filePath+" from --"+flagName+" does not exist!\n\n")
                    c.PrintUsage()
                }
            }
        }
        if flag.GetNFlags() & CLIFlagTypeString > 0 || flag.GetNFlags() & CLIFlagTypePathFile > 0 {
            c.parsedFlags[flagName] = *(flagSetPtrs[flagName]).(*string)
        }
        if flag.GetNFlags() & CLIFlagTypeBool > 0 {
            if *(flagSetPtrs[flagName]).(*bool) == true {
                c.parsedFlags[flagName] = "true"
            } else {
                c.parsedFlags[flagName] = "false"
            }
        }
    }
}
func (c *CLI) Run() {
    if len(os.Args[1:]) < 1 {
        c.PrintUsage()
    }
    for _, cmd := range c.GetCmds() {
        if cmd.String() == os.Args[1] {
            c.parseFlags(c.GetCmd(cmd.String()))

            os.Exit(c.GetCmd(cmd.String()).Run(c))
        }
    }
    c.PrintUsage()
}
func (c *CLI) Flag(n string) string {
    return c.parsedFlags[n]
}
func NewCLI(n string, d string, a string) *CLI {
    c := &CLI{ name: n, desc: d, author: a }
    return c
}

