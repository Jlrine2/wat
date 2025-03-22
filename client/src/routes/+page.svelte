<script>
    import { Button, Input, Modal, Label, Checkbox, Table, TableBody, TableBodyCell, TableBodyRow, TableHead, TableHeadCell, Select } from 'flowbite-svelte';
    import { ClapperboardPlaySolid, OpenDoorSolid, TrashBinSolid } from 'flowbite-svelte-icons';
    let formModal = false;
    /**
     * @type {any}
     */
    let selectedVideo;
    /**
     * @type {any}
     */
    let selectedSubtitle;
    /**
     * @type {string}
     */
    let watchPartyName;

    /**
     * @type {any[]}
     */
    let videos = [];
    /**
     * @type {any[]}
     */
    let subtitles = [];

    async function loadMediaFiles() {
        const response = await fetch('/media');
        if (response.ok) {
            const files = await response.json();
            videos = files
                .filter((/** @type {string} */ file) => file.endsWith('.mp4'))
                .map((/** @type {any} */ file) => ({value: file, name: file}));
            subtitles = files
                .filter((/** @type {string} */ file) => file.endsWith('.vtt'))
                .map((/** @type {any} */ file) => ({value: file, name: file}));
        }
        subtitles.push({value: "None", name: "None"});
    }

    loadMediaFiles();

    /**
     * @type {{ [key: string]: { name: string, video: string, subtitles: string } }}
     */
    let watchParties = {};
    
    async function loadWatchParties() {
        const response = await fetch('/watch-parties');
        if (response.ok) {
            watchParties = await response.json();
        }
    }

    async function handleSubmit() {
        const response = await fetch('/watch-parties', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                name: watchPartyName,
                video: selectedVideo,
                subtitles: selectedSubtitle
            })
        });

        if (response.ok) {
            formModal = false;
            watchPartyName = '';
            selectedVideo = undefined;
            selectedSubtitle = undefined;
            loadWatchParties();
        } else {
            alert('Failed to create watch party');
        }
    }

    /**
     * @param {string} watchPartyId
     */
    async function handleDelete(watchPartyId) {
        const response = await fetch(`/watch-parties?watchPartyId=${watchPartyId}`, {
            method: 'DELETE'
        });

        if (response.ok) {
            loadWatchParties();
        } else {
            alert('Failed to delete watch party');
        }
    }

    loadWatchParties();
  </script>

<div class="flex justify-center bg-gray-700">
    <div class="flex justify-between items-center w-full">
        <div class="flex-grow text-center">
            <p class="text-2xl font-bold dark:text-white">Your Watch Parties</p>
        </div>
    </div>
</div>

<Table>
    <TableHead>
        <TableHeadCell>Name</TableHeadCell>
        <TableHeadCell>Video Name</TableHeadCell>
        <TableHeadCell>Subtitle Name</TableHeadCell>       
        <TableHeadCell>
            <span class="sr-only">Edit</span>
        </TableHeadCell>
    </TableHead>
    <TableBody tableBodyClass="divide-y">
        {#each Object.entries(watchParties) as [watchPartyId, watchParty]}
        <TableBodyRow>
            <TableBodyCell>{watchParty.name}</TableBodyCell>
            <TableBodyCell>{watchParty.video}</TableBodyCell>
            <TableBodyCell>{watchParty.subtitles}</TableBodyCell>
            <TableBodyCell>
              <Button href={`/watch/?videoName=${watchParty.video}&subtitlesName=${watchParty.subtitles}&watchPartyId=${watchPartyId}`}>
                <ClapperboardPlaySolid class="w-5 h-5 me-2" />
                Join
              </Button>
              <Button on:click={() => handleDelete(watchPartyId)}>
                <TrashBinSolid class="w-5 h-5 me-2" />
                Delete
              </Button>
            </TableBodyCell>
        </TableBodyRow>
        {/each}
    </TableBody>
    <tfoot>
        <tr class="font-semibold text-gray-900 dark:text-white">
            <Button on:click={() => (formModal = true)}>Create New Watch Party</Button>
        </tr>
    </tfoot>
  </Table>

  

<Modal bind:open={formModal} size="xs" autoclose={false} class="w-full">
  <form class="flex flex-col space-y-6" on:submit|preventDefault={handleSubmit}>
    <h3 class="mb-4 text-xl font-medium text-gray-900 dark:text-white">Create a new watch party</h3>
    <Label class="space-y-2">
      <span>Name</span>
      <Input type="text" bind:value={watchPartyName} placeholder="Name your watch party something fun!" required />
    </Label>
    <Label class="space-y-2">
      <span>Pick a video to watch</span>
      <Select class="mt-2" items={videos} bind:value={selectedVideo} required />
    </Label>
    <Label class="space-y-2">
      <span>Pick subtitles for the video!</span>
      <Select class="mt-2" items={subtitles} bind:value={selectedSubtitle} required />
    </Label>
    <Button type="submit" class="w-full1">Create Watch Party</Button>
    <div class="text-sm font-medium text-gray-500 dark:text-gray-300">
      Need to upload videos or subtitles? <a href="/manageMedia" class="text-primary-700 hover:underline dark:text-primary-500"> Click here! </a>
    </div>
  </form>
</Modal>