package opactl.examples.filter

__comment = "Examples of rules filtering stdin"

__filter = "Filter ls output with mode label (parameter mod=\"-rwxr-xr-x\")"

filter = { line |
  line := input.stdin[_]
  texts := split(line, " ")
  texts[0] == input.mod
}

__json_filter = "Filter json output  (parameter sweetness=\"high\")"

json_filter = { fruit_name |
  not json_error_found
  fruit_spec := input.json_stdin[fruit_name]
  fruit_spec.sweetness == input.sweetness
}

json_error_found {
  input.json_stdin.error
  print("[Error]", input.json_stdin.error)
}