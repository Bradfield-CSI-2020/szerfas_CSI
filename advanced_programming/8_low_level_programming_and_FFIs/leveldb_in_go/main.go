package main

/*
#cgo CFLAGS: -I/Users/szerfas/src/bradfield/szerfas_CSI/advanced_programming/8_low_level_programming_and_FFIs/leveldb_in_go/leveldb/include
#cgo LDFLAGS: -L/usr/local/Cellar/leveldb/1.22/lib -lleveldb
#include "leveldb/c.h"
 */
import "C"

func main() {
	db := DB{}
	db.Open()
}

type DB struct {
	db *C.leveldb_t
	options *C.leveldb_options_t
	write_options *C.leveldb_writeoptions_t
}

func (db *DB) Open() {
	options := C.leveldb_options_create()
	// confirm we want this

	// figure out what uint8 to pass into this function, e.g., (C.uint8)1
	C.leveldb_options_set_create_if_missing(options, C.uint8_t(1))
	db.options = options

	db.write_options = C.leveldb_writeoptions_create()

	name := (C.CString)("path_to_dir")

	// similar to array of strings
	var errStr *C.char

	db.db = C.leveldb_open(options, name, &errStr)
}

func (db *DB)  Put(key []byte, val []byte) C.leveldb_put {
	C.leveldb_put()
	(leveldb_t* db, const leveldb_writeoptions_t* options, const char* key, size_t keylen, const char* val, size_t vallen, char** errptr);
}

/*
Questions:
- What does LEVELDB_EXPORT mean? It means that this function is either exported or imported depending on the definition of another variable.
For our purposes, this is likely an export.
*/



// methods to work with:
// open (skip and just do in C w/o go wrapper for now?)
// put
// get


// import leveldb C library
// start database:
// write wrapper in c that takes our Go primitives and assembles them into C structs that we can then pass into C API

// note: we need to be able to take C.x notation primitives and assemble them into the appropraite C structures
// then we need to be able to take c structures and disassemble them into C.x notation primitives (which can be re-assembled into meaningful Go data structures)

// UPDATE: only need to do this if we need to stick something in there and pull it out (which we hopefully do not need to do)
// to use open, we need to have leveldb_options, which we can create using leveldb_options_create
// -----
// we need to take the "leveldb_t*" output of the open function and make it meaningful Go object (or at last have our other Go wrappers reference it)
// struct leveldb_t {
//  DB* rep;
//};
//


// implement a wrapper like bz2compress below for the leveldb function

//LEVELDB_EXPORT leveldb_t* leveldb_open(const leveldb_options_t* options,
//                                       const char* name, char** errptr);

//int bz2compress(bz_stream *s, int action,
//                char *in, unsigned *inlen, char *out, unsigned *outlen) {
//  s->next_in = in;
//  s->avail_in = *inlen;
//  s->next_out = out;
//  s->avail_out = *outlen;
//  int r = BZ2_bzCompress(s, action);
//  *inlen -= s->avail_in;
//  *outlen -= s->avail_out;
//  s->next_in = s->next_out = NULL;
//  return r;
//}

// implement function here