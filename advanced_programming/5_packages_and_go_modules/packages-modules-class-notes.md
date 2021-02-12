Note: these are Elliott's notes he took while conducting the class, not mine.

## Objectives

By the end of this class you should understand some challenges of managing large projects, and Go's approach.

In particular, you should be comfortable with:

- Go's approach to handling imports, packages, dependencies
- Capabilities of the `go` tool
- Challenges of versioning and backwards compatibility, and how Go's approach has evolved over time

## Notes

* Mistakes Elliott has made in approaching this topic
	* compared to effective coworkers
* They're just tools
* Not glamorous but very important / effective
* Recommended reading: https://www.benkuhn.net/blub/

## Agenda

1. Lightning review of basic concepts

What does it mean that packages provide "namespaces" and "encapsulation"?
	namespacing
		example: Contains vs. StringContains inside a strings package
	encapsulation
		within a package, hide a lot of logic
		give user access to finite / small number of APIs
		in go:
			capital letters = exported

Aside: other examples of "conventions"
	Member variables start with m in Android codebase that Lady Red works on
	hungarian notation

What were some of the suggestions for good package / member names?
	conciseness
		strconv not stringconversionpackage

What sorts of things get exported from a package?
	structs, constants, variables, types, functions

What are the two meanings of an import path?
	How is the string used / interpreted?
		first part is organization (global uniqueness across ecosystem of packages)
		last part the package name matches the path
	import "github.com/.../caesar"
		can be used as a remote repository
			tool for downloading 3rd party dependencies uses it as a URL
		local file path
			tool for building locally

What are renaming, `_`, and . imports used for?
	_ used if not referencing package, but some part could rely on it
		Different types of image files
		Automatically calls Init method
	. import?
		import . "strings"

How does the name of a package relate to the import path?
	last portion of import path is name of package
	(pathological foo/bar example)

How does the download location relate to the import path?
	it goes to that location

What are the $GOPATH and $GOROOT directories?
	GOROOT: where language / sdk is
	GOPATH: context of dependencies
		src, pkg, bin?

What other relevant environment variables are there?
	

What do these commands do at a high level?
	go get
	go build
	go install

2. Address questions from Slack / weekly check-in

	What is "export data"?
		Listing info about a package (including where export data is):
			go list -export -json hash/fnv
		golang.org/x/tools/go/gcexportdata to dump export data

	What sort of caching does Go do?


	What was Go's original approach towards versioning?

3. Walkthrough of tools

	goimports
		How does it automatically find imports???
	go get
		Where does it download packages?
		How does it interact with version control tools?
		Where does it refuse if go.mod file doesn't match package import string?

Module question
	Why can't you go get@v2 if you've already done that at v1?
	If you create a v2 package with a different path, do you need to tag version still
		Does it matter
			github.com/alice/v2/caesar@v1.0.0
			github.com/alice/v2/caesar@v2.0.0

Questions:
	Is the compiler good enough to "prune" / not link functions / etc. that aren't used
	GO mod version numbering system / how it switches versions?
