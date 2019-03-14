package main

func PrintHelloHandler() int {
    return 0
}

func CheckJSONKeyHandler() int {
    return 0
}

func main() {
    cli := NewCLI("Example CLI", "Simple CLI application", "Nicholas Gasior <nameless@laatu.org>")

    cmdHello := NewCLICmd("print_hello", "Prints out Hello in specified language", PrintHelloHandler)
    cmdJSONKey := NewCLICmd("check_json_key", "Checks if JSON file has key", CheckJSONKeyHandler)

    flagJSONFile := NewCLIFlag("json-file", "JSON file", CLIFlagTypePathFile | CLIFlagMustExist | CLIFlagRequired)
    flagJSONKey := NewCLIFlag("json-key", "JSON key", CLIFlagTypeString | CLIFlagRequired)

    flagLanguage := NewCLIFlag("language", "Language", CLIFlagTypeString | CLIFlagRequired)
    flagColor := NewCLIFlag("color", "Add color", CLIFlagTypeBool)

    cmdHello.AddFlag(flagLanguage)
    cmdHello.AddFlag(flagColor)

    cmdJSONKey.AddFlag(flagJSONFile)
    cmdJSONKey.AddFlag(flagJSONKey)

    cli.AddCmd(cmdHello)
    cli.AddCmd(cmdJSONKey)
    cli.Run()
}

