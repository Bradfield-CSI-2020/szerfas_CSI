
##Optimizing a Gaussian blur
#### Memory needs and cache sizes
######Estimating kernal memory allocation:
```
9 float64 per row		9 per row * 8 bytes per float = 72 bytes per row
9 rows				9 * 72 = 648 bytes in the kernel matrix
```
######Estimating image memory allocation:
```
image dimension 1280 x 1631		implies i max = 1280 and j max = 1631
each are floats, so:			1280 * 64 = 81920, 82kb bytes on x axis (# columns in a row, equal to i)
				        1631 * 64 = 104384 bytes, 104kb on y axis (# rows, equal to j)

12380 * 1631 * 64 = 133611520 bytes, 133MB to make up the image
It's only 94kb on disk due to JPG compression, which I'm assuming are re-expanded when we're working on it 
```
Sequence of thoughts:
* I think this implies that we're looping out of order! Should go in order of j then i, not i then j
* BUT 64 bytes * 104K rows is still only 6Mb or so, which is < the 12Mb we can store in L3 cache. This implies we will NOT see a change in performance due to traversing the image in column order rather than row order first.
* Well, that's only true if we're concerned about the difference between the L3 cache and main memory. If we want to "tile" by our L1 cache (which still  sees ~10x improvement over L3), we can only go  32kb total cache size / 64 bytes per row = 500 rows
* Let's take  this L1 optimization idea further. What is everything we need to store in it at one time to prevent misses? 
The kernel, plus a few variables for calculations, then the remainder can be the image values. 
So let's take 32 kb cache size - 1024 bytes (648 bytes + a buffer) = 31kb to pull from the image at once


######My system's memory hierarchy
Output from `sudo sysctl -a | grep cache | grep size`
```
hw.cachesize: 34359738368 32768 262144 12582912 0 0 0 0 0 0
hw.cachelinesize: 64				# 64 byte cacheline - so each float takes a full cacheline
hw.l1icachesize: 32768
hw.l1dcachesize: 32768			        # 32 Kb in l1 cache, meaning we can hold 5x the kernel in L1 alone!
hw.l2cachesize: 262144			        # 256 Kb in L2 cache
hw.l3cachesize: 12582912			# 12 MB in l3 cache
machdep.cpu.cache.linesize: 64
machdep.cpu.cache.size: 256
```

######Questions captured while working
* How might we align a float to the cacheline?
* Why are we right shifting by 8 in rr += float64(r>>8) * weight?
* Why does my benchmark test function exit in .007s? It's clearly not running the blur function - I think it has to do with successful importing of main package. Perhaps I could run a basic test first.
For Rohan
* Can you walk me through how you first went about this: book first or experiment first? Initial thoughts? Sequence of thouhts?
* When did you come up with each of the ideas in your README?
* When did you implement?
* Did you remove after implementing?
* How did you split the kernel into cache aligned chunks?
* How did you split the kernel into 1/4 it's size and using it's symmetry?
* How did you split the image into cache-aligned tiles?
* How do you estimate we're using 192 bits per pixel?
* What exactly do you mean by destructuring in your code?
* Main ideas I took away from this: (1) we don't know what will optimize until we experiment, especially because there can be some hard-to-anticipate tradeoffs between computation and memory access patterns, (2) we might try restructuring the underlying datastore to see what happens, (3) go routines can increase speed
* What is waitgroup?


Notes
* Parallelism - cores
* Cache locality
* Loop unrolling and pipelining
* 


######Ideas recorded live as I go
* loop unroll - where?
* loop ordering? tiling?
* temp variables? these seem fine
* look at assembly
* could we couple and reduce function calls?
* could we rewrite if-statement in conditional style?
* go routines to make anything concurrent? how?

######initial benchmark
Output from `time ./blur`:
```
real	0m6.384s
user	0m6.390s
sys	0m0.047s
```