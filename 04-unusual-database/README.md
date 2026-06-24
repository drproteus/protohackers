# 4. Unusual Database Program

https://protohackers.com/problem/4

# Proof
```
[Wed Jun 24 18:20:39 2026 UTC] [0simple.test] NOTE:check starts
[Wed Jun 24 18:20:39 2026 UTC] [0simple.test] PASS
[Wed Jun 24 18:20:40 2026 UTC] [1users.test] NOTE:check starts
[Wed Jun 24 18:20:40 2026 UTC] [1users.test] NOTE:inserting 20 usernames and checking them
[Wed Jun 24 18:20:44 2026 UTC] [1users.test] PASS
[Wed Jun 24 18:20:45 2026 UTC] [2products.test] NOTE:check starts
[Wed Jun 24 18:20:45 2026 UTC] [2products.test] NOTE:inserting 500 JSON product records
[Wed Jun 24 18:20:55 2026 UTC] [2products.test] NOTE:checking 100 JSON product records selected at random
[Wed Jun 24 18:21:00 2026 UTC] [2products.test] NOTE:no value for 'product.85494' (response packet dropped?)
[Wed Jun 24 18:21:00 2026 UTC] [2products.test] NOTE:no value for 'product.85456' (response packet dropped?)
[Wed Jun 24 18:21:00 2026 UTC] [2products.test] NOTE:no value for 'product.85200' (response packet dropped?)
[Wed Jun 24 18:21:00 2026 UTC] [2products.test] NOTE:no value for 'product.85139' (response packet dropped?)
[Wed Jun 24 18:21:00 2026 UTC] [2products.test] NOTE:no value for 'product.85187' (response packet dropped?)
[Wed Jun 24 18:21:00 2026 UTC] [2products.test] PASS
[Wed Jun 24 18:21:01 2026 UTC] [3version.test] NOTE:check starts
[Wed Jun 24 18:21:01 2026 UTC] [3version.test] NOTE:checking 'version' implementation
[Wed Jun 24 18:21:02 2026 UTC] [3version.test] NOTE:got server version 'Ken's Key-Value Store 1.0'
[Wed Jun 24 18:21:02 2026 UTC] [3version.test] PASS
[Wed Jun 24 18:21:03 2026 UTC] [4edgecases.test] NOTE:check starts
[Wed Jun 24 18:21:03 2026 UTC] [4edgecases.test] NOTE:checking edge cases
[Wed Jun 24 18:21:05 2026 UTC] [4edgecases.test] NOTE:inserting value for SaneFrank112, and waiting for it to update
[Wed Jun 24 18:21:05 2026 UTC] [4edgecases.test] NOTE:overwriting value for SaneFrank112, and waiting for it to update
[Wed Jun 24 18:21:05 2026 UTC] [4edgecases.test] PASS
```
