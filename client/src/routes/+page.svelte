<script>
    import { Button, Input, Modal, Label, Checkbox, Table, TableBody, TableBodyCell, TableBodyRow, TableHead, TableHeadCell, Select } from 'flowbite-svelte';
    import { ClapperboardPlaySolid, TrashBinSolid } from 'flowbite-svelte-icons';
    let formModal = false;
    /**
     * @type {any}
     */
    let selectedVideo;
    /**
     * @type {any}
     */
    let selectedSubtitle;

    let videos = [{value: "1", name: "Video 1"}, {value: "2", name: "Video 2"}, {value: "3", name: "Video 3"}];
    let subtitles = [{value: "1", name: "Subtitle 1"}, {value: "2", name: "Subtitle 2"}, {value: "3", name: "Subtitle 3"}];

    let watchParties = [{name: "Watch Party 1", video: "Video 1", subtitle: "Subtitle 1"}, {name: "Watch Party 2", video: "Video 2", subtitle: "Subtitle 2"}, {name: "Watch Party 3", video: "Video 3", subtitle: "Subtitle 3"}];
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
        {#each watchParties as watchParty}
        <TableBodyRow>
                <TableBodyCell>{watchParty.name}</TableBodyCell>
                <TableBodyCell>{watchParty.video}</TableBodyCell>
                <TableBodyCell>{watchParty.subtitle}</TableBodyCell>
                <TableBodyCell>
                  <Button><ClapperboardPlaySolid class="w-5 h-5 me-2" /> Join</Button>
                  <Button class="p2!" alt="Delete Watch Party"><TrashBinSolid /></Button>
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
  <form class="flex flex-col space-y-6" action="#">
    <h3 class="mb-4 text-xl font-medium text-gray-900 dark:text-white">Create a new watch party</h3>
    <Label class="space-y-2">
      <span>Name</span>
      <Input type="text" name="name" placeholder="Name your watch party something fun!" required />
    </Label>
    <Label class="space-y-2">
      <span>Pick a video to watch</span>
      <Select class="mt-2" items={videos} bind:value={selectedVideo} />
    </Label>
    <Label class="space-y-2">
      <span>Pick subtitles for the video!</span>
      <Select class="mt-2" items={subtitles} bind:value={selectedSubtitle} />
    </Label>
    <Button type="submit" class="w-full1">Create Watch Party</Button>
    <div class="text-sm font-medium text-gray-500 dark:text-gray-300">
      Need to upload videos or subtitles? <a href="/manageMedia" class="text-primary-700 hover:underline dark:text-primary-500"> Click here! </a>
    </div>
  </form>
</Modal>