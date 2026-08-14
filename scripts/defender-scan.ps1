# scripts/defender-scan.ps1 — 用 Windows Defender 扫描给定路径(EXE/MSIX)。
#
# 用法: powershell -ExecutionPolicy Bypass -File scripts/defender-scan.ps1 -Path <file>
# 退出码: 0 = 干净, 1 = 检测到威胁或扫描失败。
#
# 设计:不关闭/不禁用 Defender;扫描后检查 Get-MpThreatDetection 是否出现
# 指向该路径的检测记录。

param(
  [Parameter(Mandatory = $true)]
  [string]$Path
)

$ErrorActionPreference = "Stop"

if (-not (Test-Path $Path)) {
  Write-Error "path not found: $Path"
  exit 2
}

$resolved = (Resolve-Path $Path).Path

# 记录扫描前的检测数量,用于识别本次新增的检测。
$allBefore = @(Get-MpThreatDetection -ErrorAction SilentlyContinue)
$before = 0
foreach ($d in $allBefore) {
  if ($d.Resources -like "*$resolved*") { $before = $before + 1 }
}

Write-Host "Defender: scanning $resolved"
Start-MpScan -ScanType CustomScan -ScanPath $resolved | Out-Null
# 给实时保护/引擎一点时间刷新检测。
Start-Sleep -Seconds 5

$allAfter = @(Get-MpThreatDetection -ErrorAction SilentlyContinue)
$detections = @()
foreach ($d in $allAfter) {
  if ($d.Resources -like "*$resolved*") { $detections += $d }
}

$new = $detections.Count - $before
if ($new -gt 0) {
  Write-Host "DEFENDER_DETECTED ($new new detection(s))"
  foreach ($d in $detections) {
    Write-Host ("ThreatID: {0}  Resources: {1}" -f $d.ThreatID, ($d.Resources -join "; "))
  }
  exit 1
}

Write-Host "DEFENDER_CLEAN"
exit 0
