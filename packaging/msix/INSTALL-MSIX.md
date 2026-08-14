# 安装 Temper(MSIX)

Temper v0.3.0 以 MSIX 形式发行。MSIX 是 Windows 10/11 的现代打包格式,
支持干净安装/升级/卸载,数据保存在用户目录(卸载不删除你的 workspace)。

## 前置:信任开发证书

v0.3.0 使用开发自签名证书 `CN=Temper Development`(公钥随发行物提供)。
首次安装前,必须以**管理员身份**信任该证书:

```powershell
# 以管理员身份打开 PowerShell
Import-Certificate -FilePath .\Temper-Development.cer `
  -CertStoreLocation Cert:\LocalMachine\TrustedPeople
```

> 信任仅影响本机;该证书只用于签名 Temper 开发构建,不用于其他用途。

## 安装

```powershell
Add-AppxPackage -Path .\Temper-0.3.0-windows-x64.msix
```

安装完成后,从**开始菜单**搜索 `Temper` 启动。

## 验证

- 签名校验:`signtool verify /pa /v Temper-0.3.0-windows-x64.msix`
- 完整性校验:见 `Temper-0.3.0-SHA256SUMS.txt`

## 数据位置

| 内容 | 路径 |
|------|------|
| 运行时状态 | `%APPDATA%\Temper\runtime` |
| CoWork 业务库 | `%APPDATA%\Temper\cowork`(temper.db) |
| 缓存 | `%LOCALAPPDATA%\Temper\cache` |
| 你的 workspace | 你选择的位置(移除项目不会删除) |

## 卸载

设置 → 应用 → Temper → 卸载。卸载后 workspace 文件与
`%APPDATA%\Temper\cowork` 中的数据保留(如需彻底清除请手动删除)。

## 常见问题

- **安装报错 "证书不受信任"**:按上述步骤导入 `Temper-Development.cer`。
- **安装报错 "已安装更高版本"**:先卸载旧版本再装。
- **SmartScreen 提示**:自签名开发证书会触发提示;确认 SHA256 校验一致后
  选择"仍要运行"。正式发布版建议使用商业代码签名证书以消除该提示。
