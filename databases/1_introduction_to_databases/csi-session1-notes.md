## Objectives

By the end of this session, you should:

* Understand the path a SQL query takes through Postgres
* Know the important high-level components of an RDBMS
* Identify some personal open questions / areas of interest for the rest of the course

## Questions

* Differences between RDBMS and NoSQL?
	* We can do it informally as they come up
	* But probably will be talking about his more in distributed systems / system design

* What do "Simple Query" / "Extended Query" mean? (Details of wire protocol)

* Debugger commands

* How did I jump to nodeSeqscan

## Explorations

### Part 0, Test Data

### Part 1, Network

Launch Wireshark
Start capturing Loopback traffic
Set pgsql filter
psql -h 127.0.0.1 -p 5432 csi
SELECT * FROM foo;
Inspect the requests and responses

### Part 2, File Formats

	# pg_waldump -f <WAL file>
	# pg_filedump -i -f -D int,varchar,age <heap file> | less
	# pg_filedump -i -f <index file> | less

	CREATE TABLE foo (id int, name varchar(255), age smallint);
	SELECT pg_relation_filepath('foo');

	BEGIN;
	SELECT txid_current();
	INSERT INTO foo (id, name, age) VALUES (...);
	ABORT;

	BEGIN;
	SELECT txid_current();
	INSERT INTO foo (id, name, age) VALUES (...);
	COMMIT;

	CHECKPOINT;

	CREATE INDEX idx_foo_name ON foo (name);
	SELECT pg_relation_filepath('idx_foo_name');

	INSERT INTO foo (id, name, age) VALUES(...);

	CHECKPOINT;

	UPDATE foo SET name = '...' WHERE id = ...;

	CHECKPOINT;

* Why would there be so much empty space in this file?
	* Whole page with only one record

* Why is the record at the end?
	* Metadata about the rows at the top
		* Maybe size, handling nulls

### Part 3, Debugger

	psql -h 127.0.0.1 -p 5432 csi
	ps x | grep postgres
	lldb
	attach <pid>
	b ExecutorRun
	b ExecLimit
	b ExecSeqScan
	SELECT * FROM foo LIMIT 5;
	bt
	p *queryDesc
	c
	bt
	p *pstate->plan
	p *pstate->lefttree
	p *pstate->lefttree->plan
	c
	bt

	p *slot
		Look up that data in the heap file
	p *node
		offset, position, count

### Other things we can try?

* kill postgres in the middle of a transaction
* kill postgres before doing CHECKPOINT
* DELETE
* import a larger dataset and look at some query plans

## Discussion

