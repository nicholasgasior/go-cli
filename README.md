# go-cli

Package `nmls/go-cli` is meant to make handling command line interface easier. 

You define commands with flags, attach a handler to it and package does all
the parsing.

### Install

Ensure you have your 
[workspace directory](https://golang.org/doc/code.html#Workspaces) created and
run the following:

```
go get -u github.com/nmls/go-cli
```

### Example

Let's start with an example covering everything. First, let's create main
`CLI` instance and commands:

```
func main() {
    myCLI      := cli.NewCLI("Example CLI", "Silly app", "Author <a@example.com>")

    cmdHello   := myCLI.AddCmd("hello", "Prints out Hello", HelloHandler)
    cmdJSONKey := myCLI.AddCmd("json_key", "Checks if JSON file has key", JSONKeyHandler)
}
```

Next, let's add flags to our commands:

```
    cmdJSONKey.AddFlag("json-file", "JSON file", cli.CLIFlagTypePathFile | cli.CLIFlagMustExist | cli.CLIFlagRequired)
    cmdJSONKey.AddFlag("json-key", "JSON key", cli.CLIFlagTypeString | cli.CLIFlagRequired)

    cmdHello.AddFlag("language", "Language", cli.CLIFlagTypeString | cli.CLIFlagRequired)
    cmdHello.AddFlag("color", "Add color", cli.CLIFlagTypeBool)
```

Third argument to `NewCLIFlag` is used to define what is the type of flag, is
it required etc. It's an integer value and the following `const`s are
available:

* `CLIFlagTypePathFile` - flag is a path to a file (string);
* `CLIFlagMustExist` - if added along with `CLIFlagTypePathFile` then path must exist;
* `CLIFlagRequired` - flag is required (this does not work with bool flag);
* `CLIFlagTypeString` - flag is a string;
* `CLIFlagTypeBool` - flag is boolean and will have a value of "true" or "false".

Finally, let's create functions to handle our commands. In below code, you can
see that method `Flag` on `CLI` instance (passed as first argument) can be
used to get a flag value.

```
func HelloHandler(c *cli.CLI) int {
    fmt.Fprintf(os.Stdout, "Language: %s\n", c.Flag("language"))
    fmt.Fprintf(os.Stdout, "Color: %s\n", c.Flag("color"))

    return 0
}

func JSONKeyHandler(c *cli.CLI) int {
    fmt.Fprintf(os.Stdout, "JSON key: %s\n", c.Flag("json-key"))
    fmt.Fprintf(os.Stdout, "JSON file: %s\n", c.Flag("json-file"))
    return 0
}
```

And in the end of `main()` func:

```
    os.Exit(myCLI.Run())
```

