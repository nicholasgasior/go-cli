package cli

import (
	"fmt"
	"os"
	"testing"
)

func HelloHandler(c *CLI) int {
	fmt.Fprintf(os.Stdout, "Language: %s\n", c.Flag("language"))
	fmt.Fprintf(os.Stdout, "Color: %s\n", c.Flag("color"))
	return 0
}

func JSONKeyHandler(c *CLI) int {
	fmt.Fprintf(os.Stdout, "JSON key: %s\n", c.Flag("json-key"))
	fmt.Fprintf(os.Stdout, "JSON file: %s\n", c.Flag("json-file"))
	return 0
}

func createCLI() *CLI {
	myCLI := NewCLI("Example CLI", "Silly app", "Author <a@example.com>")
	cmdHello := myCLI.AddCmd("hello", "Prints out Hello", HelloHandler)
	cmdJSONKey := myCLI.AddCmd("json_key", "Checks if JSON file has key", JSONKeyHandler)
	cmdJSONKey.AddFlag("json-file", "JSON file", CLIFlagTypePathFile|CLIFlagMustExist|CLIFlagRequired)
	cmdJSONKey.AddFlag("json-key", "JSON key", CLIFlagTypeString|CLIFlagRequired)
	cmdHello.AddFlag("language", "Language", CLIFlagTypeString|CLIFlagRequired)
	cmdHello.AddFlag("color", "Add color", CLIFlagTypeBool)
	cmdHello.AddFlag("many-int-default", "Many int separated with comma", CLIFlagTypeInt|CLIFlagAllowMany)
	cmdHello.AddFlag("many-int-colon", "Many int separated with colon", CLIFlagTypeInt|CLIFlagAllowMany|CLIFlagManySeparatorColon)
	cmdHello.AddFlag("many-int-semi-colon", "Many int separated with semi-colon", CLIFlagTypeInt|CLIFlagAllowMany|CLIFlagManySeparatorSemiColon)
	cmdHello.AddFlag("many-float-default", "Many float separated with comma", CLIFlagTypeFloat|CLIFlagAllowMany)
	cmdHello.AddFlag("many-float-colon", "Many float separated with colon", CLIFlagTypeFloat|CLIFlagAllowMany|CLIFlagManySeparatorColon)
	cmdHello.AddFlag("many-float-semi-colon", "Many float separated with semi-colon", CLIFlagTypeFloat|CLIFlagAllowMany|CLIFlagManySeparatorSemiColon)
	cmdHello.AddFlag("many-alphanumeric-default", "Many alphanumeric separated with comma", CLIFlagTypeAlphanumeric|CLIFlagAllowMany)
	cmdHello.AddFlag("many-alphanumeric-colon", "Many alphanumeric separated with colon", CLIFlagTypeAlphanumeric|CLIFlagAllowMany|CLIFlagManySeparatorColon)
	cmdHello.AddFlag("many-alphanumeric-semi-colon", "Many alphanumeric separated with semi-colon", CLIFlagTypeAlphanumeric|CLIFlagAllowMany|CLIFlagManySeparatorSemiColon)
    cmdHello.AddFlag("many-alphanumeric-semi-colon-underscore-dot", "Many alphanumeric separated with semi-colon allowing underscore and dot", CLIFlagTypeAlphanumeric|CLIFlagAllowMany|CLIFlagManySeparatorSemiColon|CLIFlagAllowDots|CLIFlagAllowUnderscore)
	return myCLI
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

func TestFlags(t *testing.T) {
	cli := createCLI()

	t.Run("exit with code 1 when flags are missing", func(t *testing.T) {
		assertExitCode(t, cli, []string{"test", "hello"}, 1)
	})
	t.Run("exit with code 0 when none flags are missing", func(t *testing.T) {
		assertExitCode(t, cli, []string{"test", "hello", "--language", "fi"}, 0)
	})
	t.Run("exit with code 1 when file path flag points to a missing file", func(t *testing.T) {
		assertExitCode(t, cli, []string{"test", "json_key", "--json-key", "key1", "--json-file", "/var/non-existing"}, 1)
	})
	t.Run("exit with code 0 when file path flag points to an existing file", func(t *testing.T) {
		assertExitCode(t, cli, []string{"test", "json_key", "--json-key", "key1", "--json-file", "cli.go"}, 0)
	})
	t.Run("exit with code 0 when none flags are missing and once flag is allowed to have comma-separated integers", func(t *testing.T) {
		assertExitCode(t, cli, []string{"test", "hello", "--language", "fi", "--many-int-default", "44,45,66,56,888"}, 0)
	})
	t.Run("exit with code 0 when none flags are missing and once flag is allowed to have colon-separated integers", func(t *testing.T) {
		assertExitCode(t, cli, []string{"test", "hello", "--language", "fi", "--many-int-colon", "44:45:66:56:888"}, 0)
	})
	t.Run("exit with code 0 when none flags are missing and once flag is allowed to have semi-colon-separated integers", func(t *testing.T) {
		assertExitCode(t, cli, []string{"test", "hello", "--language", "fi", "--many-int-semi-colon", "44;45;66;56;888"}, 0)
	})
	t.Run("exit with code 0 when none flags are missing and once flag is allowed to have comma-separated floats", func(t *testing.T) {
		assertExitCode(t, cli, []string{"test", "hello", "--language", "fi", "--many-float-default", "44.4,45.2,66.333,56.444,888.222"}, 0)
	})
	t.Run("exit with code 0 when none flags are missing and once flag is allowed to have colon-separated floats", func(t *testing.T) {
		assertExitCode(t, cli, []string{"test", "hello", "--language", "fi", "--many-float-colon", "44.1:45.3:66.2:56.33:888.11"}, 0)
	})
	t.Run("exit with code 0 when none flags are missing and once flag is allowed to have semi-colon-separated floats", func(t *testing.T) {
		assertExitCode(t, cli, []string{"test", "hello", "--language", "fi", "--many-float-semi-colon", "44.5;45.3;66.2;56.3333;888.44440"}, 0)
	})
	t.Run("exit with code 0 when none flags are missing and once flag is allowed to have comma-separated alphanumerics", func(t *testing.T) {
		assertExitCode(t, cli, []string{"test", "hello", "--language", "fi", "--many-alphanumeric-default", "abc44,de45,eeee66,56,888"}, 0)
	})
	t.Run("exit with code 0 when none flags are missing and once flag is allowed to have colon-separated alphanumerics", func(t *testing.T) {
		assertExitCode(t, cli, []string{"test", "hello", "--language", "fi", "--many-alphanumeric-colon", "fff44:ddd45:66:56ee:888"}, 0)
	})
	t.Run("exit with code 0 when none flags are missing and once flag is allowed to have semi-colon-separated alphanumerics", func(t *testing.T) {
		assertExitCode(t, cli, []string{"test", "hello", "--language", "fi", "--many-alphanumeric-semi-colon", "44;4afaw5;66;56;8d8s8"}, 0)
	})
	t.Run("exit with code 0 when none flags are missing and once flag is allowed to have semi-colon-separated alphanumerics with underscore and dots", func(t *testing.T) {
		assertExitCode(t, cli, []string{"test", "hello", "--language", "fi", "--many-alphanumeric-semi-colon-underscore-dot", "44;4afaw5;6___6;5.6;8d__._8s8"}, 0)
	})
}
