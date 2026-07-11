package main

import (
    "encoding/json"
    "os"
    "path/filepath"
    "strings"
    "sync"
)

type AppConfig struct {
    Name     string `json:"name"`
    Protocol string `json:"protocol"`
    Mode     string `json:"mode"`
}

type Config struct {
    ActiveMode    string      `json:"active_mode"`
    Addr          string      `json:"addr"`
    Port          string      `json:"port"`
    DNS           string      `json:"dns"`
    WinSize       string      `json:"win_size"`
    SplitMode     string      `json:"split_mode"`
    BypassLan     bool        `json:"bypass_lan"`
    SpoofArgs     string      `json:"spoof_args"`
    Apps          []AppConfig `json:"apps"`
    DNSFilter     bool        `json:"dns_filter"`      // ✅ yeni alan
    LaunchBrowser bool        `json:"launch_browser"`  // ✅ yeni alan
}

type ProxifyreItem struct {
    AppNames            []string `json:"appNames"`
    Socks5ProxyEndpoint string   `json:"socks5ProxyEndpoint"`
    SupportedProtocols  []string `json:"supportedProtocols"`
}

type ProxifyreConfig struct {
    LogLevel  string          `json:"logLevel"`
    BypassLan bool            `json:"bypassLan"`
    Proxies   []ProxifyreItem `json:"proxies"`
    Excludes  []string        `json:"excludes"`
}

var defaultConfig = Config{
    ActiveMode:    "Wireguard VPN", // ✅ default mod
    DNS:           "filtered",
    Addr:          "127.0.0.1",
    Port:          "1233",
    BypassLan:     true,
    SpoofArgs:     "--app-mode socks5 --https-split-mode chunk --https-chunk-size 15",
    DNSFilter:     true,
    LaunchBrowser: false,
    Apps: []AppConfig{
        {Name: "Discord", Protocol: "TCP", Mode: "Tunnel"},
        {Name: "Webcord", Protocol: "TCP", Mode: "Tunnel"},
        {Name: "Roblox", Protocol: "TCP+UDP", Mode: "Tunnel"},
    },
}

var (
    configMu   sync.RWMutex
    currentCfg = defaultConfig
)

func resetToDefaults() {
    configMu.Lock()
    currentCfg = defaultConfig
    configMu.Unlock()
    saveSettings()
}

func loadSettings() {
    configMu.Lock()
    defer configMu.Unlock()
    confPath := filepath.Join(workDir, "config.json")
    if data, err := os.ReadFile(confPath); err == nil {
        _ = json.Unmarshal(data, &currentCfg)
    }
}

func saveSettings() {
    configMu.Lock()
    confData := currentCfg

    confPath := filepath.Join(workDir, "config.json")
    if newData, err := json.MarshalIndent(confData, "", "  "); err == nil {
        _ = os.WriteFile(confPath, newData, 0644)
    }

    var proxies []ProxifyreItem
    var excludes []string

    for _, a := range confData.Apps {
        if a.Mode == "Exclude" {
            excludes = append(excludes, a.Name)
        }
    }

    for _, p := range []string{"TCP", "UDP", "TCP+UDP"} {
        var names []string
        for _, a := range confData.Apps {
            if a.Mode == "Tunnel" && a.Protocol == p {
                names = append(names, a.Name)
            }
        }
        if len(names) > 0 {
            proxies = append(proxies, ProxifyreItem{
                AppNames:            names,
                Socks5ProxyEndpoint: confData.Addr + ":" + confData.Port,
                SupportedProtocols:  strings.Split(p, "+"),
            })
        }
    }

    pConfig := ProxifyreConfig{LogLevel: "error", BypassLan: confData.BypassLan, Proxies: proxies, Excludes: excludes}
    proxifyreConfPath := filepath.Join(workDir, "app-config.json")
    if pNewData, err := json.MarshalIndent(pConfig, "", "  "); err == nil {
        _ = os.WriteFile(proxifyreConfPath, pNewData, 0644)
    }
    configMu.Unlock()
}
