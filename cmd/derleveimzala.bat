@echo off
chcp 65001 > nul
title ProksiTool Detaylı Derleyici ve İmzalayıcı

:: OTO YÖNETİCİ KONTROLÜ
net session >nul 2>&1
if %errorLevel% neq 0 (
    echo [BİLGİ] Yönetici yetkisi alınıyor...
    powershell -Command "Start-Process '%~f0' -Verb RunAs"
    exit /b
)

:: BAŞKAN: Klasör şaşırma bugını çözen kritik satır burası!
:: Script yönetici olarak açılsa bile projenin olduğu gerçek klasöre geri döner.
cd /d "%~dp0"

:: =====================================================================================
:: BAŞKAN: SADECE AŞAĞIDAKİ ÇİFT TIRNAK İÇİNDEKİ ALANLARI DEĞİŞTİR!
:: =====================================================================================
set "S_PROGRAM_ADI=ProksiTool"
set "S_SIRKET_ADI=ZuLuSoft Development Inc"
set "S_DEPARTMAN=Yazilim Departmani"
set "S_SEHIR=Rize"
set "S_BOLGE=Karadeniz"
set "S_ULKE=TR"
set "S_YIL=5"
:: =====================================================================================

echo ============================================================
echo [1/3] Go Projesi Derleniyor (Semboller Kazınıyor)...
echo ============================================================
go build -ldflags="-s -w -H windowsgui" -buildvcs=false -trimpath -o ProksiTool_1.06.exe

if %errorlevel% neq 0 (
    echo [HATA] Go derleme hatası verdi! Kodları kontrol et başkan.
    pause
    exit /b
)
timeout 1
echo.
echo ============================================================
echo [2/3] PowerShell ile Özel Bilgili Sertifika Basılıyor...
echo ============================================================

set "SUBJECT=CN=%S_PROGRAM_ADI%, O=%S_SIRKET_ADI%, OU=%S_DEPARTMAN%, L=%S_SEHIR%, S=%S_BOLGE%, C=%S_ULKE%"

powershell -Command "$cert = New-SelfSignedCertificate -Type CodeSigningCert -Subject '%SUBJECT%' -FriendlyName '%S_PROGRAM_ADI% Resmi Sertifikasi' -NotAfter (Get-Date).AddYears(%S_YIL%); Move-Item -Path \"Cert:\CurrentUser\My\$($cert.Thumbprint)\" -Destination 'Cert:\CurrentUser\Root' -Force -ErrorAction SilentlyContinue; Move-Item -Path \"Cert:\LocalMachine\My\$($cert.Thumbprint)\" -Destination 'Cert:\LocalMachine\Root' -Force -ErrorAction SilentlyContinue; $signResult = Set-AuthenticodeSignature -FilePath 'ProksiTool_1.06.exe' -Certificate $cert; Write-Output $signResult"

echo.
echo ============================================================
echo [3/3] İşlem Tamamlandı!
echo ============================================================
echo exe dosyana sağ tıklayıp Dijital İmzalar sekmesini kontrol edebilirsin.
echo.
pause