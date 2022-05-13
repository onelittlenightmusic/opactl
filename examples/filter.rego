package opactl

filter = { line |
  line := input.stdin[_]
  texts := split(line, " ")
  texts[0] == input.mod
}

