# opactl

`opactl` executes your own OPA policy as command. 

This is how it works. You define a rule in OPA policy, for example `rule1`. Then, `opactl` detects your rule and turns it into subcommand such as `opactl rule1`.

Options are supported for various usage. Also, you can preset configuration file, then `opactl` reads it.

## Execute a rule as subcommand

When you define a rule `filter` as follows, 

```rego
package opactl

# pick up only lines which includes specific mod
filter = { line |
  # load each line of stdin
  line := input.stdin[_]
  # split into words
  texts := split(texts, " ")
  # check the first word equals to parameter `mod`
  texts[0] == input.mod
}
```

you can run a subcommand `opactl filter` like this.

```sh
# Run subcommand filter with using stdin (-i) and parameter (mod=...)
ls -l | opactl -i filter -p mod="-rwxr-xr-x"
[
  "-rwxr-xr-x  1 hiroyukosaki  staff  8055840 May 12 01:04 opactl"
]
```

## Options

```
Flags:
  -a, --all                 Show all commands
  -b, --base string         OPA base path which will be evaluated (default "data.opactl")
      --config string       config file (default is $HOME/.opactl.yaml)
  -d, --directory strings   directories
  -h, --help                help for opactl
  -i, --input               Accept stdin as input.stdin
  -p, --parameter strings   parameter (key=value)
  -q, --query string        Input your own query script (example: { rtn | rtn := 1 }
  -v, --verbose             Toggle verbose mode on/off
```

## Configuration

You can create an `.opactl` configuration file. When you run `opactl` command in the same directory, `opactl` loads the configuration and set options. 

Each field in `.opactl` is connected to one option. For example, `parameter` field is read as `--parameter` option.

```
directory:
- examples
base: data.opactl
parameter:
- item=1
```