# huus

Huus engine is a persistent, b+ tree based, kv-store written in golang.
It is a project meant to demonstrate b plus trees data structures 
using modern implementations and optimizations using go routines.

### History

B trees were first developed so that memory like hard disks could
easily handle search queries really quickly without needing to read from 
multiple different places, since the architecture of a b tree allows storage
in a block type format.

B+ trees is a variation of B trees that allows easier range queries by
storing all the values only in the leaf nodes and linking all the nodes
to form a chain.


