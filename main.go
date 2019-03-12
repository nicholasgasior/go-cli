package main

func HelloWorldHandler() int {
    return 0
}

func AhoyHandler() int {
    return 0
}

func main() {
    cli := NewCLI("Example CLI", "Simple CLI application", "Nicholas Gasior <nameless@laatu.org>")

    cmdHelloWorld := NewCLICmd("hello_world", "Prints out hello world", HelloWorldHandler)
    cmdAhoy := NewCLICmd("ahoy", "Prints ahoy", AhoyHandler)

    flagConfig := NewCLIFlag("config", "Config file", CLIFlagTypePathFile | CLIFlagMustExist | CLIFlagRequired)
    flagColorful := NewCLIFlag("color", "Make it colorful", CLIFlagTypeBool)

    cmdHelloWorld.AddFlag(flagConfig)
    cmdAhoy.AddFlag(flagConfig)
    cmdAhoy.AddFlag(flagColorful)

    cli.AddCmd(cmdHelloWorld)
    cli.AddCmd(cmdAhoy)
    cli.Run()
}

