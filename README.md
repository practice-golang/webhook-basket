FTP, SFTP uploader for Gitea webhook

## How it works

* Begin work when receive webhook signal from Gitea (or Gogs or Github or Gitlab)
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

* Build from source
    * Windows
    ```powershell
    # Because of pkg name, not work yet. just download this repo
    # go get github.com/practice-golang/webhook-basket
    build.ps1
    ```
    * Linux
    ```sh
    # Because of pkg name, not work yet. just download this repo
    # go get github.com/practice-golang/webhook-basket
    build.sh
    ```


## Behavior

* Run 1 process per 1 target web-server
    * If target ftp/sftp servers are more than 1, several `webhook-basket` should be run with each of `ini` files
    * Listening port in each `ini` files have to be set different number
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

* There's no delete & flush target ftp directory. Only appending or overwriting


## Webhook setting

* Should be set like below picture. Also `secret` should be added if set it in `ini` file

![gitea](/doc/gitea.png)

* Set target URL like below
```powershell
http://localhost:7749/deploy
# Files will be copied to DEPLOYMENT_ROOT/json_repository_name
# DEPLOYMENT_ROOT variable is in ini file
# json_repository_name is in requested json data
```
* Add `deploy-name` parameter when site name is different from name of the repository
```powershell
http://localhost:7749/deploy?deploy-name=othername
# Files will be copied to DEPLOYMENT_ROOT/othername
# DEPLOYMENT_ROOT variable is in ini file
```
* Add `deploy-root` parameter when deployment root is different from the path in `webhook-basket.ini`
```powershell
http://localhost:7749/deploy?deploy-root=/home/newroot
# Files will be copied to /home/newroot/json_repository_name
# json_repository_name is in requested json data
```
* All above
```powershell
http://localhost:7749/deploy?deploy-name=othername&deploy-root=/home/newroot
# Files will be copied to /home/newroot/othername/
```
* When root is root(/)
```powershell
http://localhost:7749/deploy?deploy-name=othername&deploy-root=/
# Files will be copied to /othername/
```


* `webhook-basket` use only following data from Gitea sending
    * See about webhook payload
        * https://docs.gitea.io/en-us/webhooks
        * https://docs.github.com/en/developers/webhooks-and-events/webhooks/webhook-events-and-payloads
        * https://gogs.io/docs/features/webhook
        * https://docs.gitlab.com/ee/user/project/integrations/webhook_events.html
```json
{
    "repository": {
        "name": "sample-repo",
        "full_name": "practice-golang/sample-repo",
        "clone_url": "http://localhost:3000/practice-golang/sample-repo.git",
    }
}
```

* Secret - `webhook-basket` read one header of following signatures which generated with `secret` in `ini` and in `secret form` in `webhook`.
```
X-Gitea-Signature: 2f8e..
X-Gogs-Signature: 2f8e..
X-Hub-Signature-256: sha256=2f8e..
```


## Exclude file(s) when upload to ftp/sftp
* Add `.wbignore` to the target repository
* Syntax is similar with `.gitignore`.
```.gitignore
/.git
.gitignore
.wbignore

README.md
```


## Delete temporary repository root
* Send following request
```sh
DELETE uri/deploy-root
Secret secRET12345 # optional. If not set in ini, this will be ignored
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


## API
* `GET health` - Health check
* `POST deploy` - Receive webhook and deploy
* `DELETE repos-root` - Delete temporary repository root
* See `requests.http`


## License

[3-Clause BSD](https://opensource.org/licenses/BSD-3-Clause)
