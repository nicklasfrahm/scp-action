# 🚀 SCP for GitHub Actions

[![Go Report Card](https://goreportcard.com/badge/github.com/nicklasfrahm/scp-action)](https://goreportcard.com/report/github.com/nicklasfrahm/scp-action)

[GitHub Action](https://github.com/features/actions) for copying files and artifacts via SSH.

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
        action: upload
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
        action: download
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

See [action.yml](./action.yml) for more detailed information.

* `host` - ssh host
* `port` - ssh port, default is `22`
* `username` - ssh username, default is `root`
* `timeout` - timeout for ssh to remote host, default is `30s`
* `action_timeout` - timeout for action, default is `10m`
* `key` - content of ssh private key. ex raw content of `~/.ssh/id_rsa`
* `fingerprint` - fingerprint SHA256 of the host public key, see [Using host fingerprint verification](#using-host-fingerprint-verification)
* `source` - a list of files to copy
* `target` - a folder to copy to, default is `.`
* `direction` - either _upload_ or _download_

SSH Proxy Settings:

* `proxy_host` - proxy host
* `proxy_port` - proxy port, default is `22`
* `proxy_username` - proxy username, default is `root`
* `proxy_key` - content of ssh proxy private key.
* `proxy_fingerprint` - fingerprint SHA256 of the proxy host public key, see [Using host fingerprint verification](#using-host-fingerprint-verification)

## Using host fingerprint verification

Setting up SSH host fingerprint verification can help to prevent Person-in-the-Middle attacks. Before setting this up, run the command below to get your SSH host fingerprint. Remember to replace `ed25519` with your appropriate key type (`rsa`, `dsa`, etc.) that your server is using and `example.com` with your host. In modern OpenSSH releases, the _default_ key types to be fetched are `rsa` (since version 5.1), `ecdsa` (since version 6.0), and `ed25519` (since version 6.7).

```bash
ssh example.com ssh-keygen -l -f /etc/ssh/ssh_host_ed25519_key.pub | cut -d ' ' -f2
```

## Contributing

We would ❤️ for you to contribute to `nicklasfrahm/scp-action`, pull requests are welcome!

## License

This project is licensed under the [MIT license](./LICENSE.md).
