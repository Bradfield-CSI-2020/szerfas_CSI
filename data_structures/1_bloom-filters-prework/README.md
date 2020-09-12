# Bloom filters

### Description
Bloom filters are probabilistic data structures offering fast insertion and lookup times, and the ability to check with a value has conclusively not already been inserted and, with some probability, whether or not a value has alraedy been inserted.

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

For 15% accuracy initial estimates come out to be

	n			m			k				p
	235886		        280,386		        1.897119985		        0.15

We can't have a non-integer number of hash functions, so this is only approximate

##### Other questions to ask:
```
Which hash functions to use? 	        	A: Ideally fast and independent; fnv and murmur given as examples
How many to use?				A: See estimates above. Started with two and made this work, could do more
Separate bit vector or combined?		A: Literature suggests using a combined bit vector is more memory efficient
What bit vector sizes work well?		A: See equations above
```

##### Results
Tested different bitvector sizes with two hash functions on 230K insertions
```
const BLOOMFILTER_SIZE = 280386 	// Elapsed time: 85.720142ms	Memory usage: 35048 bytes 	False positive rate: 32.35%
const BLOOMFILTER_SIZE = 300386 	// Elapsed time: 89.715896ms 	Memory usage: 37548 bytes 	False positive rate: 29.70%
const BLOOMFILTER_SIZE = 600000 	// Elapsed time: 81.917925ms 	Memory usage: 75000 bytes 	False positive rate: 10.56% <-- meets requirements
```

#### Observations
Adding a third, blended hash function, substantially increased the time to run but further decreasedb the false positive rate, highlighting one of the tradeoffs of using more hashes.
```
const BLOOMFILTER_SIZE = 600000 	// Elapsed time: 124.68085ms 	Memory usage: 75000 bytes 	False positive rate: 8.88%
```

time estimate to complete: ~6.5 hours