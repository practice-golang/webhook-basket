FTP, SFTP uploader for Gitea (or Gogs) webhook

## How it works
* Begin work when receive webhook signal from Gitea
* Clone the repository at path under `CLONED_REPO_ROOT` in `webhook-basket.ini`
* Copy cloned files to target web-server via `ftp` or `sftp`

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
    * When target ftp/sftp servers are more than 1, run `webhook-basket` with each of `ini` files
    * Listening port of each `ini` files must be set different number
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
    * See https://docs.gitea.io/en-us/webhooks
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


## Trouble Shooting
* Response nothing, Request timeout
    * See https://docs.gitea.io/en-us/config-cheat-sheet/#webhook-webhook
    * Append or modify the following options in `app.ini` of Gitea
    ```ini
    [webhook]
    ALLOWED_HOST_LIST = *
    DELIVER_TIMEOUT = 120
    ```


## Todo
* [ ] Auth(Secret) header
* [ ] Exclude files parameter
* [ ] Remove post-sample which is not required


## License

[3-Clause BSD](https://opensource.org/licenses/BSD-3-Clause)
