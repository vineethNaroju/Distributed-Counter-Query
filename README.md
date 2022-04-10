# Distributed-Leaderless-Counter
Gossip protocol based counter in a leaderless cluster with customizable chatter size and query aging.

# Write Operations
Users will send inc(key, val) operations to any random node. Each node (single replica for simplicity) maintains it's own local store.

# Read Operations
User will query(key) a random node and a go-routine will display the query results periodically until the user is satisfied or not-enough changes are seen.

# Internals
Node gossips with random nodes for query results and we pass query pointer object (for db simulation) - to perform atomic operations.
Gossip queries die away due to aging (or we can signal nodes to ignore / kill query).

## Tracker
1. Maintains atomic store of nodes with name as primary key.
2. Publishes node list as heartbeat to all nodes every 1 second (daemon).

## Node
1. Maintains atomic store of key-value pairs.
2. A single process-query loop daemon to process incoming queries (daemon).
3. Based on query visit frequency, can either discard query or process query.
4. Query process simply updates query with info and sends query to random un-visited chatter-size nodes.

## Query
1. Maintains atomic store of visited node-response pairs.
2. A status daemon to display node-reponse periodically.

# Results
1. Tested with 1000 node cluster with perpetual requests (inc key)every 1 millisecond to random node.
2. Tested for a single query (get above key)
3. See loadtest_response.log simulated with 10ms delay of network latency in same datacenter.

# TODO
1. Implement convergence detection and end query daemon (priority).
2. Loadtest under heavy writes and decent read requests (priority).
3. All configurable replicas to each node and write ops (sync + async) to subset replicas (later).
4. Configurable process-query loop count for each node (boring).
