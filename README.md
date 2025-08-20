# morpher-agent

<p align="center">
  <img src="https://github.com/morpher-vm.png" alt="morpher-vm logo" width="200"/>
</p>

**morpher-agent** is a lightweight agent that collects VM information for migration.

# Installation

You can install the agent using the provided arch-specific install scripts.  
These scripts will download the release tarball, extract the binary, create a
systemd service, and start it automatically.

> **Note:** Requires `sudo` privileges.

---

## Quick Start

### x86_64
```bash
wget -O install_amd64.sh https://raw.githubusercontent.com/morpher-vm/morpher-agent/main/scripts/install_amd64.sh
chmod +x install_amd64.sh
sudo MORPHER_CONTROLLER_IP=<controller-ip> ./install_amd64.sh
```

### ARM64
```bash
wget -O install_arm64.sh https://raw.githubusercontent.com/morpher-vm/morpher-agent/main/scripts/install_arm64.sh
chmod +x install_arm64.sh
sudo MORPHER_CONTROLLER_IP=<controller-ip> ./install_arm64.sh
```

## Environment Variables

- `MORPHER_CONTROLLER_IP`: The IP address of the Morpher controller to connect to.

## Service Management

After installation, the morpher-agent service is managed via systemd:

```bash
sudo systemctl start morpher-agent
sudo systemctl enable morpher-agent
sudo systemctl status morpher-agent
sudo journalctl -u morpher-agent -f
```

## License

Licensed under the [MIT License](LICENSE).
