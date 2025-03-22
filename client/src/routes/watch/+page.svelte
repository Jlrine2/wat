<script>
// @ts-nocheck

    import { authConfig } from "$lib/auth";
    import { Alert } from "flowbite-svelte";
    import { onDestroy } from "svelte";

    const syncSocket = new WebSocket(`/ws`);
    const urlParams = new URLSearchParams(window.location.search)
    const videoName = urlParams.get("videoName")
    const subtitleName = urlParams.get("subtitleName")
    const title = document.getElementById("title")
    const watchPartyId = urlParams.get("watchPartyId")
    let joined = false

    syncSocket.onopen = async () => {
        await syncSocket.send(JSON.stringify({"event": "join", "watchPartyId": watchPartyId}));
    }


    async function playEventHandler(time) {
        const message = JSON.stringify({"event": "play", "time": time, "watchPartyId": watchPartyId});
        await syncSocket.send(message)
        console.log("sending ws message: ", message)
    }
    async function pauseEventHandler(time) {
        const message = JSON.stringify({"event": "pause", "time": time, "watchPartyId": watchPartyId});
        await syncSocket.send(message)
        console.log("sending ws message: ", message)
    }
    let lastSeekTime = 0
    async function seekEventHandler(time) {
        // Only send seek events if user manually sought the video
        // and not from programmatic seeks from other events
        if (Math.abs(time - lastSeekTime) > 0.5) {
            const message = JSON.stringify({"event": "seek", "time": time, "watchPartyId": watchPartyId});
            await syncSocket.send(message)
            console.log("sending ws message: ", message)
        }
        lastSeekTime = time;
    }

    // Handle current time request from new joiners
    function joinResponse() {
        let video = document.getElementById("video-player");
        if (video) {
            const message = JSON.stringify({
                "event": "join_response",
                "time": video.currentTime,
                "watchPartyId": watchPartyId,
                "isPlaying": !video.paused
            });
            syncSocket.send(message);
        }
    }

    syncSocket.onmessage = (event) => {
        let video = document.getElementById("video-player")
        const msg = JSON.parse(event.data)
        console.log("recieved ws message: ", msg)
        const eventType = msg.event
        const time = msg.time
        if (msg.watchPartyId !== watchPartyId) {
            return;
        }
        switch (eventType) {
            case "join":
                joinResponse();
                break;
            case "join_response":
                if (!joined) {
                    video.currentTime = time;
                    if (msg.isPlaying) {
                        video.play();
                    } else {
                        video.pause();
                    }
                    joined = true;
                }
                break;
            case "play":
                video.currentTime = time;
                video.play();
                break;
            case "pause": 
                video.currentTime = time;
                video.pause();
                break;
            case "seek":
                lastSeekTime = time;
                video.currentTime = time;
                break;
        }
    }

    onDestroy(() => {
        syncSocket.close();
    });
</script>

<div class="flex justify-center">
    <video id="video-player"
        width="90%" 
        height="auto" 
        controls
        class="rounded-lg shadow-lg my-8 max-w-6xl mx-auto bg-black"
        preload="metadata"
        on:play={(e) => playEventHandler(e.target.currentTime)}
        on:pause={(e) => pauseEventHandler(e.target.currentTime)}
        on:seeked={(e) => seekEventHandler(e.target.currentTime)}
    >
        <source src="/media/{videoName}" type="video/mp4">
        <track kind="captions" src="/media/{subtitleName}">
        Your browser does not support the video tag.
    </video>
</div>
