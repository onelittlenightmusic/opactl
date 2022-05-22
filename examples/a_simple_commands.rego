package opactl.examples.a_simple

__comment = "Simple examples."

__simple_object = "Return simple object defined with JSON"
simple_object = {
  "test": "test"
}

__first_line = "Return simple filter returning the first line of stdin (Requires stdin)"
first_line = rtn {
  rtn := input.stdin[0]
} else = {}
# To define default return value is strongly recommended.

__set = "set (Kind of list. Elements are unique. No order.)(Requires stdin)"
set[rtn] {
  rtn := input.stdin[_]
}

__array = "array (Kind of list. Elements are not necessary unique. The order is preserved.)(Requires stdin)"
array = [rtn|
  rtn := input.stdin[_]
]