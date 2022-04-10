# Distributed-Leaderless-Counter
Gossip protocol based counter in a leaderless cluster with customizable chatter size and query aging.

# Write Operations
Users will send inc(key, val) operations to any random node. Each node (single replica for simplicity) maintains it's own local store.

# Read Operations
User will query(key) a random node and a go-routine will display the query results periodically until the user is satisfied or not-enough changes are seen.

# Internals

Tracker -
1. Maintains atomic store of nodes with name as primary key.
2. Publishes node list as heartbeat to all nodes every 1 second (daemon).

Node - 
1. Maintains atomic store of key-value pairs.
2. A single process-query loop daemon to process incoming queries (daemon).
3. Based on query visit frequency, can either discard query or process query.
4. Query process simply updates query with info and sends query to random un-visited chatter-size nodes.

Query -
1. Maintains atomic store of visited node-response pairs.
2. A status daemon to display node-reponse periodically.

# TODO
1. Implement convergence detection and end query daemon. (priority)
2. All configurable replicas to each node and write ops (sync + async) to subset replicas. (replicas will be separate repo)
3. Configurable process-query loop count for each node. (boring)