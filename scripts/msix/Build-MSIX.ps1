# Build-MSIX.ps1 — 构建 unsigned MSIX(供 Sign-MSIX 签名)。
#
# 流程(Master O01):
#   clean staging
#   → Wails production build (Temper.exe)
#   → Defender raw scan (if available)
#   → copy exe/assets
#   → render AppxManifest version
#   → MakeAppx pack
#   → unsigned MSIX
#
# 用法:
#   powershell -ExecutionPolicy Bypass -File scripts/msix/Build-MSIX.ps1 [-Version 0.3.0.0]
# 输出:
#   dist/Temper-<ver>-windows-x64-unsigned.msix

param(
  [string]$Version = "0.3.0.0"
)

$ErrorActionPreference = "Stop"
$root = Resolve-Path (Join-Path $PSScriptRoot "..\..")
$desktopDir = Join-Path $root "desktop"
$distDir = Join-Path $root "dist"
$staging = Join-Path $distDir "msix-staging"
$template = Join-Path $root "packaging\msix\AppxManifest.xml.template"
$assetsDir = Join-Path $root "packaging\msix\assets"

# --- 1. clean staging ---
if (Test-Path $staging) { Remove-Item -Recurse -Force $staging }
New-Item -ItemType Directory -Path $staging -Force | Out-Null

# --- 2. Wails production build ---
Push-Location $desktopDir
try {
  Write-Host "==> wails build"
  # wails 是 .cmd 包装,PS 5 的 $LASTEXITCODE 不可靠;用 Start-Process 拿真实退出码。
  $wailsCmd = (Get-Command wails.cmd -ErrorAction SilentlyContinue).Source
  if (-not $wailsCmd) { $wailsCmd = (Get-Command wails -ErrorAction Stop).Source }
  $p = Start-Process -FilePath $wailsCmd -ArgumentList "build" -Wait -PassThru -NoNewWindow
  if ($p.ExitCode -ne 0) { throw "wails build failed (exit $($p.ExitCode))" }
} finally {
  Pop-Location
}
$exe = Join-Path $desktopDir "build\bin\Temper.exe"
if (-not (Test-Path $exe)) { throw "Temper.exe not found after wails build: $exe" }
Write-Host "Temper.exe: $((Get-Item $exe).Length / 1MB) MB"

# --- 3. Defender raw scan (best-effort; 不阻塞构建) ---
try {
  & (Join-Path $PSScriptRoot "..\defender-scan.ps1") -Path $exe
  if ($LASTEXITCODE -ne 0) { throw "Defender detected raw Temper.exe" }
} catch {
  Write-Warning "Defender scan skipped: $($_.Exception.Message)"
}

# --- 4. copy exe + assets ---
Copy-Item $exe (Join-Path $staging "Temper.exe")
if (Test-Path $assetsDir) {
  $stagingAssets = Join-Path $staging "assets"
  New-Item -ItemType Directory -Path $stagingAssets -Force | Out-Null
  Copy-Item (Join-Path $assetsDir "*") $stagingAssets -Recurse -Force
}

# --- 5. render AppxManifest version ---
$manifest = Join-Path $staging "AppxManifest.xml"
(Get-Content $template -Raw).Replace("__MSIX_VERSION__", $Version) | Set-Content $manifest -Encoding UTF8

# --- 6. locate MakeAppx (最高 SDK 版本 x64, 不硬编码) ---
$kitsRoot = Join-Path ${env:ProgramFiles(x86)} "Windows Kits\10\bin"
$sdk = Get-ChildItem $kitsRoot -Directory | Where-Object { $_.Name -match '^\d+\.\d+\.\d+\.\d+$' } |
  Sort-Object { [version]$_.Name } -Descending | Select-Object -First 1
if (-not $sdk) { throw "No Windows SDK version dir found under $kitsRoot" }
$makeAppx = Join-Path $sdk.FullName "x64\MakeAppx.exe"
if (-not (Test-Path $makeAppx)) { throw "MakeAppx.exe not found under $kitsRoot" }

# --- 7. MakeAppx pack ---
$unsigned = Join-Path $distDir "Temper-$Version-windows-x64-unsigned.msix"
if (Test-Path $unsigned) { Remove-Item -Force $unsigned }
$p = Start-Process -FilePath $makeAppx -ArgumentList @("pack", "/d", $staging, "/p", $unsigned) -Wait -PassThru -NoNewWindow
if ($p.ExitCode -ne 0) { throw "MakeAppx failed (exit $($p.ExitCode))" }
Write-Host "Unsigned MSIX: $unsigned ($([math]::Round((Get-Item $unsigned).Length / 1MB, 1)) MB)"
Write-Host "BUILD_MSIX_OK"
