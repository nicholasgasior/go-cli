package cli

import (
	"flag"
	"fmt"
	"os"
	"path"
	"reflect"
	"regexp"
    "sort"
)

// CLI is main CLI application definition. It has a name, description, author
// (which are used only when printing usage syntax), commands and pointers
// to File instances to which standard output or errors are printed (named
// respectively stdout and stderr).
type CLI struct {
	name               string
	desc               string
	author             string
	cmds               map[string]*CLICmd
	globalArgs         map[string]*CLIFlag
	globalArgsOrder    map[int]string
	nextGlobalArgIdx int
	globalFlags        map[string]*CLIFlag
	parsedFlags        map[string]string
	parsedArgs         map[int]string
	stdout             *os.File
	stderr             *os.File
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
		c.cmds = make(map[string]*CLICmd)
	}
	c.cmds[n] = cmd
}

// AttachGlobalFlag attaches instance of CLIFlag to CLI.
func (c *CLI) AttachGlobalFlag(flag *CLIFlag) {
	n := reflect.ValueOf(flag).Elem().FieldByName("name").String()
	if c.globalFlags == nil {
		c.globalFlags = make(map[string]*CLIFlag)
	}
	c.globalFlags[n] = flag
}

// AttachGlobalArg attache instance of CLIFlag to CLI as Argument.
func (c *CLI) AttachGlobalArg(flag *CLIFlag) {
	n := reflect.ValueOf(flag).Elem().FieldByName("name").String()
	if c.globalArgs == nil {
		c.globalArgs = make(map[string]*CLIFlag)
	}
    if c.globalArgsOrder == nil {
        c.globalArgsOrder = make(map[int]string)
    }
	c.globalArgs[n] = flag
	c.globalArgsOrder[c.nextGlobalArgIdx] = n
	c.nextGlobalArgIdx++
}

// GetCmd returns instance of CLICmd of command k.
func (c *CLI) GetCmd(k string) *CLICmd {
	return c.cmds[k]
}

// GetGlobalFlag returns instance of CLIFlag of flag k.
func (c *CLI) GetGlobalFlag(k string) *CLIFlag {
	return c.globalFlags[k]
}

// GetGlobalArg returns instance of CLIFlag of flag k.
func (c *CLI) GetGlobalArgByName(k string) *CLIFlag {
	return c.globalArgs[k]
}

func (c *CLI) GetGlobalArgByIdx(i int) *CLIFlag {
	return c.globalArgs[c.globalArgsOrder[i]]
}

// GetStdout returns stdout property.
func (c *CLI) GetStdout() *os.File {
	return c.stdout
}

// GetStderr returns stderr property.
func (c *CLI) GetStderr() *os.File {
	return c.stderr
}

func (c *CLI) GetSortedCmds() []string {
    cmds := reflect.ValueOf(c.cmds).MapKeys()
    scmds := make([]string, len(cmds))
    for i, cmd := range cmds {
        scmds[i] = cmd.String()
    }
    sort.Strings(scmds)
    return scmds
}

func (c *CLI) GetSortedGlobalFlags() []string {
    fs := reflect.ValueOf(c.globalFlags).MapKeys()
    sfs := make([]string, len(fs))
    for i, f := range fs {
        sfs[i] = f.String()
    }
    sort.Strings(sfs)
    return sfs
}

// GetGlobalFlags returns list of flags.
func (c *CLI) GetGlobalFlags() []reflect.Value {
	return reflect.ValueOf(c.globalFlags).MapKeys()
}

func (c *CLI) GetGlobalArgsOrder() map[int]string {
	return c.globalArgsOrder
}

func (c *CLI) GetGlobalFlagsUsage(e map[string]bool) string {
    s := ""
    for i:=0; i<2; i++ {
        for _, n := range c.GetSortedGlobalFlags() {
            if e[n] != true {
                f := c.GetGlobalFlag(n)
                if (f.IsRequired() && i == 0) || (!f.IsRequired() && i == 1) {
                    s += " " + f.GetUsage("--", true)
                }
            }
        }
    }
    return s
}

func (c *CLI) GetCmdArgsUsage(cmd *CLICmd) string {
    s := ""
    e := cmd.GetExcludedGlobalArgs()
    for i:=0; i<2; i++ {
        for j:=0; j<len(reflect.ValueOf(c.globalArgs).MapKeys()); j++ {
            n := c.globalArgsOrder[j]
            if e[n] != true {
                f := c.GetGlobalArgByIdx(j)
                if (f.IsRequired() && i == 0) || (!f.IsRequired() && i == 1) {
                    s += " " + f.GetUsage("", false)
                }
            }
        }
        if i == 0 {
            s += cmd.GetArgsUsage(true)
        } else {
            s += cmd.GetArgsUsage(false)
        }
    }
    return s
}

// PrintUsage prints usage information like available commands to CLI stdout.
func (c *CLI) PrintUsage() {
	fmt.Fprintf(c.stdout, c.name+" by "+c.author+"\n"+c.desc+"\n\n")
	fmt.Fprintf(c.stdout, "Available commands:\n")
    for _, i := range c.GetSortedCmds() {
		cmd := c.GetCmd(i)
        fmt.Fprintf(c.stdout, path.Base(os.Args[0])+" "+cmd.GetName()+c.GetGlobalFlagsUsage(cmd.GetExcludedGlobalFlags())+cmd.GetFlagsUsage()+c.GetCmdArgsUsage(cmd)+"\n")
	}
}

// AddCmd creates a new command with name n, description d and handler of f.
// It creates instance of CLICmd, attaches it to CLI and returns it.
func (c *CLI) AddCmd(n string, d string, f func(cli *CLI) int) *CLICmd {
	cmd := NewCLICmd(n, d, f)
	c.AttachCmd(cmd)
	return cmd
}

// AddGlobalFlag creates a new global flag and attaches it to CLI.
func (c *CLI) AddGlobalFlag(n string, d string, dv string, nf int32) {
	f := NewCLIFlag(n, d, dv, nf)
	c.AttachGlobalFlag(f)
}

// AddGlobalArg creates a new global arg and attaches it to CLI.
func (c *CLI) AddGlobalArg(n string, d string, nf int32) {
	f := NewCLIFlag(n, d, "", nf)
	c.AttachGlobalArg(f)
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
		if flag.GetNFlags()&CLIFlagTypeString > 0 ||
			flag.GetNFlags()&CLIFlagTypePathFile > 0 ||
			flag.GetNFlags()&CLIFlagTypeInt > 0 ||
			flag.GetNFlags()&CLIFlagTypeFloat > 0 ||
			flag.GetNFlags()&CLIFlagTypeAlphanumeric > 0 {
			flagSetPtrs[flagName] = flagSet.String(flagName, "", flag.GetDesc())
		} else if flag.GetNFlags()&CLIFlagTypeBool > 0 {
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
		var flagValue string
		if flag.GetNFlags()&CLIFlagTypeBool == 0 {
			flagValue = *(flagSetPtrs[flagName]).(*string)
		}
		if flag.GetNFlags()&CLIFlagRequired > 0 && (flag.GetNFlags()&CLIFlagTypeString > 0 || flag.GetNFlags()&CLIFlagTypePathFile > 0) {
			if flagValue == "" {
				fmt.Fprintf(c.stderr, "ERROR: Flag --"+flagName+" is missing!\n\n")
				c.PrintUsage()
				return 1
			}
			if flag.GetNFlags()&CLIFlagTypePathFile > 0 && flag.GetNFlags()&CLIFlagMustExist > 0 {
				filePath := flagValue
				if _, err := os.Stat(filePath); os.IsNotExist(err) {
					fmt.Fprintf(c.stderr, "ERROR: File "+filePath+" from --"+flagName+" does not exist!\n\n")
					c.PrintUsage()
					return 1
				}
			}
		}
		if (flag.GetNFlags()&CLIFlagRequired > 0 || flagValue != "") && flag.GetNFlags()&CLIFlagTypeInt > 0 {
			valuePattern := "[0-9]+"
			var reToMatch string
			if flag.GetNFlags()&CLIFlagAllowMany > 0 {
				if flag.GetNFlags()&CLIFlagManySeparatorColon > 0 {
					reToMatch = "^" + valuePattern + "(:" + valuePattern + ")*$"
				} else if flag.GetNFlags()&CLIFlagManySeparatorSemiColon > 0 {
					reToMatch = "^" + valuePattern + "(;" + valuePattern + ")*$"
				} else {
					reToMatch = "^" + valuePattern + "(," + valuePattern + ")*$"
				}
			} else {
				reToMatch = "^" + valuePattern + "$"
			}
			matched, err := regexp.MatchString(reToMatch, flagValue)
			if err != nil || !matched {
				fmt.Fprintf(c.stderr, "ERROR: Flag --"+flagName+" is not a valid integer!\n\n")
				c.PrintUsage()
				return 1
			}
		}
		if (flag.GetNFlags()&CLIFlagRequired > 0 || flagValue != "") && flag.GetNFlags()&CLIFlagTypeFloat > 0 {
			valuePattern := "[0-9]{1,16}\\.[0-9]{1,16}"
			var reToMatch string
			if flag.GetNFlags()&CLIFlagAllowMany > 0 {
				if flag.GetNFlags()&CLIFlagManySeparatorColon > 0 {
					reToMatch = "^" + valuePattern + "(:" + valuePattern + ")*$"
				} else if flag.GetNFlags()&CLIFlagManySeparatorSemiColon > 0 {
					reToMatch = "^" + valuePattern + "(;" + valuePattern + ")*$"
				} else {
					reToMatch = "^" + valuePattern + "(," + valuePattern + ")*$"
				}
			} else {
				reToMatch = "^" + valuePattern + "$"
			}
			matched, err := regexp.MatchString(reToMatch, flagValue)
			if err != nil || !matched {
				fmt.Fprintf(c.stderr, "ERROR: Flag --"+flagName+" is not a valid float!\n\n")
				c.PrintUsage()
				return 1
			}
		}
		if (flag.GetNFlags()&CLIFlagRequired > 0 || flagValue != "") && flag.GetNFlags()&CLIFlagTypeAlphanumeric > 0 {
			var valuePattern string
			if flag.GetNFlags()&CLIFlagAllowUnderscore > 0 && flag.GetNFlags()&CLIFlagAllowDots > 0 {
				valuePattern = "[0-9a-zA-Z_\\.]+"
			} else if flag.GetNFlags()&CLIFlagAllowUnderscore > 0 {
				valuePattern = "[0-9a-zA-Z_]+"
			} else if flag.GetNFlags()&CLIFlagAllowDots > 0 {
				valuePattern = "[0-9a-zA-Z\\.]+"
			} else {
				valuePattern = "[0-9a-zA-Z]+"
			}
			var reToMatch string
			if flag.GetNFlags()&CLIFlagAllowMany > 0 {
				if flag.GetNFlags()&CLIFlagManySeparatorColon > 0 {
					reToMatch = "^" + valuePattern + "(:" + valuePattern + ")*$"
				} else if flag.GetNFlags()&CLIFlagManySeparatorSemiColon > 0 {
					reToMatch = "^" + valuePattern + "(;" + valuePattern + ")*$"
				} else {
					reToMatch = "^" + valuePattern + "(," + valuePattern + ")*$"
				}
			} else {
				reToMatch = "^" + valuePattern + "$"
			}
			matched, err := regexp.MatchString(reToMatch, flagValue)
			if err != nil || !matched {
				fmt.Fprintf(c.stderr, "ERROR: Flag --"+flagName+" is not a valid alphanumeric value!\n\n")
				c.PrintUsage()
				return 1
			}
		}
		if flag.GetNFlags()&CLIFlagTypeString > 0 || flag.GetNFlags()&CLIFlagTypePathFile > 0 {
			c.parsedFlags[flagName] = flagValue
		}
		if flag.GetNFlags()&CLIFlagTypeBool > 0 {
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
	c.stdout = stdout
	c.stderr = stderr

	if len(os.Args[1:]) < 1 {
		c.PrintUsage()
		return 1
	}

    return 0

	/*for _, cmd := range c.GetCmds() {
		if cmd.String() == os.Args[1] {
			exitCode := c.parseFlags(c.GetCmd(cmd.String()))
			if exitCode > 0 {
				return exitCode
			}
			return c.GetCmd(cmd.String()).Run(c)
		}
	}
	c.PrintUsage()
	return 1*/
}

// Flag returns value of flag.
func (c *CLI) Flag(n string) string {
	return c.parsedFlags[n]
}

// Arg returns value of argument.
func (c *CLI) Arg(i int) string {
	return c.parsedArgs[i]
}

// NewCLI creates new instance of CLI with name n, description d and author a
// and returns it.
func NewCLI(n string, d string, a string) *CLI {
	c := &CLI{name: n, desc: d, author: a}
	return c
}
