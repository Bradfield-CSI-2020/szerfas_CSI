By the end of this session, you should be able to:

* Understand how Go's `reflect` package represents types and values at runtime
* Read and write functions that work on arbitrary, possibly unknown types
* Stay away from reflection unless absolutely necessary

Hopefully the discussion will also provide some clarification of types and interfaces!

## Exercises

- Warm up: Map

How would you implement "map" in Go?
	What's the best you could do without Reflect?
		What limitations?
	What can you do with Reflect?

Let's come back and discuss at 5:10 Pacific

	- Benchmark

- DeepEquals
	- NaN?
	- nil vs. empty?
	- Handling cycles?
- JSON encode / decode
	- struct tags
- Look at standard library versions

## Notes / discussion

Reflection, type switches, empty interface: what's the role of each?
Why doesn't `heap` use reflection?
	Why doesn't it work with `[]interface{}`?

How does TypeOf work?
	(see implementation in reflect/type.go)
Escaping dependency hell by copying code
	(see rtype struct in reflect/type.go)

Warnings
	nil interfaces
		(nil_iface.go)
	equal pointers but unequal types
		(arr.go)

## Use Cases

- fmt
- json
- template
- gopter
