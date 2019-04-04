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
    stdout *os.File
    stderr *os.File
}

// GetName returns CLI name.
func (c *CLI) GetName() string {
    return c.name
}

// GetDesc returns CLI description.
func (c *CLI) GetDesc() string {
    return c.desc
}

// GetAuthor returns CLI author.
func (c *CLI) GetAuthor() string {
    return c.author
}

// AttachCmd attaches instance of CLICmd to CLI.
func (c *CLI) AttachCmd(cmd *CLICmd) {
    n := reflect.ValueOf(cmd).Elem().FieldByName("name").String()
    if c.cmds == nil { 
        c.cmds = make(map[string]*CLICmd); 
    }
    c.cmds[n] = cmd
}

// GetCmd returns instance of CLICmd of command k.
func (c *CLI) GetCmd(k string) *CLICmd {
    return c.cmds[k]
}

// GetCmds returns list of commands.
func (c *CLI) GetCmds() ([]reflect.Value) {
    return reflect.ValueOf(c.cmds).MapKeys()
}

// PrintUsage prints usage information like available commands to CLI stdout.
func (c *CLI) PrintUsage() {
    fmt.Fprintf(c.stdout, c.name + " by " + c.author + "\n" + c.desc + "\n\n")
    fmt.Fprintf(c.stdout, "Available commands:\n")
    for _, i := range c.GetCmds() {
        cmd := c.GetCmd(i.String())
        fmt.Fprintf(c.stdout, path.Base(os.Args[0]) + " " + cmd.GetName() + cmd.GetFlagsUsage() + "\n")
    }
}

// AddCmd creates a new command with name n, description d and handler of f.
// It creates instance of CLICmd, attaches it to CLI and returns it.
func (c *CLI) AddCmd(n string, d string, f func(cli *CLI) int) *CLICmd {
    cmd := NewCLICmd(n, d, f)
    c.AttachCmd(cmd)
    return cmd
}

// getFlagSetPtrs creates flagset instance, parses flags and returns list of
// pointers to results of parsing the flags.
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

// parseFlags iterates over flags and validates them.
// In case of error it prints out to CLI stderr.
func (c *CLI) parseFlags(cmd *CLICmd) int {
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
                fmt.Fprintf(c.stderr, "ERROR: Flag --"+flagName+" is missing!\n\n")
                c.PrintUsage()
                return 1
            }
            if flag.GetNFlags() & CLIFlagTypePathFile > 0 && flag.GetNFlags() & CLIFlagMustExist > 0 {
                filePath := *(flagSetPtrs[flagName]).(*string)
                if _, err := os.Stat(filePath); os.IsNotExist(err) {
                    fmt.Fprintf(c.stderr, "ERROR: File "+filePath+" from --"+flagName+" does not exist!\n\n")
                    c.PrintUsage()
                    return 1
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
    return 0
}

// Run parses the arguments, validates them and executes command handler. In
// case of invalid arguments, error is printed to stderr and 1 is returned.
// Return value behaves like exit code.
func (c *CLI) Run(stdout *os.File, stderr *os.File) int {
    c.stdout = stdout; c.stderr = stderr
    if len(os.Args[1:]) < 1 {
        c.PrintUsage()
        return 1
    }
    for _, cmd := range c.GetCmds() {
        if cmd.String() == os.Args[1] {
            exitCode := c.parseFlags(c.GetCmd(cmd.String()))
            if exitCode > 0 {
                return exitCode
            }
            return c.GetCmd(cmd.String()).Run(c)
        }
    }
    c.PrintUsage()
    return 1
}

// Returns value of flag.
func (c *CLI) Flag(n string) string {
    return c.parsedFlags[n]
}

// NewCLI creates new instance of CLI with name n, description d and author a 
// and returns it.
func NewCLI(n string, d string, a string) *CLI {
    c := &CLI{ name: n, desc: d, author: a }
    return c
}

