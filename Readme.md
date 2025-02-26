# Watch Anything Together

Watch anything Together is a simple service that allows you to watch video files in 
sync with a group of people across any distance.

## Usage
### Setup the server
In order to set up this app we will need a directory on our machine with the media we want to serve. This guide
will assume that there is a `./media` folder with mp4 video content.

```shell
# before we can use we must build the image
docker built -t wat:latest .
docker run -it -p 8080 -v ./media:/media wat:latest -p 8080 -m /media
```

### Check that media can be served
Open `http://localhost:8080/media/<Name of some video file including .mp4 suffix>`

If that opens a video player things should be working

### Go to the watch party
Open `http://locahost:8080/video?videoName=<Name of some video file including .mp4 suffix>`
When opening your client will be registered with the watch party automatically and any play/pause events
done by other members will be synced with your player and vice versa.


## Features
This project is very much a work in progress and is not complete, here is a breakdown of the current
implemented/planned features:

✅ Video streaming at "original quality"

✅ Player Sync for Play/Pause Events

🚧 Player Sync for Join Events

🚧 Dynamic creation of "watch parties" and support for multiple concurrent "watch parties"

🚧 Auth, UI, UX, etc