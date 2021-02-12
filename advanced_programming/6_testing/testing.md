Make-up session?
	- Wednesday 2/10
	- Wednesday 2/17 (this one won)

---

By the end of this session you should be able to understand the capabilities and "philosophy" of:
- Go's testing framework in general
- Property-based testing

You should also be able to:
- Use Go's testing package / tool for writing tests
- Identify opportunities to gradually incorporate "property-based testing" into your projects

---

What are some general principles for writing good tests / avoiding bad tests?
	Prefer to start with "black-box" testing
		Want specific reason to do "white-box" / "clear-box" testing
	Having verbose messaging when you have failures
	Not trying to test too much
		Only look at properties of complex object that matter
		Don't do exact string matching, look for patterns
		Avoid "brittle" tests
	Jake's principles for unit vs. integration tests: only test one boundary at a time
		e.g. web app with frontend
			http request to backend
			route config
			auth layer
			business logic, talks to database
		instead of launching a backend and db, hitting backend with network request
			test each boundary separately
	Lady Red's strong opinions about mocking:
		Don't make every test write their own separate mock of the same component
		Consider offering a mock in addition to full component
			KVStoreService
				(need to initialize with host:port)
				Get -> network call
			FakeKVStoreService
				(no initialization data needed)
				Get -> in memory in a map or something
		What's the value of this?
			If someone writes own mock for KVStoreService
			Someone makes update to API (now Get requires version)

What’s the “table-driven” testing pattern in Go?
	One implementation of a test, data in a table
		Input, expected output, message?
	Alyssa: pytest parametrize
		decorator approach, expected output, list combinations in rows
	Stephen's consideration: split tables?
	Jake: haven't had good experience with Node, Jest
		harder to debug, stops on first failure
		bad feedback messages
			assertError: true != false

```python
"""
{
	"name": "alice-smith",
	"expected-callback": handle-hyphen,
}

Three kinds of things in the table:
	first name only
	first + last with no hyphen
	first + last with hyphen
"""
def testNames(last, first):
	if not last:
		pass
	else:


```

When might it make sense to consider random tests?
	Lady Red's "perfect example"
		Difference of two complex JSON objects ("subtract one from other")
		Also had like a "merge" function
	Use slow / easy implementation to compare
		Generate lots of random inputs

Follow-up what's the difference between a random test and a "property-based test"
	"property-based test" = random test + shrinking

What are some useful checks you can do with property-based testing?
	if you have some rules around "what's a valid object"
		you can make sure all your outputs are valid
		e.g. list is sorted
		e.g. JSON object conforms to some schema
	does it crash with arbitrary inputs?
	mathematical properties
		serialize / deserialize
		compress / decompress
		encrypt / decrypt
		idempotency: doing an operation twice, second time does nothing
		order doesn't matter
	other invariants from Oz's John Hughes video:
		doing things in a different order?
	Ansel: SQLite example
		Fuzzing
		Alyssa: FactoryBoy, for fuzzy data

Stephen: using these ideas in a broader context
Jake: ChaosMonkey for disaster recovery / emergency response

When might it make sense to create a separate package for tests?

What’s “coverage”, and to what extent is it useful to measure?

In 11.2.3 "White Box Testing" the book uses the idea of swapping in test functionality by having global variables. What do you think of this approach?

---

How can you incrementally introduce property-based testing?
	Find the right use case
		Oz: nice API that hides lots of complexity?
		If your unit test has arbitrary hardcoded stuff, generate it instead
		Really well-defined acceptance criteria

---

it should (be sorted)
vs.
it should equal [1, 1, 2, 6, 9]
fix the random seed

What is property-based testing?
What are some advantages / disadvantages relative to "typical" unit tests?

---

Exercises
- Mode
- Product
- Range Count
- Skip List

Real-world examples?
- Standard library?

Advanced / real-world example:
	https://www.cis.upenn.edu/~bcpierce/papers/mysteriesofdropbox.pdf

---

Interesting Twitter threads:
	https://twitter.com/hillelogram/status/1190408364725067776
	https://twitter.com/DRMacIver/status/1140615803752079360

Nice examples / writeup: 
	https://increment.com/testing/in-praise-of-property-based-testing/
