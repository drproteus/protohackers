# 0. Smoke Test
> Deep inside Initrode Global's enterprise management framework lies a component that writes data to a server and expects to read the same data back. (Think of it as a kind of distributed system delay-line memory). We need you to write the server to echo the data back.
>
> Accept TCP connections.
>
> Whenever you receive data from a client, send it back unmodified.
>
> Make sure you don't mangle binary data, and that you can handle at least 5 simultaneous clients.
>
> Once the client has finished sending data to you it shuts down its sending side. Once you've reached end-of-file on your receiving side, and sent back all the data you've received, close the socket so that the client knows you've finished. (This point trips up a lot of proxy software, such as ngrok; if you're using a proxy and you can't work out why you're failing the check, try hosting your server in the cloud instead).
>
> Your program will implement the TCP Echo Service from RFC 862.



# Proof
```log
[Wed Jun 10 22:48:48 2026 UTC] [main.test] NOTE:check starts
[Wed Jun 10 22:48:48 2026 UTC] [main.test] NOTE:connected to 186.98.167.103 port 1337
[Wed Jun 10 22:48:48 2026 UTC] [main.test] NOTE:passed content.0
[Wed Jun 10 22:48:48 2026 UTC] [main.test] NOTE:connected to 186.98.167.103 port 1337
[Wed Jun 10 22:48:48 2026 UTC] [main.test] NOTE:passed content.1
[Wed Jun 10 22:48:48 2026 UTC] [main.test] NOTE:connected to 186.98.167.103 port 1337
[Wed Jun 10 22:48:49 2026 UTC] [main.test] NOTE:passed content.2
[Wed Jun 10 22:48:49 2026 UTC] [main.test] NOTE:connected to 186.98.167.103 port 1337
[Wed Jun 10 22:48:49 2026 UTC] [main.test] NOTE:passed content.3
[Wed Jun 10 22:48:49 2026 UTC] [main.test] PASS
[Wed Jun 10 22:48:50 2026 UTC] [multiclient.test] NOTE:check starts
[Wed Jun 10 22:48:50 2026 UTC] [multiclient.test] NOTE:connected to 186.98.167.103 port 1337
[Wed Jun 10 22:48:50 2026 UTC] [multiclient.test] NOTE:connected to 186.98.167.103 port 1337
[Wed Jun 10 22:48:50 2026 UTC] [multiclient.test] NOTE:connected to 186.98.167.103 port 1337
[Wed Jun 10 22:48:50 2026 UTC] [multiclient.test] NOTE:connected to 186.98.167.103 port 1337
[Wed Jun 10 22:48:50 2026 UTC] [multiclient.test] NOTE:connected to 186.98.167.103 port 1337
[Wed Jun 10 22:48:50 2026 UTC] [multiclient.test] NOTE:passed connection 0
[Wed Jun 10 22:48:50 2026 UTC] [multiclient.test] NOTE:passed connection 1
[Wed Jun 10 22:48:50 2026 UTC] [multiclient.test] NOTE:passed connection 2
[Wed Jun 10 22:48:50 2026 UTC] [multiclient.test] NOTE:passed connection 3
[Wed Jun 10 22:48:50 2026 UTC] [multiclient.test] NOTE:passed connection 4
[Wed Jun 10 22:48:50 2026 UTC] [multiclient.test] PASS
```
