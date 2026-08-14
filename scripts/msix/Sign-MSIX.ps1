# Sign-MSIX.ps1 — 对 MSIX 做 Authenticode 签名(SignTool sign /fd SHA256)。
#
# 本地优先证书存储 thumbprint(CN=Temper Development);
# CI 用临时 PFX(从 GitHub Secrets 解码,用完删除,绝不落盘/日志)。
#
# 用法:
#   powershell -ExecutionPolicy Bypass -File scripts/msix/Sign-MSIX.ps1 `
#     -Msix dist\Temper-0.3.0.0-windows-x64-unsigned.msix `
#     [-Thumbprint <cert-thumbprint>] [-PfxPath <temp.pfx>] [-PfxPassword <pw>]

param(
  [Parameter(Mandatory = $true)][string]$Msix,
  [string]$Thumbprint,
  [string]$PfxPath,
  [string]$PfxPassword,
  [string]$OutputName = ""
)

$ErrorActionPreference = "Stop"
$root = Resolve-Path (Join-Path $PSScriptRoot "..\..")

# --- 定位 SignTool(最高 SDK 版本 x64) ---
$kitsRoot = Join-Path ${env:ProgramFiles(x86)} "Windows Kits\10\bin"
$sdk = Get-ChildItem $kitsRoot -Directory | Where-Object { $_.Name -match '^\d+\.\d+\.\d+\.\d+$' } |
  Sort-Object { [version]$_.Name } -Descending | Select-Object -First 1
if (-not $sdk) { throw "No Windows SDK version dir found under $kitsRoot" }
$signTool = Join-Path $sdk.FullName "x64\SignTool.exe"
if (-not (Test-Path $signTool)) { throw "SignTool.exe not found under $kitsRoot" }

$output = [System.IO.Path]::Combine([System.IO.Path]::GetDirectoryName($Msix),
  ([System.IO.Path]::GetFileName($Msix) -replace "-unsigned", ""))
if ($OutputName -ne "") {
  $output = [System.IO.Path]::Combine([System.IO.Path]::GetDirectoryName($Msix), $OutputName)
}

$sigParams = @("sign", "/fd", "SHA256")
if ($PfxPath -and (Test-Path $PfxPath)) {
  if (-not $PfxPassword) { throw "PFX password required when signing from a PFX" }
  $sigParams = $sigParams + @("/f", $PfxPath, "/p", $PfxPassword)
} elseif ($Thumbprint) {
  $sigParams = $sigParams + @("/sha1", $Thumbprint)
} else {
  # 本地:自动找 CN=Temper Development 代码签名证书
  $cert = Get-ChildItem Cert:\CurrentUser\My -ErrorAction SilentlyContinue | Where-Object {
    $_.Subject -eq "CN=Temper Development" -and $_.HasPrivateKey
  } | Select-Object -First 1
  if (-not $cert) { throw "No Temper Development cert found; run New-DevCertificate.ps1" }
  $sigParams = $sigParams + @("/sha1", $cert.Thumbprint)
}

# SignTool 就地签名(此版本不支持 /out);签名后重命名为正式名。
$sigParams = $sigParams + @($Msix)
Write-Host "==> SignTool sign /fd SHA256 (in place) -> $Msix"
$oldPref = $ErrorActionPreference
$ErrorActionPreference = "Continue"

& $signTool @sigParams
$code = $LASTEXITCODE
$ErrorActionPreference = $oldPref
if ($code -ne 0) { throw "SignTool sign failed (exit $code)" }

# 就地签名完成 → 重命名为正式名(移除 -unsigned)。
if (Test-Path $output) { Remove-Item -Force $output }
Move-Item -Force $Msix $output
Write-Host "SIGN_OK"
Write-Host "SIGNED_MSIX=$output"

# 清理:绝不保留 PFX 引用
if ($PfxPath -and (Test-Path $PfxPath)) { Remove-Item -Force $PfxPath }
Write-Host "SIGN_OK"
Write-Host "SIGNED_MSIX=$output"
