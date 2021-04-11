## Objectives

By the end of this session, you should understand:

* Challenges and tradeoffs involved in designing a file format
* Postgres heap table file format
* (If time allows) How Postgres uses a buffer cache to speed up file access

## Agenda

* Discuss prework

	* Ansel: some difficulties:
		* Accounting for future proofing?
			* Resizing stuff, variable length fields
		* Instead went with basic version of what Postgres does
			* Pointers at front, data at back of page
	* Lady Red:
		* How to handle variable schema?
			* enum for type
			* Record takes pointer for schema
		* Representing schema on disk:
			* array of enums
			* Could potentially write schema to file
		* How to interact with type system?
	* Stephen:
		* Process of coming up with file format seems straightforward
			* For example:
				* Specify page size
					* Do we want to split our data into pages?
				* Specify # of records on page
					* 
				* Use (offset, length) tuples within a page
		* Going from there to implementation seems very tricky
			* What directory?
			* How to implement the FileScan node?
			* Catalog for storing schema?

* Inspect heap files

	* What's the structure of Postgres heap table files?


	* Files divided into 8kb pages
	* Pages have a header, line pointers, free space, heap tuples
	* Deep dive into each part (some will require context from later sessions)
	* What happens if we do these operations, and what's the cost?
		* Insert
		* Update
		* Delete
	* How does a scan work?
	* How do we handle variable-sized 

```
pg_filedump -f -i -D int,varchar,varchar ... | less
```

* Experiment with buffer usage

	* Shut down and restart Postgres to wipe out cache

```
SELECT n.nspname, c.relname, count(*) AS buffers
             FROM pg_buffercache b JOIN pg_class c
             ON b.relfilenode = pg_relation_filenode(c.oid) AND
                b.reldatabase IN (0, (SELECT oid FROM pg_database
                                      WHERE datname = current_database()))
             JOIN pg_namespace n ON n.oid = c.relnamespace
             WHERE n.nspname = 'public'
             GROUP BY n.nspname, c.relname
             ORDER BY 3 DESC;
```

## To-Do

* Durability

## Resources

How Postgres handles very large fields:
	https://wiki.postgresql.org/wiki/TOAST

Considerations around ensuring durability when writing to disk:
	https://www.postgresql.org/docs/current/wal-reliability.html

E-book with details about file structure
	https://www.interdb.jp/pg/pgsql01.html

https://malisper.me/the-file-layout-of-postgres-tables/

Heap Only Tuples
	https://www.interdb.jp/pg/pgsql07.html
	https://github.com/postgres/postgres/blob/master/src/backend/access/heap/README.HOT
