package opactl

# object
get_test_object = {
  "test": "test"
}

get_first_line = rtn {
  rtn := input.stdin[0]
} else = {}
# To define default return value is strongly recommended.

# set (Kind of list. Elements are unique. No order.)
select_unique_lines[rtn] {
  rtn := input.stdin[_]
}

# array (Kind of list. Elements are not necessary unique. The order is preserved.)
lines = [rtn|
  rtn := input.stdin[_]
]