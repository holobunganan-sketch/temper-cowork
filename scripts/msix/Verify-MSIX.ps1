# Verify-MSIX.ps1 — 验证 MSIX 签名(SignTool verify /pa /v)并计算 SHA256。
#
# 用法:
#   powershell -ExecutionPolicy Bypass -File scripts/msix/Verify-MSIX.ps1 -Msix <signed.msix>
# 退出码: 0 = 签名有效。

param(
  [Parameter(Mandatory = $true)][string]$Msix
)

$ErrorActionPreference = "Stop"
$root = Resolve-Path (Join-Path $PSScriptRoot "..\..")

if (-not (Test-Path $Msix)) { throw "MSIX not found: $Msix" }

$kitsRoot = Join-Path ${env:ProgramFiles(x86)} "Windows Kits\10\bin"
$sdk = Get-ChildItem $kitsRoot -Directory | Where-Object { $_.Name -match '^\d+\.\d+\.\d+\.\d+$' } |
  Sort-Object { [version]$_.Name } -Descending | Select-Object -First 1
if (-not $sdk) { throw "No Windows SDK version dir found under $kitsRoot" }
$signTool = Join-Path $sdk.FullName "x64\SignTool.exe"
if (-not (Test-Path $signTool)) { throw "SignTool.exe not found" }

Write-Host "==> SignTool verify /pa /v"
$verifyArgs = @("verify", "/pa", "/v", $Msix)
$oldPref = $ErrorActionPreference
$ErrorActionPreference = "Continue"
& $signTool @verifyArgs
$code = $LASTEXITCODE
$ErrorActionPreference = $oldPref
if ($code -ne 0) { throw "Signature verification failed (exit $code)" }

$hash = Get-FileHash -Algorithm SHA256 $Msix
Write-Host "SHA256: $($hash.Hash)  $($hash.Path)"
Write-Host "VERIFY_OK"
exit 0
