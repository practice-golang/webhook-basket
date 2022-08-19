FTP, SFTP uploader for Gitea (or Gogs) webhook

## How it works
* Begin work when receive webhook signal from Gitea
* Clone the repository at path in `CLONED_REPO_ROOT` in `webhook-basket.ini`
* Copy cloned files to web-server using `ftp` or `sftp`

## Usage
* Run
```sh
webhook-basket
```
* Help
```sh
webhook-basket -h
```

## Limit
* 1 process per 1 target web-server
    * When ftp/sftp servers are more than 1, run `webhook-basket` with each of `ini` files
    * Example usage
        * Linux
        ```sh
        nohup webhook-basket -ini config_a.ini &
        nohup webhook-basket -ini config_b.ini &
        ...
        ```
        * Windows
        ```powershell
        start webhook-basket.exe -ini config_a.ini
        start webhook-basket.exe -ini config_b.ini
        ...
        ```


## Webhook setting
* Should be set like below picture

![gitea](/doc/gitea.png)

* `webhook-basket` read following data from Gitea sending
```json
{
    "repository": {
        "name": "sample-repo",
        "full_name": "practice-golang/sample-repo",
        "clone_url": "http://localhost:3000/practice-golang/sample-repo.git",
    },
    "pusher": {
        "username": "practice-golang",
        "email": "practice-golang@noreply.example.org",
    },
}
```

* Set target URL like below
```
http://localhost:7749/deploy
```

* Add `deploy-name` parameter when site name is different from name of the repository
```
http://localhost:7749/deploy?deploy-name=othername
```
* Add `deploy-root` parameter when deployment root is different from the path in `webhook-basket.ini`
```
http://localhost:7749/deploy?deploy-root=/home/newroot
```
* All above
```
http://localhost:7749/deploy?deploy-name=othername&deploy-root=/home/newroot
```
* See https://docs.gitea.io/en-us/webhooks


## Trouble Shooting
* Response nothing
    * See https://docs.gitea.io/en-us/config-cheat-sheet/#webhook-webhook
    * Add below option at app.ini of Gitea
    ```ini
    [webhook]
    ALLOWED_HOST_LIST = *
    ```
* Request timeout
    * See https://docs.gitea.io/en-us/config-cheat-sheet/#webhook-webhook
    * Add below option at app.ini of Gitea
    ```ini
    [webhook]
    DELIVER_TIMEOUT = 120
    ```


## Todo
* [ ] Use name mapper feature of `go-ini/ini`


## License

[3-Clause BSD](https://opensource.org/licenses/BSD-3-Clause)
