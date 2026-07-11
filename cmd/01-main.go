package main

import (
    "fmt"
    "os"
    "context"
    "time"
    "path/filepath"
    "github.com/getlantern/systray"
    "fyne.io/fyne/v2/app"
)

var myApp = app.New()
var (
    mStatus  *systray.MenuItem
    mInfo    *systray.MenuItem
    mDns     *systray.MenuItem
    mDnsFilt *systray.MenuItem
    mDnsUn   *systray.MenuItem
    mModes   *systray.MenuItem
    mSpoof   *systray.MenuItem
    mWG      *systray.MenuItem
)

func main() {
    _, cancel := context.WithCancel(context.Background())
    defer cancel()

    confPath := filepath.Join(workDir, "config.json")
    if _, err := os.Stat(confPath); os.IsNotExist(err) {
        saveSettings()
    }

    loadSettings()
    initGUIPindows()
    go systray.Run(onReady, onExit)
    myApp.Run()
}

func updateStatus() {
    dnsDurum := "Filtreli"
    if currentCfg.DNS == "unfiltered" {
        dnsDurum = "Filtresiz"
    }
    mStatus.SetTitle(fmt.Sprintf("🟢 MOD: %s | DNS: %s", currentCfg.ActiveMode, dnsDurum))
    mInfo.SetTitle(fmt.Sprintf("📍 ADRES: %s:%s", currentCfg.Addr, currentCfg.Port))
}

func onReady() {
    systray.SetTooltip("ProksiTool | Proxy & DPI bypass aracı")
    mStatus = systray.AddMenuItem("", "")
    mStatus.Disable()
    mInfo = systray.AddMenuItem("", "")
    mInfo.Disable()
    updateStatus()
    systray.AddSeparator()

    // DNS toggle
    mDns = systray.AddMenuItem("🛡️ DNS Mod Seç", "DNS sağlayıcısını belirle")
    mDnsFilt = mDns.AddSubMenuItemCheckbox("Filtreli (Güvenli)", "", currentCfg.DNS == "filtered")
    mDnsUn   = mDns.AddSubMenuItemCheckbox("Filtresiz (Ham)", "", currentCfg.DNS == "unfiltered")

    // Proxy mod seçimi (SpoofDPI ↔ WireProxy)
    mModes = systray.AddMenuItem("🔄 Proxy Mod Seç", "Bağlantı yöntemini belirle")
    mSpoof = mModes.AddSubMenuItemCheckbox("SpoofDPI", "", currentCfg.ActiveMode == "SpoofDPI")
    mWG    = mModes.AddSubMenuItemCheckbox("WireProxy", "", currentCfg.ActiveMode == "WireProxy")

    systray.AddSeparator()
    mQuit := systray.AddMenuItem("❌ Çıkış", "Uygulamayı kapatır")

    go func() {
        for {
            select {
            case <-mDnsFilt.ClickedCh:
                go applyDnsMode("filtered", mDnsFilt, mDnsUn)
            case <-mDnsUn.ClickedCh:
                go applyDnsMode("unfiltered", mDnsFilt, mDnsUn)
            case <-mSpoof.ClickedCh:
                go applyNewMode("SpoofDPI")
            case <-mWG.ClickedCh:
                go applyNewMode("WireProxy")
            case <-mQuit.ClickedCh:
                onExit()
            }
        }
    }()
}

func applyDnsMode(mode string, filtItem, unItem *systray.MenuItem) {
    currentCfg.DNS = mode
    saveSettings()
    updateStatus()
    cleanAllEngines()
    startService() // DNSCrypt yeniden başlatılır

    if mode == "filtered" {
        filtItem.Check()
        unItem.Uncheck()
    } else {
        filtItem.Uncheck()
        unItem.Check()
    }
}

func applyNewMode(mode string) {
    currentCfg.ActiveMode = mode
    saveSettings()
    updateStatus()
    cleanAllEngines()
    startService()

    // Menüde seçilen mode’u işaretle
    mSpoof.Uncheck()
    mWG.Uncheck()
    switch mode {
    case "SpoofDPI":
        mSpoof.Check()
    case "WireProxy":
        mWG.Check()
    }
}

func onExit() {
    cleanAllEngines()
    if err := restoreSystemDNS(); err != nil {
        fmt.Println("DNS geri yüklenirken hata:", err)
    }
    time.Sleep(500 * time.Millisecond)
    os.Exit(0)
}
