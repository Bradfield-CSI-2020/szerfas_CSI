# Merkle Hash Trees

## Notes
### [Hash functions](https://nakamoto.com/hash-functions/)
Terminology
* Pre-image - input
* Digest - output

Three characteristics of hash functions
* Deterministic
* Uniformly distributed
* Fixed output size

Three more characteristics of cryptographic hash functions
* One-way - computionally intractable to invert the hash and compute the pre-image
* Avalanche effect - change one bit and the whole hash changes. (One way to do this in implemtnation is rotate all the bits upon a change)
* Collision resistant - should be very difficult if not impossible to produce the same digest from different pre-images.

More on collision resistance
* One way to produce a collision is via pigeon-hole principle - if you're considering more inputs than possible outputs, then you'll have to have a collision somewhere.
For a hash with 256 characters, the digest set is 2^256. Given the number of atoms in the universe is close to 2^260, this is an enormous number, and all but statistically impossible we'd see a collision due to too many inputs.
* Another way to produce a collision is through analytical attempts or attacks. For example, rather than concerning ourselves with producing more outputs than a hash can handle, we can take a "birthday attack" approach. (See below)
* Producing random collisions in of itself is unlikely to have practical implications, but it's a sign that we're approaching a point in which a hash function might be cracked.
* Practically significant collisions are not random, but a chose prefix-collision, where given two leading prefixes, someone is able to determine two different messages that follow from the prefixes that will produce an identical digest.

Practical notes
* A chosen prefix collision was demonstrated for SHA-1 in 2019
* SHA-2 is still considered strong, and the default for many cyptographic applications

##### The birthday phenomenon and its implications for hash collisions
The birthday phenomenon goes like this: on average, how many guests need to be at a party so that two guests will share the same birthday? The answer is, surprisingly, 23. At 23 guests, there is a 50.7% probability that two people share a birthday. At 70 guests, it is 99.9% likely that two share the same birthday.

![Image of Birthday probability distribution](https://s31991.pcdn.co/wp-content/uploads/2019/02/Birthday_Paradox.png.webp)
Image credit: [Relatively Interesting](https://www.relativelyinteresting.com/the-birthday-paradox/?utm_source=org).

This is due to pair-wise relationship. With 23 people, how many possible pairs are there? This is 23-choose-2, or 23*22 / 2. What is the chance that a pair shares a birthday? 1/365. (The first person can have any birthday, the second person has a 1/365 of having that exact birthday.)

It's tempting to jump from here and conclude that with 253 pairs, we're at 253/365 or 69% chance some two people have the same birthday. 
But this would be a common probability flaw. (It might sharpen intuition by asking, would you expect 100% probability if we had 365 pairs? 
No, we'd still expect some chance that we'd not have a pair yet.) 
We really want to proceed as follows:
* The probability that two people have the same birthday is 1 - P(no two people have the same birthday)
* The probability a pair does not have the same birthday is therefore 364/365 = 99.7...%
* With 253 pairs, 99.7...% ^ 253 = 49.95...% probability that no two birthdays match.
* 1 - this number is 50.05%

It turns out that √n is a good approximation for the 50% likelihood threshold for a birthday collision (I've skipped over the math why).

This same logic can applied to finding collisions in hashes. So for a 128-bit hash output space, and because we only need √n time and √n space to generate a collision, we only need 2^64 iterations and space to generate a collision with 50% probability.


##### Observations on normal vs cryptographic hashing
* We don't use crypto hashes for things that don't need it because they're much slower
* That said, what if an attacker can anticipate your hashing function in a hash map and produce collisions?
They could intentionally insert lots of collisions and degrade performance.




### [Merkle Trees](https://nakamoto.com/introduction-to-cryptocurrency/)
What it is
* Hashing data upwards in a tree - bottom layer of blocks are hashed, then those digests are hashed, to create the next node, and so on

Benefits/use cases
* Checksum, if hash off, can drill down in log(N) time to see which part of the file wrong
* Sending reference to a unique combination of data without sending that data
* Inclusion proofs
![Image of merkle tree](https://i.imgur.com/IhX00ja.png)

Terms
* Cryptographic accumulator

In the block chain, a merkle root plays as follows:
![Image of merkle tree in blockchain](https://nakamoto.com/content/images/2019/12/image-13.png)


