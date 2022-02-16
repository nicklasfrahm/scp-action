# 🚀 SCP File Transfer for GitHub Actions

[![Go Report Card](https://goreportcard.com/badge/github.com/nicklasfrahm/scp-action)](https://goreportcard.com/report/github.com/nicklasfrahm/scp-action)
[![container](https://github.com/nicklasfrahm/scp-action/actions/workflows/container.yml/badge.svg?branch=main)](https://github.com/nicklasfrahm/scp-action/actions/workflows/container.yml)

A [GitHub Action](https://github.com/features/actions) to upload and download files via SCP.

## Usage

Please note that if you only specify a single file as source, the target must be a file name and not a folder.

### 🔼 Uploading local files to remote target

```yaml
name: upload

on:
  - push

jobs:
  upload:
    name: Upload
    runs-on: ubuntu-latest
    steps:
      - name: Check out repository
        uses: actions/checkout@master

      - name: Upload file via SSH
        uses: nicklasfrahm/scp-action@main
        with:
          direction: upload
          host: ${{ secrets.SSH_TARGET_HOST }}
          fingerprint: ${{ secrets.SSH_TARGET_FINGERPRINT }}
          username: ${{ secrets.SSH_USER }}
          key: ${{ secrets.SSH_PRIVATE_KEY }}
          proxy_host: ${{ secrets.SSH_PROXY_HOST }}
          proxy_fingerprint: ${{ secrets.SSH_PROXY_FINGERPRINT }}
          proxy_username: ${{ secrets.SSH_USER }}
          proxy_key: ${{ secrets.SSH_PRIVATE_KEY }}
          source: |
            path/to/source/a.txt
            path/to/source/b.txt
          target: path/to/target
```

### 🔽 Downloading remote files to local target

```yaml
name: download

on:
  - push

jobs:
  download:
    name: Download
    runs-on: ubuntu-latest
    steps:
      - name: Check out repository
        uses: actions/checkout@master

      - name: Download file via SSH
        uses: nicklasfrahm/scp-action@main
        with:
          direction: download
          host: ${{ secrets.SSH_TARGET_HOST }}
          fingerprint: ${{ secrets.SSH_TARGET_FINGERPRINT }}
          username: ${{ secrets.SSH_USER }}
          key: ${{ secrets.SSH_PRIVATE_KEY }}
          proxy_host: ${{ secrets.SSH_PROXY_HOST }}
          proxy_fingerprint: ${{ secrets.SSH_PROXY_FINGERPRINT }}
          proxy_username: ${{ secrets.SSH_USER }}
          proxy_key: ${{ secrets.SSH_PRIVATE_KEY }}
          source: path/to/source/a.txt
          target: path/to/target/b.txt
```

## Input variables

See [action.yml](./action.yml) for more detailed information. **Please note that all input variables must have string values. It is thus recommend to always use quotes.**

| Input variable                      | Default value | Description                                                                                                                      |
| ----------------------------------- | ------------- | -------------------------------------------------------------------------------------------------------------------------------- |
| `host`                              | _required_    | SSH host                                                                                                                         |
| `port`                              | `22`          | SSH port                                                                                                                         |
| `username`                          | `root`        | SSH username                                                                                                                     |
| `passphrase`                        | _none_        | SSH passphrase                                                                                                                   |
| `insecure_password`                 | _none_        | SSH password, not recommended for security reasons                                                                               |
| `timeout`                           | `30s`         | Timeout for SSH connection to remote host                                                                                        |
| `action_timeout`                    | `10m`         | Timeout for action                                                                                                               |
| `key`                               | _none_        | Content of ssh private key, raw content of `~/.ssh/id_rsa`                                                                       |
| `fingerprint`                       | _none_        | Fingerprint SHA256 of the host public key, see [Using host fingerprint verification](#using-host-fingerprint-verification)       |
| `insecure_ignore_fingerprint`       | `false`       | Skip fingerprint verification of the host public key, not recommended for security reasons                                       |
| `source`                            | _required_    | A list of files to copy                                                                                                          |
| `target`                            | `.`           | A folder to copy to                                                                                                              |
| `direction`                         | _none_        | Transfer direction, must be either `upload` or `download`                                                                        |
| `proxy_host`                        | _none_        | SSH proxy host                                                                                                                   |
| `proxy_port`                        | `22`          | SSH proxy port                                                                                                                   |
| `proxy_username`                    | `root`        | SSH proxy username                                                                                                               |
| `proxy_passphrase`                  | _none_        | SSH proxy passphrase                                                                                                             |
| `insecure_proxy_password`           | _none_        | SSH proxy password                                                                                                               |
| `proxy_key`                         | _none_        | Content of SSH proxy private key                                                                                                 |
| `proxy_fingerprint`                 | _none_        | Fingerprint SHA256 of the proxy host public key, see [Using host fingerprint verification](#using-host-fingerprint-verification) |
| `insecure_proxy_ignore_fingerprint` | _none_        | Skip fingerprint verification of the proxy host public key, not recommended for security reasons                                 |

## Using host fingerprint verification

Setting up SSH host fingerprint verification can help to prevent Person-in-the-Middle attacks. Before setting this up, run the command below to get your SSH host fingerprint. Remember to replace `ed25519` with your appropriate key type (`rsa`, `dsa`, etc.) that your server is using and `example.com` with your host. In modern OpenSSH releases, the _default_ key types to be fetched are `rsa` (since version 5.1), `ecdsa` (since version 6.0), and `ed25519` (since version 6.7).

```bash
ssh example.com ssh-keygen -l -f /etc/ssh/ssh_host_ed25519_key.pub | cut -d ' ' -f2
```

## Contributing

We would ❤️ for you to contribute to `nicklasfrahm/scp-action`, pull requests are welcome!

## License

This project is licensed under the [MIT license](./LICENSE.md).
