package cli

import (
    "testing"
    "os"
    "fmt"
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
    myCLI      := NewCLI("Example CLI", "Silly app", "Author <a@example.com>")
    cmdHello   := myCLI.AddCmd("hello", "Prints out Hello", HelloHandler)
    cmdJSONKey := myCLI.AddCmd("json_key", "Checks if JSON file has key", JSONKeyHandler)
    cmdJSONKey.AddFlag("json-file", "JSON file", CLIFlagTypePathFile | CLIFlagMustExist | CLIFlagRequired)
    cmdJSONKey.AddFlag("json-key", "JSON key", CLIFlagTypeString | CLIFlagRequired)
    cmdHello.AddFlag("language", "Language", CLIFlagTypeString | CLIFlagRequired)
    cmdHello.AddFlag("color", "Add color", CLIFlagTypeBool)
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
}
