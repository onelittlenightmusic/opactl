package opactl.examples.test

__comment = "test tool"

__is_an_array = "check if stdin is an array"
default is_an_array = false
is_an_array {
  is_array(input.json_stdin)
}


__has_orange = "check if stdin has orange"
default has_orange = false
has_orange {
  "orange" == input.json_stdin[_]
}

__all = "check all tests"
default all = false
all {
  is_an_array
  has_orange
}