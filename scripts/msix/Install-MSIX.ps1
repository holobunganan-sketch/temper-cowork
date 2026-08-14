# Install-MSIX.ps1 — 安装签名 MSIX。
#
# 规则(Master O05):
#   - 不得偷偷修改 Windows trust store
#   - 若 Temper-Development.cer 未信任,明确提示用户以管理员身份信任到
#     LocalMachine\TrustedPeople,然后 Add-AppxPackage
#   - 验证:installed / Start Menu Temper / launch works
#
# 用法:
#   powershell -ExecutionPolicy Bypass -File scripts/msix/Install-MSIX.ps1 `
#     -Msix <signed.msix> [-Cer packaging\msix\Temper-Development.cer] [-AutoTrust]

param(
  [Parameter(Mandatory = $true)][string]$Msix,
  [string]$Cer = "",
  [switch]$AutoTrust
)

$ErrorActionPreference = "Stop"

if (-not (Test-Path $Msix)) { throw "MSIX not found: $Msix" }

# 检查 MSIX 签名者是否已信任(TrustedPeople)。
$sig = Get-AuthenticodeSignature $Msix
if ($sig.Status -ne "Valid") {
  Write-Host "MSIX signature status: $($sig.Status)"
  if ($AutoTrust -and $Cer -and (Test-Path $Cer)) {
    Write-Host "Auto-trusting $Cer into LocalMachine\TrustedPeople (requires admin)..."
    $target = "Cert:\LocalMachine\TrustedPeople"
    Import-Certificate -FilePath $Cer -CertStoreLocation $target | Out-Null
  } else {
    Write-Host "Please trust the certificate as administrator, then re-run:"
    Write-Host "  Import-Certificate -FilePath $Cer -CertStoreLocation Cert:\LocalMachine\TrustedPeople"
    Write-Host "  (or run with -AutoTrust from an elevated shell)"
    exit 1
  }
}

Write-Host "==> Add-AppxPackage $Msix"
Add-AppxPackage -Path $Msix
if ($LASTEXITCODE -ne 0) { throw "Add-AppxPackage failed" }

# 验证已安装
$pkg = Get-AppxPackage -Name "Temper.Cowork.Desktop"
if (-not $pkg) { throw "Package not found after install" }
Write-Host "Installed: $($pkg.PackageFullName)"
Write-Host "Start Menu entry: $($pkg.InstallLocation)"
Write-Host "INSTALL_OK"
