Let’s think through the cache misses

Running sudo sysctl -a | grep cache | grep size gives: 
```
vfs.generic.nfs.server.reqcache_size: 64
hw.cachesize: 34359738368 32768 262144 12582912 0 0 0 0 0 0
hw.cachelinesize: 64
hw.l1icachesize: 32768		# 32kb
hw.l1dcachesize: 32768		# 32kb
hw.l2cachesize: 262144		# 256kb
hw.l3cachesize: 12582912		# 12mb
machdep.cpu.cache.linesize: 64
machdep.cpu.cache.size: 256
```
4Kx4K array of ints should be 16M * 4bytes = 64mb for the array; 16kb for a row
How many rows can we fit into the L1 cache? 32kb / 16kb = 2 at a time
How many rows can we fit into the L3 cache? 12mb / 16kb = 750 at a time

How many L3 misses might we expect given we’re traversing by columns? Well, since we’ll be traversing an entire row at a time, every 750 reads we’ll need to pull from memory. Since there are 4000 reads in a given column, that’s 4000/750 = 3.66 misses per column.

For the matrix, that’s 3.66 * 4000 columns = 14,640 misses, or about 0.0915% of the overall matrix. 

This is not what the stats are showing (0.8% overall LLd miss rate), suggesting I’m missing something:
```
(base) ResNet-59-66:memory-hierarchy-prework stephenzerfas$ valgrind --tool=cachegrind ./loop-order2
==98716== Cachegrind, a cache and branch-prediction profiler
==98716== Copyright (C) 2002-2017, and GNU GPL'd, by Nicholas Nethercote et al.
==98716== Using Valgrind-3.16.0.GIT and LibVEX; rerun with -h for copyright info
==98716== Command: ./loop-order2
==98716==
--98716-- warning: L3 cache found, using its data for the LL simulation.
--98716-- warning: specified LL cache: line_size 64  assoc 16  total_size 12,582,912
--98716-- warning: simulated LL cache: line_size 64  assoc 24  total_size 12,582,912
==98716==
==98716== I   refs:      230,926,422
==98716== I1  misses:          5,436
==98716== LLi misses:          2,897
==98716== I1  miss rate:        0.00%
==98716== LLi miss rate:        0.00%
==98716==
==98716== D   refs:      130,542,388  (97,733,283 rd   + 32,809,105 wr)
==98716== D1  misses:     16,023,911  (    21,911 rd   + 16,002,000 wr)
==98716== LLd misses:      1,008,764  (     7,471 rd   +  1,001,293 wr)
==98716== D1  miss rate:        12.3% (       0.0%     +       48.8%  )
==98716== LLd miss rate:         0.8% (       0.0%     +        3.1%  )
==98716==
==98716== LL refs:        16,029,347  (    27,347 rd   + 16,002,000 wr)
==98716== LL misses:       1,011,661  (    10,368 rd   +  1,001,293 wr)
==98716== LL miss rate:          0.3% (       0.0%     +        3.1%  )
```
