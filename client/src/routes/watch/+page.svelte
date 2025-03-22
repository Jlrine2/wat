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


    function playEventHandler(time) {
        const message = JSON.stringify({"event": "play", "time": time});
        syncSocket.send(message)
        console.log("sending ws message: ", message)
    }
    function pauseEventHandler(time) {
        const message = JSON.stringify({"event": "pause", "time": time});
        syncSocket.send(message)
        console.log("sending ws message: ", message)
    }
    let lastSeekTime = 0
    function seekEventHandler(time) {
        // Only send seek events if user manually sought the video
        // and not from programmatic seeks from other events
        if (Math.abs(time - lastSeekTime) > 0.5) {
            const message = JSON.stringify({"event": "seek", "time": time});
            syncSocket.send(message)
            console.log("sending ws message: ", message)
        }
        lastSeekTime = time;
    }
    syncSocket.onmessage = (event) => {
        let video = document.getElementById("video-player")
        const msg = JSON.parse(event.data)
        console.log("recieved ws message: ", msg)
        const eventType = msg.event
        const time = msg.time
        switch (eventType) {
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

{#if authConfig.isAuthenticated}
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

{:else} 
    <Alert color="red">
        <span class="font-medium">Access Denied</span>
        You must be logged in to watch
    </Alert>
{/if}
