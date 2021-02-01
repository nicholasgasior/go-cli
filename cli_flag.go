package cli

import (
	"errors"
	"os"
	"regexp"
)

const (
	// Required sets flag to be required.
	Required = 1
	// TypeString sets flag to be string.
	TypeString = 8
	// TypePathFile sets flag to be path to a file.
	TypePathFile = 16
	// TypeBool sets flag to be boolean.
	TypeBool = 32
	// TypeInt sets flag to be integer.
	TypeInt = 64
	// TypeFloat sets flag to be float.
	TypeFloat = 128
	// TypeAlphanumeric sets flag to be alphanumeric.
	TypeAlphanumeric = 256
	// MustExist sets flag must point to an existing file (required TypePathFile to be added as well).
	MustExist = 512
	// AllowMany allows flag to have more than one value separated by comma by default.
	// For example: AllowMany with TypeInt allows values like: 123 or 123,455,666 or 12,222
	// AllowMany works only with TypeInt, TypeFloat and TypeAlphanumeric.
	AllowMany = 1024
	// ManySeparatorColon works with AllowMany and sets colon to be the value separator, instead of colon.
	ManySeparatorColon = 2048
	// ManySeparatorSemiColon works with AllowMany and sets semi-colon to be the value separator.
	ManySeparatorSemiColon = 4096
	// AllowDots can be used only with TypeAlphanumeric and additionally allows flag to have dots.
	AllowDots = 8192
	// AllowUnderscore can be used only with TypeAlphanumeric and additionally allows flag to have underscore chars.
	AllowUnderscore = 16384
	// AllowHyphen can be used only with TypeAlphanumeric and additionally allows flag to have hyphen chars.
	AllowHyphen         = 32768
	TypeEmail           = 65536
	TypeFQDN            = 131072
	TypePathDir         = 262144
	TypePathRegularFile = 524288
)

// CLIFlag represends flag. It has a name, alias, description, value that is
// shown when printing help and configuration which is an integer value. It can
// be for example Required|TypePathFile|MustExist.
type CLIFlag struct {
	name      string
	alias     string
	helpValue string
	desc      string
	nflags    int32
}

// GetName returns flag name.
func (c *CLIFlag) GetName() string {
	return c.name
}

// GetAlias returns flag alias (its short version, eg. -h).
func (c *CLIFlag) GetAlias() string {
	return c.alias
}

// GetHelpValue returns value that is shown when help is printed.
func (c *CLIFlag) GetHelpValue() string {
	return c.helpValue
}

// GetDesc return flag description.
func (c *CLIFlag) GetDesc() string {
	return c.desc
}

// GetNFlags return flag configuration.
func (c *CLIFlag) GetNFlags() int32 {
	return c.nflags
}

// GetHelpLine returns flag usage info that is used when printing help.
func (c *CLIFlag) GetHelpLine() string {
	s := " "
	if c.GetAlias() == "" {
		s += " \t"
	} else {
		s += " -" + c.GetAlias() + ",\t"
	}
	s += " --" + c.GetName() + " " + c.GetHelpValue() + " \t" + c.GetDesc() + "\n"
	return s
}

// IsRequired returns true when flag is required.
func (c *CLIFlag) IsRequired() bool {
	return c.nflags&Required > 0
}

// IsRequireValue returns true when flag requires a value (only bool one returns false).
func (c *CLIFlag) IsRequireValue() bool {
	return c.nflags&TypeString > 0 || c.nflags&TypePathFile > 0 || c.nflags&TypeInt > 0 || c.nflags&TypeFloat > 0 || c.nflags&TypeAlphanumeric > 0
}

// IsTypeBool returns true when flag is of bool type.
func (c *CLIFlag) IsTypeBool() bool {
	return c.nflags&TypeBool > 0
}

// IsTypeInt returns true when flag is of int type.
func (c *CLIFlag) IsTypeInt() bool {
	return c.nflags&TypeInt > 0
}

// IsTypeFloat returns true when flag is of float type.
func (c *CLIFlag) IsTypeFloat() bool {
	return c.nflags&TypeFloat > 0
}

// IsTypeAlphanumeric returns true when flag is of alphanumeric type.
func (c *CLIFlag) IsTypeAlphanumeric() bool {
	return c.nflags&TypeAlphanumeric > 0
}

// IsTypeString returns true when flag is of string type.
func (c *CLIFlag) IsTypeString() bool {
	return c.nflags&TypeString > 0
}

// IsTypePathFile returns true when flag should be path to a file.
func (c *CLIFlag) IsTypePathFile() bool {
	return c.nflags&TypePathFile > 0
}

// IsTypePathRegularFile returns true when flag should be path to a regular file.
func (c *CLIFlag) IsTypePathRegularFile() bool {
	return c.nflags&TypePathRegularFile > 0
}

// IsTypePathDir returns true when flags should be path to a directory.
func (c *CLIFlag) IsTypePathDir() bool {
	return c.nflags&TypePathDir > 0
}

// ValidateValue takes value coming from --NAME and -ALIAS and validates it.
func (c *CLIFlag) ValidateValue(isArg bool, nz string, az string) error {
	// both alias and name cannot be set
	if nz != "" && az != "" {
		return errors.New("Both -" + c.GetAlias() + " and --" + c.GetName() + " passed")
	}

	label := "Flag"
	if isArg {
		label = "Argument"
	}

	nlabel := c.GetName()
	if isArg {
		nlabel = c.GetHelpValue()
	}

	// empty
	if c.IsRequired() && (nz == "" && az == "") {
		if c.IsTypeString() || c.IsTypePathFile() || c.IsTypePathDir() || c.IsTypeInt() || c.IsTypeFloat() || c.IsTypeAlphanumeric() {
			return errors.New(label + " " + nlabel + " is missing")
		}
	}
	// string does not need any additional checks apart from the above one
	if c.IsTypeString() {
		return nil
	}
	v := az
	if nz != "" {
		v = nz
	}
	if c.IsRequired() || v != "" {
		// if flag is a file and have to exist
		if c.IsTypePathFile() {
			if _, err := os.Stat(v); os.IsNotExist(err) {
				return errors.New("File " + v + " from " + nlabel + " does not exist")
			}
			return nil
		}
		// if flag is a regular file and have to exist
		if c.IsTypePathRegularFile() {
			fileInfo, err := os.Stat(v)
			if os.IsNotExist(err) {
				return errors.New("File " + v + " from " + nlabel + " does not exist")
			}
			if !fileInfo.Mode().IsRegular() {
				return errors.New("Path " + v + " from " + nlabel + " is not a regular file")
			}
			return nil
		}
		// if flag is a directory and have to exist
		if c.IsTypePathDir() {
			fileInfo, err := os.Stat(v)
			if os.IsNotExist(err) {
				return errors.New("Directory " + v + " from " + nlabel + " does not exist")
			}
			if !fileInfo.IsDir() {
				return errors.New("Path " + v + " from " + nlabel + " is not a directory")
			}
			return nil
		}
		// int, float, alphanumeric - single or many, separated by various chars
		var reType string
		var reValue string
		// set regexp part just for the type (eg. int, float, anum)
		if c.IsTypeInt() {
			reType = "[0-9]+"
		} else if c.IsTypeFloat() {
			reType = "[0-9]{1,16}\\.[0-9]{1,16}"
		} else if c.IsTypeAlphanumeric() {
			// alphanumeric + additional characters
			if c.nflags&AllowHyphen > 0 && c.nflags&AllowUnderscore > 0 && c.nflags&AllowDots > 0 {
				reType = "[0-9a-zA-Z_\\.\\-]+"
			} else if c.nflags&AllowUnderscore > 0 && c.nflags&AllowDots > 0 {
				reType = "[0-9a-zA-Z_\\.]+"
			} else if c.nflags&AllowUnderscore > 0 && c.nflags&AllowHyphen > 0 {
				reType = "[0-9a-zA-Z_\\-]+"
			} else if c.nflags&AllowDots > 0 && c.nflags&AllowHyphen > 0 {
				reType = "[0-9a-zA-Z\\.\\-]+"
			} else if c.nflags&AllowUnderscore > 0 {
				reType = "[0-9a-zA-Z_]+"
			} else if c.nflags&AllowDots > 0 {
				reType = "[0-9a-zA-Z\\.]+"
			} else {
				reType = "[0-9a-zA-Z]+"
			}
		}
		// create the final regexp depending on if single or many values are allowed
		if c.nflags&AllowMany > 0 {
			var d string
			if c.nflags&ManySeparatorColon > 0 {
				d = ":"
			} else if c.nflags&ManySeparatorSemiColon > 0 {
				d = ";"
			} else {
				d = ","
			}
			reValue = "^" + reType + "(" + d + reType + ")*$"
		} else {
			reValue = "^" + reType + "$"
		}
		m, err := regexp.MatchString(reValue, v)
		if err != nil || !m {
			return errors.New(label + " " + nlabel + " has invalid value")
		}
	}
	return nil
}

// NewCLIFlag creates instance of CLIFlag and returns it.
func NewCLIFlag(n string, a string, hv string, d string, nf int32) *CLIFlag {
	f := &CLIFlag{name: n, alias: a, helpValue: hv, desc: d, nflags: nf}
	return f
}
