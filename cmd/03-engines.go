package main

import (
    "log"
    "os/exec"
    "syscall"

    A "github.com/wiresock/ndisapi-go"
    ndisapi "myvpn/pkg/ndisapi/ndisapi-go"
    "myvpn/pkg/spoofdpi"
    "myvpn/pkg/wireproxy"
    "myvpn/pkg/dnscrypt"
)

var (
    ndisRouter interface{}
    ndisApi    *A.NdisApi
)

func startService() {
    cleanAllEngines()

    // Mod seçimi
    switch currentCfg.ActiveMode {
    case "SpoofDPI":
         spoofdpi.Run(currentCfg.Addr, currentCfg.Port, currentCfg.SpoofArgs)
    case "Wireguard VPN":
         wireproxy.Run(currentCfg.Addr, currentCfg.Port)
    }

    // NDIS her zaman çalışır
    go func() {
        api, err := A.NewNdisApi()
        if err == nil {
            router, err := ndisapi.SetupProxyRouter(api)
            if err == nil {
                ndisApi = api
                ndisRouter = router
            }
        }
    }()

    // DNSCrypt her zaman çalışır
     dnscrypt.Run(currentCfg.DNSFilter)

    // WGCF (Warp config üretici) istenirse ayrı goroutine
     wgcf.Run()
}


func cleanAllEngines() {
    spoofdpi.Stop()
    wireproxy.Stop()
    dnscrypt.Stop()
    wgcf.Stop()

    if ndisRouter != nil {
        if r, ok := ndisRouter.(interface{ Close() }); ok {
            r.Close()
        }
        ndisRouter = nil
    }
    ndisApi = nil
}

func setSystemDNS(dnsMode string) error {
    var dnsServers string
    if dnsMode == "filtered" {
        dnsServers = "9.9.9.9,149.112.112.112"
    } else {
        dnsServers = "1.1.1.1,1.0.0.1"
    }

    psCommand := `Get-NetAdapter | Where-Object {$_.Status -eq 'Up'} | ForEach-Object { 
        Set-DnsClientServerAddress -InterfaceIndex $_.InterfaceIndex -ServerAddresses ("` + dnsServers + `") -ErrorAction SilentlyContinue 
    }; ipconfig /flushdns`

    cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-ExecutionPolicy", "Bypass", "-Command", psCommand)
    cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: CREATE_NO_WINDOW}
    return cmd.Run()
}

func restoreSystemDNS() error {
    psCommand := `Get-NetAdapter | Where-Object {$_.Status -eq 'Up'} | ForEach-Object { 
        Set-DnsClientServerAddress -InterfaceIndex $_.InterfaceIndex -ResetServerAddresses -ErrorAction SilentlyContinue 
    }; ipconfig /flushdns`

    cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-ExecutionPolicy", "Bypass", "-Command", psCommand)
    cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: CREATE_NO_WINDOW}
    return cmd.Run()
}
