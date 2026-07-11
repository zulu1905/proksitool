package main

import (
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "syscall"
    "unsafe"
    "golang.org/x/sys/windows/registry"
)

const CREATE_NO_WINDOW = 0x08000000
const startupName = "ProksiTool"

var (
    modkernel32           = syscall.NewLazyDLL("kernel32.dll")
    procGetShortPathNameW = modkernel32.NewProc("GetShortPathNameW")
    workDir     string
    desktopDir  string
)

func getShortPath(path string) string {
    utf16Path, _ := syscall.UTF16PtrFromString(path)
    buf := make([]uint16, 260)
    ret, _, _ := procGetShortPathNameW.Call(
        uintptr(unsafe.Pointer(utf16Path)),
        uintptr(unsafe.Pointer(&buf[0])),
        uintptr(len(buf)),
    )
    if ret == 0 { return path }
    return syscall.UTF16ToString(buf)
}

func init() {
    execPath, err := os.Executable()
    if err != nil {
        execPath, _ = os.Getwd()
    }
    rawWork := filepath.Dir(execPath)
    workDir = getShortPath(rawWork)

    rawDesktop := filepath.Join(os.Getenv("USERPROFILE"), "Desktop")
    desktopDir = getShortPath(rawDesktop)
}

func isAutoStartEnabled() bool {
    k, err := registry.OpenKey(registry.CURRENT_USER,
        `Software\Microsoft\Windows\CurrentVersion\Run`, registry.QUERY_VALUE)
    if err != nil { return false }
    defer k.Close()
    _, _, err = k.GetStringValue(startupName)
    return err == nil
}

func toggleAutoStart() bool {
    selfPath, err := os.Executable()
    if err != nil { return false }
    if isAutoStartEnabled() {
        k, err := registry.OpenKey(registry.CURRENT_USER,
            `Software\Microsoft\Windows\CurrentVersion\Run`, registry.SET_VALUE)
        if err == nil {
            _ = k.DeleteValue(startupName)
            k.Close()
        }
        return false
    } else {
        k, err := registry.OpenKey(registry.CURRENT_USER,
            `Software\Microsoft\Windows\CurrentVersion\Run`, registry.SET_VALUE)
        if err == nil {
            _ = k.SetStringValue(startupName, `"`+selfPath+`"`)
            k.Close()
            return true
        }
        return false
    }
}

func launchBrowserWithProxy() {
    browsers := []string{"chrome", "opera", "firefox", "brave", "msedge"}
    selectedBrowser := "msedge"
    for _, b := range browsers {
        checkCmd := exec.Command("cmd", "/c", "where", b)
        checkCmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: CREATE_NO_WINDOW}
        if err := checkCmd.Run(); err == nil {
            selectedBrowser = b
            break
        }
    }
    proxyArg := fmt.Sprintf("--proxy-server=socks5://%s:%s", currentCfg.Addr, currentCfg.Port)
    cmd := exec.Command("cmd", "/c", "start", selectedBrowser,
        proxyArg,
        "--no-first-run",
        "https://www.google.com")
    cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: CREATE_NO_WINDOW}
    _ = cmd.Run()
}
