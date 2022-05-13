package opactl.hierarchy

import data.sample_data
get = sample_data

search = [ name |
  pod := sample_data.pods[name]
  pod.label.app == input.label
]

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