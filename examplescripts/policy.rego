# METADATA
# description: opactl package
package opactl
import future.keywords

__comment = "Examples of rules using with other CLI commands"

__search = "Check file name with extension (.yaml) and word (\"index\" etc)"

# METADATA
# description: search rule
search = { file |
  file := __file_list[_]
  endswith(file, input.extension)
  contains(file, input.word)
}

__file_list = { _each_file |
  _files := split(input.stdin[i], " ")
  _each_file := _files[_]
}

search_spaced = concat(" ", search)

# METADATA
# description: has_keywords_in_line rule
grep {
  line := input.stdin[_]
  not keywords_missing
  every keyword in input.keywords {
    contains(line, keyword)
  }
}

keywords_missing {
  not is_array(input.keywords)
  print("[Error] Requires array parameter of keywords")
}
