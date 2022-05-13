package opactl.hierarchy.label

import data.sample_data

index = { k: summary |
  pods := sample_data.pods
  pods[_].label[k]
  summary := { v: podlist |
    v := pods[_].label[k]
    podlist = {pod| pods[pod].label[k] == v}
  }
}

count_of = { k: count_summary |
  summary := index[k]
  count_summary := { v: counts |
    counts := count(summary[v])
  }
}
# $ opactl hierarchy label count_of app analytics
# 1
# $ opactl hierarchy label count_of app
# {
#   "analytics": 1,
#   "backend": 3,
#   "frontend": 2
# }
## Label auto completion
# $ opactl hierarchy label index <tab><tab>
# $ opactl hierarchy label index access https
# [
#   "nginx"
# ]
