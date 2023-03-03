# SP Bday Updater

Did you ever want to have birthday every day on https://www.spieleplanet.eu/? Now, you can automate setting your birthday to the current day and celebrate with you SP-friends on a daily basis!

## Features

* Updates your birthday everyday on https://www.spieleplanet.eu/ at midnight
* Set your individual timezone

## Usage

### Binary

* Download binary from [Github](https://github.com/yawn77/spbdayupdater/releases)
* Set environment variables
  * `SP_USERNAME=<your username>`
  * `SP_PASSWORD=<your password>`
  * `TZ=<your timezone>` (e.g., America/Los_Angeles)

### Docker

* Pull latest docker image:
```
docker pull yawn77/spbdayupdater:latest
```

* Run docker container:
```
docker run --rm \
  -e SP_USERNAME=<your username> \
  -e SP_PASSWORD=<your password> \
  -e TZ=<your timezone> \
  yawn77/spbdayupdater:latest
```

## License

[MIT](https://choosealicense.com/licenses/mit/)
