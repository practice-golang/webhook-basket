# A shoddy ~~deployer~~ copier :-D

## Process
* Receive webhook signal from `gitea`
* Clone to temporary path of `webhook-basket` using `git`
* Copy to `web-server` using `ssh`, `scp`

## Using
* `go-scp` for file copying
* `crypto/ssh` for directory creation

## Requirement
* Git
* SSH server
    * linux : https://www.google.com/search?q=linux+install+ssh&oq=linux+install+ssh
    * windows : https://winscp.net/eng/docs/guide_windows_openssh_server#installing_sftp_ssh_server

## Build
```powershell
mingw32-make.exe
```
or
```bash
make
```

## Set target repositories
See `ini/webhook-basket.ini`

## Using webhook
See https://docs.gitea.io/en-us/webhooks

## Todo
* Block reread ini, push same repo when goroutine running
- [x] Change route
- [x] .git directory eatting bug
