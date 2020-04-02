package cli

import (
	"os"
	"testing"
)

func h(c *CLI) int {
	return 0
}

func cli() *CLI {
	c := NewCLI("Example CLI", "Silly app", "Author <a@example.com>")
	return c
}

func assertExitCode(t *testing.T, cli *CLI, a []string, c int) {
	os.Args = a
	f, _ := os.Open("/dev/null")
	defer f.Close()
	got := cli.Run(f, f)
	want := c
	if got != want {
		t.Errorf("got %d want %d\n", got, want)
	}
}

func TestRequiredFlags(t *testing.T) {
	c1 := cli()
	c1.AddGlobalFlag("flag1", "Required flag", "", CLIFlagTypeString|CLIFlagRequired)
	c1.AddGlobalFlag("flag2", "Optional flag", "", CLIFlagTypeInt)
	c1.AddGlobalFlag("flag3", "Required flag", "ignored", CLIFlagTypeInt|CLIFlagRequired)
	_ = c1.AddCmd("cmd", "Sample command", h)

	c2 := cli()
	c2.AddGlobalFlag("flag1", "Required flag", "", CLIFlagTypeString|CLIFlagRequired)
	c2.AddGlobalFlag("flag2", "Optional flag", "", CLIFlagTypeInt)
	c2.AddGlobalFlag("flag3", "Required flag", "ignored", CLIFlagTypeInt|CLIFlagRequired)
	cmd21 := c2.AddCmd("cmd", "Sample command", h)
	cmd21.AddFlag("flag4", "Required flag", "", CLIFlagTypeFloat|CLIFlagRequired)
	cmd21.AddFlag("flag5", "Optional flag", "", CLIFlagTypeString)
	cmd21.AddFlag("flag6", "Required flag", "ignored", CLIFlagTypeInt|CLIFlagRequired)
	cmd21.ExcludeGlobalFlag("flag3")

	c3 := cli()
	cmd31 := c3.AddCmd("cmd", "Sample command", h)
	cmd31.AddFlag("flag4", "Required flag", "", CLIFlagTypeFloat|CLIFlagRequired)
	cmd31.AddFlag("flag5", "Optional flag", "", CLIFlagTypeString)

	t.Run("exit with code 1 when required global flag is missing", func(t *testing.T) {
		assertExitCode(t, c1, []string{"test", "cmd"}, 1)
		assertExitCode(t, c1, []string{"test", "cmd", "--flag1=string"}, 1)
		assertExitCode(t, c1, []string{"test", "cmd", "--flag1=string", "--flag2=123"}, 1)

		assertExitCode(t, c2, []string{"test", "cmd", "--flag4=string", "--flag6=123"}, 1)
	})

	t.Run("exit with code 1 when required command flag is missing", func(t *testing.T) {
		assertExitCode(t, c3, []string{"test", "cmd"}, 1)
		assertExitCode(t, c3, []string{"test", "cmd", "--flag4=55.44"}, 1)
		assertExitCode(t, c3, []string{"test", "cmd", "--flag4=55.44", "--flag5=optional"}, 1)

		assertExitCode(t, c2, []string{"test", "cmd", "--flag1=string"}, 1)
	})

	t.Run("exit with code 0 when all required flags are passed", func(t *testing.T) {
		assertExitCode(t, c1, []string{"test", "cmd", "--flag1=string", "--flag3=5"}, 0)
		assertExitCode(t, c1, []string{"test", "cmd", "--flag3=5", "--flat1=string"}, 0)

		assertExitCode(t, c2, []string{"test", "cmd", "--flag1=string", "--flag6=6", "--flag4=11.123", "--flag5=optional", "--flag2=optional"}, 0)
		assertExitCode(t, c2, []string{"test", "cmd", "--flag4=11.123", "--flag6=6", "--flag1=string"}, 0)

		assertExitCode(t, c3, []string{"test", "cmd", "--flag4=22.44"}, 0)
	})
}

func TestRequiredArgs(t *testing.T) {
	c1 := cli()
	c1.AddGlobalArg("arg1", "Required arg", CLIFlagTypeString|CLIFlagRequired)
	c1.AddGlobalArg("arg2", "Optional arg", CLIFlagTypeFloat)
	c1.AddGlobalArg("arg3", "Required atg", CLIFlagTypeInt|CLIFlagRequired)
	_ = c1.AddCmd("cmd", "Sample command", h)

	c2 := cli()
	c2.AddGlobalArg("arg1", "Required arg", CLIFlagTypeString|CLIFlagRequired)
	c2.AddGlobalArg("arg2", "Optional arg", CLIFlagTypeFloat)
	c2.AddGlobalArg("arg3", "Required arg", CLIFlagTypeInt|CLIFlagRequired)
	cmd21 := c2.AddCmd("cmd", "Sample command", h)
	cmd21.AddArg("arg4", "Required arg", CLIFlagTypeFloat|CLIFlagRequired)
	cmd21.AddArg("arg5", "Optional arg", CLIFlagTypeString)
	cmd21.AddArg("arg6", "Required arg", CLIFlagTypeInt|CLIFlagRequired)
	cmd21.ExcludeGlobalArg("arg3")

	c3 := cli()
	cmd31 := c3.AddCmd("cmd", "Sample command", h)
	cmd31.AddArg("arg4", "Required arg", CLIFlagTypeFloat|CLIFlagRequired)
	cmd31.AddArg("arg5", "Optional arg", CLIFlagTypeString)

	t.Run("exit with code 1 when required global arg is missing", func(t *testing.T) {
		assertExitCode(t, c1, []string{"test", "cmd"}, 1)
		assertExitCode(t, c1, []string{"test", "cmd", "--arg1=string"}, 1)
		assertExitCode(t, c1, []string{"test", "cmd", "arg1", "55.33"}, 1)
		assertExitCode(t, c1, []string{"test", "cmd", "arg1"}, 1)

		assertExitCode(t, c2, []string{"test", "cmd", "11.11", "5"}, 1)
		assertExitCode(t, c2, []string{"test", "cmd", "11.11", "5", "garg1", "6"}, 1)
	})

	t.Run("exit with code 1 when required command arg is missing", func(t *testing.T) {
		assertExitCode(t, c3, []string{"test", "cmd"}, 1)
		assertExitCode(t, c3, []string{"test", "cmd", "arg5"}, 1)

		assertExitCode(t, c2, []string{"test", "cmd", "garg1"}, 1)
		assertExitCode(t, c2, []string{"test", "cmd", "garg1", "12.123"}, 1)
	})

	t.Run("exit with code 0 when all required args are passed", func(t *testing.T) {
		assertExitCode(t, c1, []string{"test", "cmd", "garg1", "5"}, 0)
		assertExitCode(t, c1, []string{"test", "cmd", "garg1", "5", "33.21"}, 0)

		assertExitCode(t, c2, []string{"test", "cmd", "garg1", "123.123", "6"}, 0)
		assertExitCode(t, c2, []string{"test", "cmd", "garg1", "123.123", "6", "33.21", "arg5"}, 0)

		assertExitCode(t, c3, []string{"test", "cmd", "44.33"}, 0)
		assertExitCode(t, c3, []string{"test", "cmd", "44.33", "arg5"}, 0)
	})

}

func TestArgsAndFlagsOrder(t *testing.T) {
	c2 := cli()

	c2.AddGlobalArg("arg1", "Required arg", CLIFlagTypeString|CLIFlagRequired)
	c2.AddGlobalArg("arg2", "Optional arg", CLIFlagTypeFloat)
	c2.AddGlobalArg("arg3", "Required arg", CLIFlagTypeInt|CLIFlagRequired)

	c2.AddGlobalFlag("flag1", "Required flag", "", CLIFlagTypeString|CLIFlagRequired)
	c2.AddGlobalFlag("flag2", "Optional flag", "", CLIFlagTypeInt)
	c2.AddGlobalFlag("flag3", "Required flag", "ignored", CLIFlagTypeInt|CLIFlagRequired)

	cmd21 := c2.AddCmd("cmd", "Sample command", h)
	cmd21.AddArg("arg4", "Required arg", CLIFlagTypeFloat|CLIFlagRequired)
	cmd21.AddArg("arg5", "Optional arg", CLIFlagTypeString)
	cmd21.AddArg("arg6", "Required arg", CLIFlagTypeInt|CLIFlagRequired)
	cmd21.ExcludeGlobalArg("arg3")

	cmd21.AddFlag("flag4", "Required flag", "", CLIFlagTypeFloat|CLIFlagRequired)
	cmd21.AddFlag("flag5", "Optional flag", "", CLIFlagTypeString)
	cmd21.AddFlag("flag6", "Required flag", "ignored", CLIFlagTypeInt|CLIFlagRequired)
	cmd21.ExcludeGlobalFlag("flag3")

	t.Run("exit with code 1 when flags and args are passed in wrong order", func(t *testing.T) {
		assertExitCode(t, c2, []string{"test", "cmd", "garg1", "44.44", "6", "--flag6=6", "--flag4=123.123", "--flag1=str"}, 1)
		assertExitCode(t, c2, []string{"test", "cmd", "--flag6=6", "garg1", "--flag4=123.123", "44.44", "--flag1=str", "6"}, 1)
	})

	t.Run("exit with code 0 when flags and args are passed in correct order", func(t *testing.T) {
		assertExitCode(t, c2, []string{"test", "cmd", "--flag1=str", "--flag2=5", "--flag4=123.123", "--flag5=string", "--flag6=50", "garg1", "44.44", "6", "33.33", "arg5"}, 0)
		assertExitCode(t, c2, []string{"test", "cmd", "--flag6=6", "--flag4=123.123", "--flag1=str", "garg1", "44.44", "6"}, 0)
	})
}

// test single
// test with many

func TestStrFlagType(t *testing.T) {
}

func TestIntFlagType(t *testing.T) {
}

func TestFloatFlagType(t *testing.T) {
}

func TestAlphanumericFlagType(t *testing.T) {
}

func TestBoolFlagType(t *testing.T) {
}

func TestPathFileFlagType(t *testing.T) {
}
