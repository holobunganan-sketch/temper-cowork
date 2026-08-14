# Configure-GitHubSigning.ps1 — 把本地开发证书导出为临时 PFX 并写入
# GitHub Secrets(TEMPER_MSIX_PFX_B64 / TEMPER_MSIX_PFX_PASSWORD)。
#
# 规则(Master O03):
#   - 读取本地 CN=Temper Development 证书
#   - 生成随机强 PFX 密码
#   - 临时 PFX → Base64 → gh secret set → 删除临时 PFX → 清空密码变量
#   - 禁止日志输出 secret
#
# 用法(需 gh 已认证且对该仓库有 secrets 权限):
#   powershell -ExecutionPolicy Bypass -File scripts/msix/Configure-GitHubSigning.ps1 `
#     -Repo holobunganan-sketch/temper-cowork

param(
  [Parameter(Mandatory = $true)][string]$Repo
)

$ErrorActionPreference = "Stop"

$cert = Get-ChildItem Cert:\CurrentUser\My -ErrorAction SilentlyContinue | Where-Object {
  $_.Subject -eq "CN=Temper Development" -and $_.HasPrivateKey
} | Select-Object -First 1
if (-not $cert) { throw "No Temper Development cert found; run New-DevCertificate.ps1 first" }

# 生成随机强密码(32 字节 base64,URL-safe)。
$bytes = New-Object byte[] 32
[System.Security.Cryptography.RandomNumberGenerator]::Fill($bytes)
$password = [Convert]::ToBase64String($bytes)

$tempPfx = Join-Path $env:TEMP "temper-dev-signing.pfx"
try {
  # 导出含私钥的 PFX(需密码保护)。
  $cert | Export-PfxCertificate -FilePath $tempPfx -Password (ConvertTo-SecureString $password -AsPlainText -Force) | Out-Null

  $pfxB64 = [Convert]::ToBase64String([System.IO.File]::ReadAllBytes($tempPfx))

  Write-Host "Setting TEMPER_MSIX_PFX_B64 ($($pfxB64.Length) chars)..."
  gh secret set TEMPER_MSIX_PFX_B64 --repo $Repo --body $pfxB64
  if ($LASTEXITCODE -ne 0) { throw "gh secret set TEMPER_MSIX_PFX_B64 failed" }

  Write-Host "Setting TEMPER_MSIX_PFX_PASSWORD..."
  gh secret set TEMPER_MSIX_PFX_PASSWORD --repo $Repo --body $password
  if ($LASTEXITCODE -ne 0) { throw "gh secret set TEMPER_MSIX_PFX_PASSWORD failed" }

  Write-Host "SECRETS_OK"
} finally {
  # 删除临时 PFX + 清空密码(绝不落日志)。
  if (Test-Path $tempPfx) { Remove-Item -Force $tempPfx }
  $password = $null
  $pfxB64 = $null
  [System.GC]::Collect()
}
