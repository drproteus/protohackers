# 2. Means to an End

https://protohackers.com/problem/2

# Proof
```
[Thu Jun 25 17:38:36 2026 UTC] [0example.test] NOTE:check starts
[Thu Jun 25 17:38:36 2026 UTC] [0example.test] NOTE:connected to 186.98.167.103 port 20400
[Thu Jun 25 17:38:36 2026 UTC] [0example.test] NOTE:sending example data
[Thu Jun 25 17:38:37 2026 UTC] [0example.test] PASS
[Thu Jun 25 17:38:38 2026 UTC] [1main.test] NOTE:check starts
[Thu Jun 25 17:38:38 2026 UTC] [1main.test] NOTE:4 simultaneous clients
[Thu Jun 25 17:38:38 2026 UTC] [1main.test] NOTE:connected to 186.98.167.103 port 20400
[Thu Jun 25 17:38:38 2026 UTC] [1main.test] NOTE:connected to 186.98.167.103 port 20400
[Thu Jun 25 17:38:39 2026 UTC] [1main.test] NOTE:connected to 186.98.167.103 port 20400
[Thu Jun 25 17:38:39 2026 UTC] [1main.test] NOTE:connected to 186.98.167.103 port 20400
[Thu Jun 25 17:38:47 2026 UTC] [1main.test] PASS
[Thu Jun 25 17:38:49 2026 UTC] [2largedata.test] NOTE:check starts
[Thu Jun 25 17:38:49 2026 UTC] [2largedata.test] NOTE:inserting 200k prices in random order
[Thu Jun 25 17:38:49 2026 UTC] [2largedata.test] NOTE:connected to 186.98.167.103 port 20400
[Thu Jun 25 17:38:52 2026 UTC] [2largedata.test] PASS
[Thu Jun 25 17:38:54 2026 UTC] [3badclient.test] NOTE:check starts
[Thu Jun 25 17:38:54 2026 UTC] [3badclient.test] NOTE:checking whether bad clients can disrupt good clients
[Thu Jun 25 17:38:54 2026 UTC] [3badclient.test] NOTE:connected to 186.98.167.103 port 20400
[Thu Jun 25 17:38:54 2026 UTC] [3badclient.test] NOTE:connected to 186.98.167.103 port 20400
[Thu Jun 25 17:38:54 2026 UTC] [3badclient.test] NOTE:connected to 186.98.167.103 port 20400
[Thu Jun 25 17:38:54 2026 UTC] [3badclient.test] NOTE:connected to 186.98.167.103 port 20400
[Thu Jun 25 17:38:54 2026 UTC] [3badclient.test] NOTE:sending an incomplete message
[Thu Jun 25 17:38:54 2026 UTC] [3badclient.test] NOTE:sending an illegal message type
[Thu Jun 25 17:38:54 2026 UTC] [3badclient.test] NOTE:disconnecting immediately after sending a query
[Thu Jun 25 17:38:55 2026 UTC] [3badclient.test] PASS
[Thu Jun 25 17:38:56 2026 UTC] [4intoverflow.test] NOTE:check starts
[Thu Jun 25 17:38:56 2026 UTC] [4intoverflow.test] NOTE:testing for integer overflow
[Thu Jun 25 17:38:56 2026 UTC] [4intoverflow.test] NOTE:connected to 186.98.167.103 port 20400
[Thu Jun 25 17:38:56 2026 UTC] [4intoverflow.test] PASS
```
