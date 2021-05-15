# 🚀 SCP for GitHub Actions

[GitHub Action](https://github.com/features/actions) for copying files and artifacts via SSH.

## Usage

Upload local files to remote target:

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
      uses: nicklasfrahm/scp-action@master
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

Download remote files to local target;

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
      uses: nicklasfrahm/scp-action@master
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
        source: |
          path/to/source/a.txt
          path/to/source/b.txt
        target: path/to/target
```

## License

This project is licensed under the [MIT license](./LICENSE.md).
