# Watch Anything Together

Watch anything Together is a simple service that allows you to watch video files in
sync with a group of people across any distance.

## Usage

### Configure
In order to use this app you will need to [create an application](https://discord.com/developers/applications) in Discords developer portal

Copy the example config file and change values as neccessary.

### Setup the server

In order to set up this app we will need a directory on our machine with the media we want to serve. This guide
will assume that there is a `./media` folder with mp4 video content.

```shell
# before we can use we must build the image
docker build -t wat:latest .
docker run -it -p 8080 -v ./media:/media -v ./config/configFile.yaml:/etc/wat/config.yaml wat:latest
```

### Check that media can be served

Open `http://localhost:8080/media/<Name of some video file including .mp4 suffix>`

If that opens a video player things should be working

### Open the Client

Open `http://locahost:8080/`

## Features

This project is very much a work in progress and is not complete, here is a breakdown of the current
implemented/planned features:

✅ Video streaming at "original quality"

✅ Player Sync for Play/Pause Events

✅ Auth with Discord

🚧 Player Sync for Join Events

🚧 Dynamic creation of "watch parties" and support for multiple concurrent "watch parties"

