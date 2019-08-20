# Transaction ID Proof of Concept

```sh
git clone https://github.com/samlitowitz/transaction-id-poc.git
cd transaction-id-poc/build && docker-compose up
```

request -> router, add unique (for how long?) transaction id to request -> <resolved route>

echo route 



Add transaction id to every request
use the transaction id to track request throughout the system and compile log 

# TODO
1. Encode messages using protobuf
1. Make subscriber which aggregates messages by transaction id
