# Hyperloglog


Q: What's the max number of leading zeros in fnv32 hash output of any word in `/usr/share/dict/words`?
A: 17

Q: If you use the naive estimate of 2^n, where n is the max number of leading zeros in any hash output, what's the approximate cardinality of `/usr/share/dict/words`?
A: Naively we'd expect 2^17 = 131,072

Q: How does the answer change if you use different hash functions?
A: Skipped for now

Q: Can you find an adversarial set of words such that the naive estimate is WAY off?
A: skipped for now; would find the word with the max leading zeroes and just input that

Q: ``
