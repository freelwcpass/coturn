# Pion TURN Server

Minimal standalone TURN/STUN server using `github.com/pion/turn` with ephemeral auth compatible with LiveKit.

## Env Vars

- `TURN_LISTEN_ADDR`: UDP listen address (default `:3478`)
- `TURN_REALM`: TURN realm (default `turn.example.com`)
- `TURN_USERNAME`: Static TURN username
- `TURN_PASSWORD`: Static TURN password
- `RELAY_PUBLIC_IP`: Public IP advertised in relay candidates
- `RELAY_BIND_ADDR`: Local bind address for relay sockets (default `0.0.0.0`)
- `RELAY_MIN_PORT` / `RELAY_MAX_PORT`: Relay UDP port range (default `49160–49200`)

## Build

```
docker build -t pion-turn:latest .
```

## Run

```
docker run -d --name pion-turn \
  -p 3478:3478/udp -p 49160-49200:49160-49200/udp \
  pion-turn:latest
```

The server automatically loads configuration from `.env` in the working directory. Edit `turn-server/.env` before building to customize values.

### Run with explicit credentials

```
docker run -d --name pion-turn \
  -p 3478:3478/udp -p 49160-49200:49160-49200/udp \
  -e TURN_REALM=turn.example.com \
  -e TURN_USERNAME=turnuser \
  -e TURN_PASSWORD=turnpass \
  -e RELAY_PUBLIC_IP=203.0.113.10 \
  pion-turn:latest
```

## LiveKit Integration (static credentials)

```
rtc:
  turn_servers:
    - host: turn.example.com
      port: 3478
      protocol: udp
      username: turnuser
      credential: turnpass
```

Ensure DNS for `turn.example.com` points to the host/node running this server. Open `3478/udp` and the relay UDP range on your firewall.
