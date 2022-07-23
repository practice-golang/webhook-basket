FTP, SFTP uploader for Gitea webhook

## Behavior

* Receive webhook signal from `gitea`
* Clone the repository at path in `CLONED_REPO_ROOT` of `webhook-basket.ini`
* Copy cloned files to web-server using `ftp` or `sftp`


## Webhook setting

* Should be set like below picture

![gitea](/doc/gitea.png)

* Read following data from Gitea sending
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

* When site name is different from name of the repository
```
http://localhost:7749/deploy?deploy-name=othername
```
* When Deployment root is different from the path in `webhook-basket.ini`
```
http://localhost:7749/deploy?deploy-root=/home/newroot
```
* All above
```
http://localhost:7749/deploy?deploy-name=othername&deploy-root=/home/newroot
```
* See https://docs.gitea.io/en-us/webhooks

## Trouble Shooting
* Request timeout
    * Add below option at app.ini of Gitea
    ```ini
    [webhook]
    DELIVER_TIMEOUT = 120
    ```

## License

[3-Clause BSD](https://opensource.org/licenses/BSD-3-Clause)
