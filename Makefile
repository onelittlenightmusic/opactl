opactl_examples=opactl -i examples

build:
	go build
	sudo cp opactl /usr/local/bin/
test:
	echo '{"orange": {"sweetness": "high"}, "cherry": {"sweetness":"middle"}}' \
	| $(opactl_examples) filter json_filter -p sweetness=high \
	| $(opactl_examples) test all