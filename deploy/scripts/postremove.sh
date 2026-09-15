#!/bin/sh
# RPM/DEB 卸载后：停止并禁用服务，重载 systemd
if command -v systemctl >/dev/null 2>&1; then
  systemctl stop linmo || true
  systemctl disable linmo || true
  systemctl daemon-reload || true
fi
