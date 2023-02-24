# Boost POW Golang library

This document describes the overall design of the project, which is not written yet. 

Boost POW is a library for embedding hash puzzles in Bitcoin script. The point is to draw the attention of Bitcoiners to information that is probably important. [Here](https://bitcoinfiles.org/t/7332808b5283f8acedcc6240a42f669cc3d305413201527852061fd5b283d0d8) is the whitepaper. 

Some other articles describing the economics and rationale for Boost POW: 

* [Proof-of-Work as a Handicap](https://bitcoinfiles.org/t/0c9544cf8650794d0221a0b11fec45ed19409e6deef9b3eeeea7ee956cdde7af)
* [Proof-of-work as an Upvote System](https://bitcoinfiles.org/t/f9e6c4f0ac7219257e1276cd23c1bff5e5088204ff4e3471786c6252fb00f01e)
* [Pow.co FAQ](https://github.com/DanielKrawisz/Entropy/blob/main/Pow.co%20FAQ.md)

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
* [**Gigamonkey**](https://github.com/Gigamonkey-BSV/Gigamonkey): A c++ library that implements Boost and Stratum with extensive tests. 
* [**boostpow-js**](https://github.com/pow-co/boostpow-js/): A JavaScript library that implements Boost POW. In particular, there are unit tests that can be traslated here. 

## Design of the Library

There will be types for a Boost output and input script. The following operations should be supported, in order of importance: 
* Boost types <=> Bitcoin scripts
* Boost types => hash puzzles in [go-work](https://github.com/DanielKrawisz/go-work). 
* {miner key, hash puzzle solution} => Boost input
* Boost output => Stratum notify message
* {Stratum subscribe response, Stratum submit request} => hash puzzle solution

## Test vectors 

We boost the Boost whitepaper, which has txid `0x7332808b5283f8acedcc6240a42f669cc3d305413201527852061fd5b283d0d8`. There are four variations, corresponding to version 1 vs. version 2 and contract vs bounty. We have output scripts for all variations as well as example spending and redeeming txs. 

### Bounty, version 1

output script hex: 

```
08626f6f7374706f7775045704000020d8d083b2d51f0652785201324105d3c39c662fa44062ccedacf883528b803273049cff631d067468656f727904890000001c746869732069732074686520426f6f737420776869746570617065727e7c557a766b7e52796b557a8254887e557a8258887e7c7eaa7c6b7e7e7c8254887e6c7e7c8254887eaa01007e816c825488537f7681530121a5696b768100a0691d00000000000000000000000000000000000000000000000000000000007e6c539458959901007e819f6976a96c88ac
```

expected values: 

```
{
  content=0x7332808b5283f8acedcc6240a42f669cc3d305413201527852061fd5b283d0d8 
  difficulty: .01 
  topic: "theory" 
  additional_data: "this is the Boost whitepaper"
  category: 1111
  user_nonce: 137 
  version: 1
}
```

spending transaction: `0xc5c7248302683107aa91014fd955908a7c572296e803512e497ddf7d1f458bd3`

redeeming transaction: `0x99bbbc28d39427bf530c05dc90db12e0953122fe5055afce6370d89a0085c28d`

redeeming script values: 

```
{
  signature: 3044022100ac4003d62ddadbf0bff9cbe63d0f6ad740494ee7fcf5f296cfc056f52f087c7c021f2f9e2db03b141ce88edc1c10850a0831dea63edd6c6a8040d80e24737e6d4a41  
  pubkey: 03097e9768554d40c0b5b18e44db2a15bbd137a373c39af46033049477bcbb79a4 
  nonce : 31497 
  timestamp : 1677268580
  extra_nonce_2 : "0f0445b186e64adc"
  extra_nonce_1 : 909479219 
  miner_pubkey_hash: 0xa3c10ac097a7da0009a786cc17edc1391a3bddf6
}
```

### Bounty, version 2

output script hex: 

```
08626f6f7374706f7775045704000020d8d083b2d51f0652785201324105d3c39c662fa44062ccedacf883528b803273049cff631d067468656f727904890000001c746869732069732074686520426f6f737420776869746570617065727e7c557a766b7e52796b567a8254887e567a820120a1697e7c7eaa7c6b7e6b04ff1f00e076836b847c6c84856c7e7c8254887e6c7e7c8254887eaa01007e816c825488537f7681530121a5696b768100a0691d00000000000000000000000000000000000000000000000000000000007e6c539458959901007e819f6976a96c88ac
```

expected values: 

```
{
  content=0x7332808b5283f8acedcc6240a42f669cc3d305413201527852061fd5b283d0d8 
  difficulty: .01 
  topic: "theory" 
  additional_data: "this is the Boost whitepaper"
  category: 1111
  user_nonce: 137 
  version: 2
}
```

spending transaction: `0x12aaef887e8348e83eac2937849de22a0bda8f7c2c819199bbcbb20b01722144`

redeeming transaction: `0x0f38e43dfc603296ef6883da389fc93815c0535bfba255b070a98bd6cc4da984`

redeeming script values: 

```
{
  signature: 304502210081cac0bdfb713e8c6632ec8c7b6f1d070b19a43c3b06e05174f25dc9065c6e910220787dd9d0f58f79cda8b7f5b436eb2f8cd6d50dc5271e6216308c286406d4166141  
  pubkey: 03e0fd48907c0117600a6326aafe7d43adbc9421a4381bb6579f1ab4912cd25e37 
  nonce : 5267719 
  timestamp : 1677269436 
  extra_nonce_2 : "b4d8e1f74255bebc" 
  extra_nonce_1 : 2329617541 
  miner_pubkey_hash: 0x81bb8505a9999135a105e2f0290d55b1b70f7d3f
}
```

### Contract, version 1

output script hex: 

```
08626f6f7374706f7775143f7d0fb7b1550d29f0e205a1359199a90585bb81045704000020d8d083b2d51f0652785201324105d3c39c662fa44062ccedacf883528b803273049cff631d067468656f727904890000001c746869732069732074686520426f6f737420776869746570617065727e7c557a766b7e52796b557a8254887e557a8258887e7c7eaa7c6b7e7e7c8254887e6c7e7c8254887eaa01007e816c825488537f7681530121a5696b768100a0691d00000000000000000000000000000000000000000000000000000000007e6c539458959901007e819f6976a96c88ac
```

expected values: 

```
{
  content=0x7332808b5283f8acedcc6240a42f669cc3d305413201527852061fd5b283d0d8 
  difficulty: .01 
  topic: "theory" 
  additional_data: "this is the Boost whitepaper"
  category: 1111
  user_nonce: 137 
  miner_address: 16nhPWCkbkR1bNACwPYULBWyvxQ5MCDZBo
  version: 1
}
```

spending transaction: `0xed122aa475c02ee049b342d9224bc140f015eee30b8411ad999c6a8378d9766e`

redeeming transaction: `0x85f9461ce4c88755052673edc6aab16d817b80aad6ce02ad2e80d36e6df78317`

redeeming script values: 

```
{
  signature: 3045022100ac1f6aa4153037920cec3d18ef6e129eaefabe88c98f8e19aea7af806b645aad02203b521118c2b7b5fa64382f559db17a5d58e0d160196d3f06989df98731215f1b41
  pubkey: 03e0fd48907c0117600a6326aafe7d43adbc9421a4381bb6579f1ab4912cd25e37 
  nonce : 4799381
  timestamp : 1677270497
  extra_nonce_2 : "2df60742aed5d329"
  extra_nonce_1 : 1035367878
}
```

### Contract, version 2

output script hex: 

```
08626f6f7374706f7775143f7d0fb7b1550d29f0e205a1359199a90585bb81045704000020d8d083b2d51f0652785201324105d3c39c662fa44062ccedacf883528b803273049cff631d067468656f727904890000001c746869732069732074686520426f6f737420776869746570617065727e7c557a766b7e52796b567a8254887e567a820120a1697e7c7eaa7c6b7e6b04ff1f00e076836b847c6c84856c7e7c8254887e6c7e7c8254887eaa01007e816c825488537f7681530121a5696b768100a0691d00000000000000000000000000000000000000000000000000000000007e6c539458959901007e819f6976a96c88ac
```

expected values: 

```
{
  content=0x7332808b5283f8acedcc6240a42f669cc3d305413201527852061fd5b283d0d8 
  difficulty: .01 
  topic: "theory" 
  additional_data: "this is the Boost whitepaper"
  category: 1111
  user_nonce: 137 
  miner_address: 16nhPWCkbkR1bNACwPYULBWyvxQ5MCDZBo
  version: 2
}
```

spending transaction: `0x6a26e314d33cfb0948e9bb12559abc9c687403c639e4302faf562d818b2ff0a2`

redeeming transaction: `0xa512c846b0154d23325f40ef87a088d747252fc2a179cca067a6026ee59c5ea6`

redeeming script values: 

```
{
  signature: 304402201f94a12ace389cd389ef129dc9b68eb1a357ff6f71a508aa0b3accd90736007702206d316fce43e5ae24a6b07acc342e0f7a5c0d0366a2a00dee00acbb25b8f4f6a941
  pubkey: 03e0fd48907c0117600a6326aafe7d43adbc9421a4381bb6579f1ab4912cd25e37 
  nonce : 3901135
  timestamp : 1677271659
  extra_nonce_2 : "4f22be6e277ead90"
  extra_nonce_1 : 3783406472
}
```
