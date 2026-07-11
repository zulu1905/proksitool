package ndisapi

import (
    "context"
    "fmt"
    ndisapi "myvpn/pkg/ndisapi/ndisapi-go"
)

func Run(ctx context.Context) {
    api, err := ndisapi.NewNdisApi()
    if err != nil {
        fmt.Println("Failed to create NDIS API:", err)
        return
    }
    defer api.Close()

    if !api.IsDriverLoaded() {
        fmt.Println("Windows Packet Filter driver is not installed")
        return
    }

    tcpAdapters, err := api.GetTcpipBoundAdaptersInfo()
    if err != nil {
        fmt.Println("GetTcpipBoundAdaptersInfo error:", err)
        return
    }
    if tcpAdapters == nil || tcpAdapters.AdapterCount == 0 {
        fmt.Println("No adapters found")
        return
    }

    handle := tcpAdapters.AdapterHandle[0]
    name := api.ConvertWindows2000AdapterName(string(tcpAdapters.AdapterNameList[0][:]))
    fmt.Println("Using adapter:", name)

    mode := &ndisapi.AdapterMode{
        AdapterHandle: handle,
        Flags:         0, // MSTCP_FLAG_* sabitlerini wrapper’da bulup buraya koy
    }

    if err := api.SetAdapterMode(mode); err != nil {
        fmt.Println("SetAdapterMode error:", err)
        return
    }

go func() {
        // EtherRequest'i oluştur (Sadece handle yeterli)
        req := &ndisapi.EtherRequest{
            AdapterHandle: handle,
        }

        for {
            select {
            case <-ctx.Done():
                return
            default:
                // Kütüphane req içini ReadPacket çağrıldığında dolduruyor
                if api.ReadPacket(req) {
                    // Paketi buradan işlemeye başla. 
                    // req.EthPacket veya req.pPacket (kütüphanendeki isim neyse) üzerinden verilere erişirsin.
                    
                    if err := api.SendPacketToAdapter(req); err != nil {
                        fmt.Println("SendPacketToAdapter error:", err)
                    }
                }
            }
        }
    }()
    }

