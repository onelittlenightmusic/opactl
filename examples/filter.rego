package opactl

__filter = "Filter ls output with mode label (parameter mod=\"-rwxr-xr-x\")"

filter = { line |
  line := input.stdin[_]
  texts := split(line, " ")
  texts[0] == input.mod
}

__filter = "Filter ls output with mode label (parameter mod=\"-rwxr-xr-x\")"

json_filter = { fruit_name |
  not json_error_found
  fruit_spec := input.json_stdin[fruit_name]
  fruit_spec.sweetness == "high"
}

json_error_found {
  input.json_stdin.error
  print("[Error]", input.json_stdin.error)
}