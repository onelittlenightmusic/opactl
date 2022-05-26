# opactl in shell script

Of course, we can use `opactl` in shell scripts. Please try a sample shell script. You can try these examples. You can create your own rules after that.

```sh
./main.sh data1

[data/x2.yaml]
column21: test21
column22: test22
-------
[data/y2yyy.yaml]
all:
  new: second
  old: first
-------
```

This is how it works. `opactl` accepts the output of `ls` command (list of file names) and filters only file names which meets criteria. 

```sh
#!/bin/sh

dir=$1
for f in $(ls $dir | opactl search_spaced) 
...
```

In another case, `opactl` receives file contents and generates `true` or `false`. For example, `opactl grep` rule checks file contents include all keywords at the same time.

```sh
./grep.sh data1

data1/x2.yaml includes all keywords
```

```sh
keywords=column21,test21
...
  check=$(cat $file | opactl grep -P keywords=$keywords)
  if [ $check ]; then
    echo $file includes all keywords
  fi
```