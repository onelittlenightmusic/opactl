package opactl.examples.hierarchy

import data.sample_data

__comment = "Example of how to process data and input."

__get = "Data loaded from YAML file."

get = sample_data

__search_label = "Rule which processes YAML data and searches elements. (parameter: app=frontend)"

search_label = [ name |
  pod := sample_data.pods[name]
  pod.label.app == input.app
]

__app_labels = "Rule which collects all labels from YAML data."

app_labels = { label |
  label := sample_data.pods[_].label.app
}

## Output examples:
# $ opactl hierarchy get -a
#[
#  "pods"
#]
# $ opactl hierarchy get pods -a
#[
#  "mariadb",
#  "mysql",
#  "nginx",
#  "nodejs",
#  "postgresql",
#  "python3"
#]
# $ opactl hierarchy get pods nginx label
#{
#  "app": "frontend"
#}
# $ opactl hierarchy search -p label=backend
#[
#  "mariadb",
#  "mysql",
#  "postgresql"
#]
# $ opactl hierarchy search -p label=frontend
#[
#  "nginx",
#  "nodejs"
#]