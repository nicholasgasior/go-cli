# go-cli

### Sample usage

To create a simple command line interface you have to:

* import the package;
* create main CLI instance;
* create commands;
* create command flags and attach them to commands;
* attach commands to flags;
* call `Run()` on the CLI instance.

See below example:

```package main

import (
    "os"
    "fmt"
    "github.com/nmls/go-cli"
)

func PrintHelloHandler(cli1 *cli.CLI) int {
    fmt.Fprintf(os.Stdout, "Language: %s\n", cli1.Flag("language"))
    fmt.Fprintf(os.Stdout, "Color: %s\n", cli1.Flag("color"))

    return 0
}

func CheckJSONKeyHandler(cli1 *cli.CLI) int {
    fmt.Fprintf(os.Stdout, "JSON key: %s\n", cli1.Flag("json-key"))
    fmt.Fprintf(os.Stdout, "JSON file: %s\n", cli1.Flag("json-file"))
    return 0
}

func main() {
    cli1 := cli.NewCLI("Example CLI", "Simple CLI application", "Nicholas Gasior <nameless@laatu.org>")

    cmdHello := cli.NewCLICmd("print_hello", "Prints out Hello in specified language", PrintHelloHandler)
    cmdJSONKey := cli.NewCLICmd("check_json_key", "Checks if JSON file has key", CheckJSONKeyHandler)

    flagJSONFile := cli.NewCLIFlag("json-file", "JSON file", cli.CLIFlagTypePathFile | cli.CLIFlagMustExist | cli.CLIFlagRequired)
    flagJSONKey := cli.NewCLIFlag("json-key", "JSON key", cli.CLIFlagTypeString | cli.CLIFlagRequired)

    flagLanguage := cli.NewCLIFlag("language", "Language", cli.CLIFlagTypeString | cli.CLIFlagRequired)
    flagColor := cli.NewCLIFlag("color", "Add color", cli.CLIFlagTypeBool)

    cmdHello.AddFlag(flagLanguage)
    cmdHello.AddFlag(flagColor)

    cmdJSONKey.AddFlag(flagJSONFile)
    cmdJSONKey.AddFlag(flagJSONKey)

    cli1.AddCmd(cmdHello)
    cli1.AddCmd(cmdJSONKey)
    cli1.Run()
}
```
