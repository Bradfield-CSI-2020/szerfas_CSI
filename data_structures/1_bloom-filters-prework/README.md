# Bloom filters

### Description
Bloom filters are probabilistic data structures offering:
 * fast insertion and lookup times
 * the ability to check with a value has conclusively not already been inserted
 * the ability to check with some probability, whether or not a value has already been inserted.

### Sources
* Quick intro and illustration: https://llimllib.github.io/bloomfilter-tutorial/#footnote2
* In-depth walkthrough of math: http://www.michaelnielsen.org/ddi/why-bloom-filters-work-the-way-they-do/. The section "Summing up Bloom filters" is especially helpful for considering which parameters to use when implementing.

## My implementation

### Estimating parameters for our bloom filter:
Goal: Insert ~2Mb file `/usr/share/dict/words` with over ~230K words to the followign:
 * <100kb memory usage
 * 15% false positive rate
 * speed < few seconds

##### Approach (blended from above articles)
* Choose a desired max probability of a false positive, p
* Choose a ballpark value for n, number of items inserted into the bloom filter
* Choose a value for m, number of bits used in the underlying bit vector, where m = n/ln(2) * log(1/p)
* Calculate the optimal value of k, number of hash functions used, where k = ln(1/p).

_(note: I found the equations in the second article conflict with those in the first - the second made more sense to me so I used those)_

For 15% accuracy initial estimates come out to be:

	n			m			k				p
	235886		        280,386		        1.897119985		        0.15

We can't have a non-integer number of hash functions, so this is only approximate. From here I ramped up size to decrease false positive probabilities. 

##### Other questions to ask:
```
Which hash functions to use? 	        	A: Ideally fast, independent*, and uniform; fnv and murmur given as examples.
How many to use?				A: See estimates above. Started with two and made this work, could do more
Separate bit vector or combined?		A: Literature suggests using a combined bit vector is more memory efficient
What bit vector sizes work well?		A: See equations above
```
*Interestingly, and somewhat counterintuitively, [Kirtz and Mitzenmacher](http://www.eecs.harvard.edu/~michaelm/postscripts/rsa2008.pdf) suggest blending two independent and uniform hash functions to increase bloom filter performance without increasing false positive rates.
I implemented their approach adding fnv and murmer hashes before taking mod to fit into the hash set. I didn't go deep on the math, but the intuition makes sense to me. If I have hash values 12 and 13, then mod 10 of 12, 13, and (12+13) all return different results. 

##### Results
Tested different bitvector sizes with two hash functions, blending results to make a third hashed value, on 230K insertions
```
const BLOOMFILTER_SIZE = 280386 	// Elapsed time: 32.021804ms	Memory usage: 35048 bytes 	False positive rate: 32.35%
const BLOOMFILTER_SIZE = 300386 	// Elapsed time: 36.116905ms 	Memory usage: 37548 bytes 	False positive rate: 29.70%
const BLOOMFILTER_SIZE = 600000 	// Elapsed time: 31.764326ms 	Memory usage: 75000 bytes 	False positive rate: 10.56% <-- meets requirements
const BLOOMFILTER_SIZE = 800000 	// Elapsed time: 33.212294ms 	Memory usage: 100000 bytes 	False positive rate: 4.55%  <-- maxes size out at 100kb

```

#### Observations
I intially used just two hashes, then I implemented a blended but re-ran the first two hashes to create the blended (effectively doubling the amount of hashing). My final results (reflected above) store and use the first two hashes for best performance.
Implementing Kirtz and Mitzenmacher's blending of hashes makes a significant difference!
```
two hash calls, two values:         const BLOOMFILTER_SIZE = 600000 	// Elapsed time: 81.917925ms 	Memory usage: 75000 bytes 	False positive rate: 10.56% <-- meets requirements
four hash calls, three values:      const BLOOMFILTER_SIZE = 600000 	// Elapsed time: 124.68085ms 	Memory usage: 75000 bytes 	False positive rate: 8.88%
two hash calls, two values:         const BLOOMFILTER_SIZE = 600000     // Elapsed time: 32.909853ms    Memory usage: 75000 bytes   False positive rate: 8.88%
```

What I'd like to try next: additional blending - could I blend with a multiple to further decrease the false positive rate?

time estimate to complete: ~6.5 hours, 7.5 hours after refactor and optimization


##### When to use
In some ways, it's like a cache: 
* I can move something up the memory hierarchy
* If I want to ship something and don't want to ship the full data with it

Clearly a bloom filter can save you time, but it may also give you some privacy:
* Imagine there's a server knows about the transactions coming in. You want it to send you data when a transaction is relevant to you.
One way is to ask the server to send along data that matches certain criteria.
Downside to that is now the server knows exact details about you and relevant transactions you might not be want it to.
Another way to do it is to ask the server to hash criteria data and send transaction to you according to some bloom filter.
You and the server can agree on hash functions, then you send a bit vector and the server sends only transactions that collide with that bit vector after hashing.
The server will send you plenty of data on false positives, but far less than if sent you everything.
* Similarly, bloom filters to check for leaked password: Someone holding a password dump could receive a bit field from you and run it against the existing bloom filter of the dump.

 