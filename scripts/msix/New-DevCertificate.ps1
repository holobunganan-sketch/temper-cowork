# New-DevCertificate.ps1 — 查找或创建 CN=Temper Development 开发代码签名证书。
#
# 规则(Master O02):
#   - 已存在 CN=Temper Development + Code Signing EKU + HasPrivateKey → 复用
#   - 否则创建:CurrentUser\My, DigitalSignature, Code Signing EKU,
#     5 年有效期, exportable private key
#   - 导出公钥到 packaging/msix/Temper-Development.cer(允许提交)
#   - 私钥永不导出到仓库
#
# 用法: powershell -ExecutionPolicy Bypass -File scripts/msix/New-DevCertificate.ps1

$ErrorActionPreference = "Stop"

$subject = "CN=Temper Development"
$cerPath = Join-Path $PSScriptRoot "..\..\packaging\msix\Temper-Development.cer"

function Find-TemperCert {
  Get-ChildItem Cert:\CurrentUser\My -ErrorAction SilentlyContinue | Where-Object {
    $_.Subject -eq $subject -and
    ($_.EnhancedKeyUsageList | Where-Object { $_.FriendlyName -eq "Code Signing" }) -and
    $_.HasPrivateKey
  } | Select-Object -First 1
}

$cert = Find-TemperCert
if ($cert) {
  Write-Host "Reusing existing cert: $($cert.Subject) thumbprint=$($cert.Thumbprint) expires=$($cert.NotAfter)"
} else {
  Write-Host "Creating new code-signing cert: $subject"
  $cert = New-SelfSignedCertificate `
    -Subject $subject `
    -CertStoreLocation Cert:\CurrentUser\My `
    -Type CodeSigningCert `
    -KeyUsage DigitalSignature `
    -KeyExportPolicy Exportable `
    -NotAfter (Get-Date).AddYears(5)
  Write-Host "Created thumbprint=$($cert.Thumbprint)"
}

# 导出公钥 CER(可提交;不含私钥)
$certDir = Split-Path $cerPath -Parent
if (-not (Test-Path $certDir)) { New-Item -ItemType Directory -Path $certDir -Force | Out-Null }
Export-Certificate -Cert $cert -FilePath $cerPath -Type CERT | Out-Null
Write-Host "Exported public cert: $cerPath"

Write-Host ("CERT_THUMBPRINT=" + $cert.Thumbprint)
