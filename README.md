# Boost POW Golang library

This document describes the overall design of the project, which is not written yet. 

Boost POW is a library for embedding hash puzzles in Bitcoin script. The point is to draw the attention of Bitcoiners to information that is probably important. [Here]() is the whitepaper. 

## Protocol Design

Boost POW encodes a hash puzzle that is identical to that used to process Bitcoin blocks, which means the same infrastructure used for Bitcoin mining can be used for redeeming Boost outputs. The fields in a Boost script correspond to parts of a Bitcoin block header and coinbase and to fields in the Stratum protocol, which is used by hashers to communicate with Bitcoin nodes. 

Boost POW scripts have two variations, resulting in 4 total patterns. 

* **Contract vs Bounty:** Bounty scripts can be mined by anybody whereas contract scripts can only be mined by a specific miner holding the private key corresponding to a Bitcoin address. 
* **Version 1 vs Version 2:** Version 2 scripts take ASIC Boost into account, whereas version 1 does not. 

## Script Fields

Output script fields: 

* **Content:** corresponds to the previous block hash in a Bitcoin block header. This contains the hash of the content being boosted. 
* **Target:** corresponds to the target or nBits field in the Bitcoin block header. This is the difficulty of the boost. 
* **category:** corresponds to version field in Bitcoin block header. In Boost version 1, this has 4 bytes available whereas in version 2, only 2 bytes are available. 
* **topic/tag:** corresponds to the first part of the coinbase. Up to 20 bytes. 
* **additional data:** corresponds to the last part of the coinbase. Can be anything, no size limit. 
* **user_nonce:** A number that should be filled in randomly. Used to ensure that two identical boost scripts are not accidentally created. 

Input or output script fields: 

* **miner pubkey hash:** Corresponds to the first part of the coinbase. Used to ensure that no Boost input script can be re-used by someone else to redeem an output. In a bounty script, this field is in the input script. In a contract script, it is in the output script. 

Input script field:

* **Time:** the time that the boost was completed, same as the time field in the Bitcoin block header. 
* **nonce:** same as in the Bitcoin block header. 
* **extra_nonce 1 and 2:** Corresponds to the middle of the coinbase, and is a thing in Stratum. Extra nonce 1 is assigned by the mining pool to the miners. Extra nonce 2 is chosen by the miner. 
* **miner pubkey:** used to verify the signature and is not part of the hash puzzle. 
* **signature:** a signature, also not part of the hash puzzle. 
* **general purpose bits:** These only appear in a version 2 script and have to do with ASIC Boost. 

## Antecedents

* [**go-work**](https://github.com/DanielKrawisz/go-work): A golang libary for hash puzzles. 
* [**go-Stratum**](https://github.com/DanielKrawisz/go-Stratum): An incomplete golang libary for Stratum messages. 
* [**Gigamonkey**](): A c++ library that implements Boost and Stratum with extensive tests. 
* [**boostpow-js**](): A JavaScript library that implements Boost POW. In particular, there are unit tests that can be traslated here. 

## Design of the Library

There will be types for a Boost output and input script. The following operations should be supported, in order of importance: 
* Boost types <=> Bitcoin scripts
* Boost types => hash puzzles in [go-work](https://github.com/DanielKrawisz/go-work). 
* {miner key, hash puzzle solution} => Boost input
* Boost output => Stratum notify message
* {Stratum subscribe response, Stratum submit request} => hash puzzle solution
