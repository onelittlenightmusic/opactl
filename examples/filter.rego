package opactl

__filter = "Filter ls output with mode label (parameter mod=\"-rwxr-xr-x\")"

filter = { line |
  line := input.stdin[_]
  texts := split(line, " ")
  texts[0] == input.mod
}

