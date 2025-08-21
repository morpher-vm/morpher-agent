#!/usr/bin/env bash
set -euo pipefail

: "${MORPHER_CONTROLLER_IP:?set MORPHER_CONTROLLER_IP (e.g. 192.168.54.3)}"

REPO="${REPO:-morpher-vm/morpher-agent}"
VERSION="${VERSION:-v0.0.1}"               # e.g. v0.1.0 or latest
SERVICE_NAME="${SERVICE_NAME:-morpher-agent}"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"
CONFIG_DIR="/etc/${SERVICE_NAME}"

ASSET_FILE="${ASSET_FILE:-${SERVICE_NAME}_Linux_arm64.tar.gz}"
BASE_URL="https://github.com/${REPO}/releases"
DL_URL="${BASE_URL}/download/${VERSION}/${ASSET_FILE}"
[[ "${VERSION}" == "latest" ]] && DL_URL="https://github.com/${REPO}/releases/latest/download/${ASSET_FILE}"

tmpdir="$(mktemp -d)"; trap 'rm -rf "$tmpdir"' EXIT
echo "[*] Download ${DL_URL}"
curl -fsSL -o "${tmpdir}/asset.tgz" "${DL_URL}"

# Download and verify checksum
cs_url="${BASE_URL}/download/${VERSION}/checksums.txt"
[[ "${VERSION}" == "latest" ]] && cs_url="https://github.com/${REPO}/releases/latest/download/checksums.txt"

curl -fsSL -o "${tmpdir}/checksums.txt" "${cs_url}"

line="$(grep -E "  ${ASSET_FILE}\$" "${tmpdir}/checksums.txt" || true)"
if [[ -z "${line}" ]]; then
  echo "[X] ${ASSET_FILE} entry not found in checksums.txt (${cs_url})"
  exit 1
fi

hash="$(echo "${line}" | awk '{print $1}')"
echo "[*] Verifying checksum for ${ASSET_FILE}"
echo "${hash}  asset.tgz" | (cd "${tmpdir}" && sha256sum -c -)

echo "[*] Extract"
tar -xzf "${tmpdir}/asset.tgz" -C "${tmpdir}"

bin_path="$(find "${tmpdir}" -type f -executable -name "${SERVICE_NAME}" | head -n1)"
[[ -z "${bin_path}" ]] && { echo "binary not found: ${SERVICE_NAME}"; exit 1; }

echo "[*] Install -> ${INSTALL_DIR}"
install -d "${INSTALL_DIR}"
install -m 0755 "${bin_path}" "${INSTALL_DIR}/${SERVICE_NAME}"

echo "[*] Config -> ${CONFIG_DIR}"
install -d "${CONFIG_DIR}"
cat > "${CONFIG_DIR}/${SERVICE_NAME}.env" <<EOF
MORPHER_AGENT_BASE_URL=http://${MORPHER_CONTROLLER_IP}:9000
EOF
chmod 600 "${CONFIG_DIR}/${SERVICE_NAME}.env"

unit="/etc/systemd/system/${SERVICE_NAME}.service"
echo "[*] systemd unit -> ${unit}"
cat > "${unit}" <<EOF
[Unit]
Description=${SERVICE_NAME}
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
EnvironmentFile=-${CONFIG_DIR}/${SERVICE_NAME}.env
ExecStart=${INSTALL_DIR}/${SERVICE_NAME} start
Restart=on-failure

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable --now "${SERVICE_NAME}.service"
systemctl --no-pager --full status "${SERVICE_NAME}.service" || true

echo "✅ Done"