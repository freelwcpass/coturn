package main

import (
    "log"
    "net"
    "os"
    "strconv"
    "strings"

    "github.com/pion/turn/v2"
)

func env(key, def string) string {
    if v := os.Getenv(key); v != "" {
        return v
    }
    return def
}

func loadDotEnv(path string) {
    b, err := os.ReadFile(path)
    if err != nil {
        return
    }
    for _, line := range strings.Split(string(b), "\n") {
        s := strings.TrimSpace(line)
        if s == "" || strings.HasPrefix(s, "#") {
            continue
        }
        idx := strings.Index(s, "=")
        if idx <= 0 {
            continue
        }
        key := strings.TrimSpace(s[:idx])
        val := strings.TrimSpace(s[idx+1:])
        if len(val) >= 2 {
            if (val[0] == '"' && val[len(val)-1] == '"') || (val[0] == '\'' && val[len(val)-1] == '\'') {
                val = val[1 : len(val)-1]
            }
        }
        _ = os.Setenv(key, val)
    }
}

func mustParseIP(label, ipStr string) net.IP {
    if ipStr == "" {
        return nil
    }
    ip := net.ParseIP(ipStr)
    if ip == nil {
        log.Fatalf("%s is not a valid IP: %s", label, ipStr)
    }
    return ip
}

func main() {
    loadDotEnv(".env")
    listen := env("TURN_LISTEN_ADDR", ":3478")
    realm := env("TURN_REALM", "turn.example.com")
    staticUser := env("TURN_USERNAME", "turnuser")
    staticPass := env("TURN_PASSWORD", "turnpass")

    publicIP := env("RELAY_PUBLIC_IP", "")
    bindAddr := env("RELAY_BIND_ADDR", "0.0.0.0")
    minPortStr := env("RELAY_MIN_PORT", "49160")
    maxPortStr := env("RELAY_MAX_PORT", "49200")

    minPort, err := strconv.Atoi(minPortStr)
    if err != nil { log.Fatalf("invalid RELAY_MIN_PORT: %v", err) }
    maxPort, err := strconv.Atoi(maxPortStr)
    if err != nil { log.Fatalf("invalid RELAY_MAX_PORT: %v", err) }
    if minPort > maxPort { log.Fatalf("RELAY_MIN_PORT must be <= RELAY_MAX_PORT") }

    udpConn, err := net.ListenPacket("udp4", listen)
    if err != nil { log.Fatalf("failed to bind %s: %v", listen, err) }

    relayGen := &turn.RelayAddressGeneratorPortRange{
        RelayAddress: mustParseIP("RELAY_PUBLIC_IP", publicIP),
        Address:      bindAddr,
        MinPort:      uint16(minPort),
        MaxPort:      uint16(maxPort),
        MaxRetries:   100,
    }

    server, err := turn.NewServer(turn.ServerConfig{
        Realm: realm,
        AuthHandler: func(username, realm string, srcAddr net.Addr) ([]byte, bool) {
            if username != staticUser {
                return nil, false
            }
            return turn.GenerateAuthKey(staticUser, realm, staticPass), true
        },
        PacketConnConfigs: []turn.PacketConnConfig{{
            PacketConn:            udpConn,
            RelayAddressGenerator: relayGen,
        }},
    })
    if err != nil { log.Fatalf("failed to start Pion TURN server: %v", err) }
    _ = server
    log.Printf("Pion TURN listening on %s, realm=%s, relay range=%d-%d, publicIP=%s, user=%s", listen, realm, minPort, maxPort, publicIP, staticUser)
    select {}
}