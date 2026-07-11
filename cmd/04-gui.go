package main

import (
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/layout"
    "fyne.io/fyne/v2/widget"
    "fyne.io/fyne/v2/canvas"
    "image/color"
    "os/exec"
)

var (
    settingsWindow     fyne.Window
    appListWindow      fyne.Window
    appListContainer   *fyne.Container
    refreshAppListFunc func()
)

func initGUIPindows() {
    settingsWindow = myApp.NewWindow("Ayarlar")
    settingsWindow.Resize(fyne.NewSize(600, 480))
    settingsWindow.CenterOnScreen()
    settingsWindow.SetCloseIntercept(func() { settingsWindow.Hide() })

    // Proxy motoru seçimi (WireProxy / SpoofDPI)
    engineSelect := widget.NewSelect([]string{"WireProxy", "SpoofDPI"}, func(s string) {
        currentCfg.ActiveMode = s
        saveSettings()
        go startService()
    })
    engineSelect.SetSelected(currentCfg.ActiveMode)

    // DNS filtre toggle
    dnsFilterCheck := widget.NewCheck("DNS filtre aktif", func(b bool) {
        currentCfg.DNSFilter = b
        saveSettings()
        go startService()
    })
    dnsFilterCheck.SetChecked(currentCfg.DNSFilter)

    // Tarayıcıyı proxy ile başlat
    browserCheck := widget.NewCheck("Tarayıcıyı proxy ile başlat", func(b bool) {
        currentCfg.LaunchBrowser = b
        saveSettings()
        if b {
            go exec.Command("chrome.exe").Start()
        }
    })
    browserCheck.SetChecked(currentCfg.LaunchBrowser)

    resetBtn := widget.NewButton("Varsayılanlara Dön", func() {
        resetToDefaults()
        engineSelect.SetSelected(currentCfg.ActiveMode)
        dnsFilterCheck.SetChecked(currentCfg.DNSFilter)
        browserCheck.SetChecked(currentCfg.LaunchBrowser)
        go startService()
    })
    resetBtn.Importance = widget.LowImportance

    saveBtn := widget.NewButton("Uygula", func() {
        saveSettings()
        go startService()
        settingsWindow.Hide()
    })
    saveBtn.Importance = widget.HighImportance

    settingsWindow.SetContent(container.NewPadded(container.NewVScroll(container.NewVBox(
        widget.NewLabel("Proxy Motoru Seçimi:"), engineSelect,
        dnsFilterCheck,
        browserCheck,
        layout.NewSpacer(),
        container.NewVBox(resetBtn, saveBtn),
    ))))

    // Uygulama listesi (NDIS modunda tünellenecek uygulamalar)
    appListWindow = myApp.NewWindow("Uygulama Listesi")
    appListWindow.Resize(fyne.NewSize(600, 600))
    appListWindow.CenterOnScreen()
    appListWindow.SetCloseIntercept(func() { appListWindow.Hide() })

    input := widget.NewEntry(); input.SetPlaceHolder("Örn: discord.exe")
    appListContainer = container.NewVBox()

    refreshAppListFunc = func() {
        appListContainer.Objects = nil
        for i, app := range currentCfg.Apps {
            idx := i
            delBtn := widget.NewButton("Sil", func() {
                currentCfg.Apps = append(currentCfg.Apps[:idx], currentCfg.Apps[idx+1:]...)
                refreshAppListFunc()
            })
            delBtn.Importance = widget.DangerImportance

            protoSel := widget.NewSelect([]string{"TCP", "UDP", "TCP+UDP"}, func(s string) { currentCfg.Apps[idx].Protocol = s })
            protoSel.SetSelected(app.Protocol)

            modeSel := widget.NewSelect([]string{"Tunnel", "Exclude"}, func(s string) {
                currentCfg.Apps[idx].Mode = s
                refreshAppListFunc()
            })
            modeSel.SetSelected(app.Mode)

            nameText := canvas.NewText(app.Name, color.White)
            if app.Mode == "Exclude" {
                nameText.Color = color.NRGBA{R: 255, G: 0, B: 0, A: 255}
            } else {
                nameText.Color = color.NRGBA{R: 0, G: 150, B: 255, A: 255}
            }
            nameText.TextStyle = fyne.TextStyle{Bold: true}

            row := container.NewHBox(delBtn, protoSel, modeSel, nameText)
            appListContainer.Add(row)
        }
        appListContainer.Refresh()
    }

    addBtn := widget.NewButton("Ekle", func() {
        if input.Text != "" {
            currentCfg.Apps = append(currentCfg.Apps, AppConfig{Name: input.Text, Protocol: "TCP+UDP", Mode: "Tunnel"})
            input.SetText("")
            refreshAppListFunc()
        }
    })
    addBtn.Importance = widget.MediumImportance

    saveAppListBtn := widget.NewButton("Kaydet ve Yeniden Başlat", func() {
        saveSettings()
        go startService()
        appListWindow.Hide()
    })
    saveAppListBtn.Importance = widget.HighImportance

    appListWindow.SetContent(container.NewBorder(container.NewVBox(input, addBtn), saveAppListBtn, nil, nil, container.NewVScroll(appListContainer)))
}

func showSettingsGUI() { fyne.Do(func() { settingsWindow.Show() }) }
func showAppListGUI() { fyne.Do(func() { refreshAppListFunc(); appListWindow.Show() }) }
