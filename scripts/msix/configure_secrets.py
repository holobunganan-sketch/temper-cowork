# 配置 TEMPER_MSIX_PFX_B64 / TEMPER_MSIX_PFX_PASSWORD 到 GitHub Secrets。
# 用法: python scripts/msix/configure_secrets.py
# 依赖: pip install pynacl;环境变量 GITHUB_PERSONAL_ACCESS_TOKEN 已认证。
import base64
import os
import secrets as pysecrets
import subprocess
import sys
import urllib.request

from nacl.public import PublicKey, SealedBox
import nacl.utils

REPO = "holobunganan-sketch/temper-cowork"
TOKEN = os.environ.get("GITHUB_PERSONAL_ACCESS_TOKEN", "")
if not TOKEN:
    sys.exit("GITHUB_PERSONAL_ACCESS_TOKEN not set")

API = f"https://api.github.com/repos/{REPO}"


def api(path, method="GET", data=None):
    req = urllib.request.Request(
        API + path,
        data=json_dumps(data).encode() if data is not None else None,
        method=method,
        headers={
            "Authorization": f"Bearer {TOKEN}",
            "Accept": "application/vnd.github+json",
            "Content-Type": "application/json",
        },
    )
    with urllib.request.urlopen(req) as resp:
        return json_loads(resp.read())


import json


def json_dumps(o):
    return json.dumps(o)


def json_loads(b):
    return json.loads(b.decode())


def main():
    # 1. 导出本地证书为临时 PFX(随机强密码)
    password = pysecrets.token_urlsafe(32)
    pfx_path = os.path.join(os.environ.get("TEMP", "."), "temper-dev-signing.pfx")
    ps = (
        "$c = Get-ChildItem Cert:\\CurrentUser\\My | Where-Object { $_.Subject -eq "
        "'CN=Temper Development' -and $_.HasPrivateKey } | Select-Object -First 1; "
        f"$c | Export-PfxCertificate -FilePath '{pfx_path}' -Password "
        f"(ConvertTo-SecureString '{password}' -AsPlainText -Force) | Out-Null"
    )
    subprocess.run(["powershell", "-NoProfile", "-Command", ps], check=True)

    with open(pfx_path, "rb") as f:
        pfx_bytes = f.read()
    os.remove(pfx_path)  # 立即删除临时 PFX

    pfx_b64 = base64.b64encode(pfx_bytes).decode()

    # 2. 获取 public key,libsodium 加密
    pub = api("/actions/secrets/public-key")
    key = base64.b64decode(pub["key"])
    key_id = pub["key_id"]

    def encrypt(value):
        sealed = SealedBox(PublicKey(key))
        return base64.b64encode(sealed.encrypt(value.encode())).decode()

    # 3. PUT secrets
    api("/actions/secrets/TEMPER_MSIX_PFX_B64", method="PUT",
        data={"encrypted_value": encrypt(pfx_b64), "key_id": key_id})
    print("TEMPER_MSIX_PFX_B64 set")
    api("/actions/secrets/TEMPER_MSIX_PFX_PASSWORD", method="PUT",
        data={"encrypted_value": encrypt(password), "key_id": key_id})
    print("TEMPER_MSIX_PFX_PASSWORD set")

    # 清内存中的敏感值
    del pfx_b64, password
    print("SECRETS_OK")


if __name__ == "__main__":
    main()
