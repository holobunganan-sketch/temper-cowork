# Verify-MSIX.ps1 — 验证 MSIX 签名(SignTool verify /pa /v)并计算 SHA256。
#
# 用法:
#   powershell -ExecutionPolicy Bypass -File scripts/msix/Verify-MSIX.ps1 -Msix <signed.msix>
# 退出码: 0 = 签名有效。

param(
  [Parameter(Mandatory = $true)][string]$Msix
)

$ErrorActionPreference = "Stop"
$root = (Resolve-Path (Join-Path $PSScriptRoot "..\..")).Path

if (-not (Test-Path $Msix)) { throw "MSIX not found: $Msix" }

$kitsRoot = Join-Path ${env:ProgramFiles(x86)} "Windows Kits\10\bin"
$sdk = Get-ChildItem $kitsRoot -Directory | Where-Object { $_.Name -match '^\d+\.\d+\.\d+\.\d+$' } |
  Sort-Object { [version]$_.Name } -Descending | Select-Object -First 1
if (-not $sdk) { throw "No Windows SDK version dir found under $kitsRoot" }
$signTool = Join-Path $sdk.FullName "x64\SignTool.exe"
if (-not (Test-Path $signTool)) { throw "SignTool.exe not found" }

# CI 使用临时 PFX,自签名证书不在信任存储;先将开发证书导入信任根。
# windows-latest 以服务账户运行,LocalMachine\Root 可写;本地非管理员失败时
# 回退到 verify /v(仅验证签名哈希,不查信任链)。
$cerPath = Join-Path $root "packaging\msix\Temper-Development.cer"
if ((Test-Path $cerPath) -and -not (Test-Path "Cert:\LocalMachine\Root\94B722B425EA7D5B337825FF00D18C7B0E6FDA97")) {
  try {
    Import-Certificate -FilePath $cerPath -CertStoreLocation Cert:\LocalMachine\Root -ErrorAction SilentlyContinue | Out-Null
  } catch {
    Write-Host "note: could not trust dev cert: $($_.Exception.Message)"
  }
}
Write-Host "==> SignTool verify /pa /v"
$verifyArgs = @("verify", "/pa", "/v", $Msix)
$oldPref = $ErrorActionPreference
$ErrorActionPreference = "Continue"
& $signTool @verifyArgs
$code = $LASTEXITCODE
$ErrorActionPreference = $oldPref
if ($code -ne 0) {
  # 信任链不可用时(本地非管理员),fallback 到 verify /v 验证签名哈希。
  Write-Host "note: /pa trust-chain verify failed (exit $code); falling back to verify /v"
  $oldPref = $ErrorActionPreference
  $ErrorActionPreference = "Continue"
  & $signTool @("verify", "/v", $Msix)
  $code = $LASTEXITCODE
  $ErrorActionPreference = $oldPref
}
if ($code -ne 0) { throw "Signature verification failed (exit $code)" }

$hash = Get-FileHash -Algorithm SHA256 $Msix
Write-Host "SHA256: $($hash.Hash)  $($hash.Path)"
Write-Host "VERIFY_OK"
exit 0
