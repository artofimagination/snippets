# Turn server POC
This snippet is created to test the configuration of a turn server using [coturn](https://github.com/coturn/coturn)
The test constists of a turn server, a client and a server written in Go using [pion-to-pion webrtc](https://github.com/pion/webrtc) example code

# How to use?
The turn server is configured to use:
  - static long term credentials
  - tls
  - pre set realm name
  - run on localhost

## Set tls certificates  
By default all certificates are pregenerated an uploaded into the repo. Feel free to regenerate them by running generateCerts.sh
To configure the tls extensions modify
  - [extcnfCA.ext](https://github.com/artofimagination/snippets/blob/master/turn-server/extcnfCA.ext) for the CA
  - [extcnf.ext](https://github.com/artofimagination/snippets/blob/master/turn-server/extcnf.ext) for all the other certificates
  
## Start system
  - `docker-compose up -d --force-recreate webrtc-client`

## Set user
Got the syetem working by setting the plain password in the webrtc clients and using the encrupted password in the `turnserver.conf`.
To get the encryption run:
  - ```docker exec -it turn-server bash -c "turnadmin -k -u username -r realm_name -p passwd"```
  - copy the generated key into turnserver.conf under "'Static' user accounts" section like the example states there

Both webrtc-client and webrtc-server has their own user registered (testUser and testUser2)
  
# Known issues
At the moment the webrtc datachannel works only if the ICE Server is configured with `ICETransportPolicyAll`.
Using `ICETransportPolicyRelay` will still try to establish relaying trhough the turn server, but it will fail at the end. The closes hint on the issue is the following corutn log

```0: session 008000000000000001: realm <pocturn.localhost> user <defaultUser>: incoming packet message processed, error 401: Unauthorized
0: write_client_connection:4292:start
0: write_client_connection: prepare to write to s 0x7f34ac028b20
0: write_client_connection:4315:end
0: read_client_connection:4614:end
0: udp_server_input_handler:628:start
0: ssl_read: before read...
0: ssl_read: after read: 0
0: shutdown_client_connection:4172:start
0: session 008000000000000001: usage: realm=<pocturn.localhost>, username=<defaultUser>, rp=4, rb=478, sp=2, sb=248
0: session 008000000000000001: peer usage: realm=<pocturn.localhost>, username=<defaultUser>, rp=0, rb=0, sp=0, sb=0
0: closing session 0x7f34ac028d20, client socket 0x7f34ac028b20 (socket session=0x7f34ac028d20)
0: session 008000000000000001: closed (2nd stage), user <defaultUser> realm <pocturn.localhost> origin <>, local 172.18.0.3:5349, remote 172.18.0.4:39682, reason: SSL read error
0: session 008000000000000001: SSL shutdown received, socket to be closed (local 172.18.0.3:5349, remote 172.18.0.4:39682)
0: shutdown_client_connection:4237:end
0: udp_server_input_handler:666:end
0: udp_server_input_handler:628:start
0: dtls_server_input_handler:268:start
0: dtls_accept_client_connection:237:start
```
```
