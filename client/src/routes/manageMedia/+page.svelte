<script>
    import { Button, Fileupload, Label, Modal, Table, TableBody, TableBodyCell, TableBodyRow, TableHead, TableHeadCell, Progressbar } from "flowbite-svelte";
    import { DownloadSolid, TrashBinSolid } from "flowbite-svelte-icons";
    let formModal = false;
    let uploadProgress = 0;
    let isUploading = false;
    /**
     * @type {FileList | undefined}
     */
    let selectedFile;

    /**
     * @param {{ preventDefault: () => void; }} event
     */
    async function handleSubmit(event) {
        event.preventDefault();
        if (!selectedFile) return;

        isUploading = true;
        uploadProgress = 0;

        const formData = new FormData();
        formData.append('file', selectedFile[0]);

        try {
            const xhr = new XMLHttpRequest();
            xhr.timeout = 1000*60*5;
            xhr.open('POST', '/media');
            xhr.upload.onprogress = (e) => {
                if (e.lengthComputable) {
                    uploadProgress = (e.loaded / e.total) * 100;
                }
            };

            xhr.onload = () => {
                if (xhr.status === 201) {
                    formModal = false;
                    loadFiles();
                    isUploading = false;
                    uploadProgress = 0;
                }
            };

            xhr.onerror = () => {
                alert('Upload failed');
                isUploading = false;
            };
            xhr.ontimeout = () => {
                alert('Upload Timed out');
                isUploading = false;
            };
            xhr.onabort = () => {
                alert('Upload aborted');
                isUploading = false;
            };


            xhr.send(formData);
        } catch (error) {
            alert('Upload failed');
            isUploading = false;
        }
    }

    /**
     * @param {string} filename
     * @param {string} filetype
     */
    function downloadFile(filename, filetype) {
        const link = document.createElement('a');
        link.href = `/media/${filename}.${filetype}`;
        link.download = `${filename}.${filetype}`;
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
    }

    /**
     * @param {string} filename
     * @param {string} filetype
     */
    async function deleteFile(filename, filetype) {
        try {
            const response = await fetch(`/media?filename=${filename}.${filetype}`, {
                method: 'DELETE'
            });
            if (response.ok) {
                loadFiles();
            } else {
                alert('Failed to delete file');
            }
        } catch (error) {
            alert('Failed to delete file');
        }
    }


    /**
     * @type {{ filename: string; filetype: string; }[]}
     */
    let files = [];
    
    async function loadFiles() {
        const response = await fetch('/media');
        const fileList = await response.json();
        files = fileList.map((/** @type {string} */ file) => {
            const parts = file.split('.');
            return {
                filename: parts.slice(0, -1).join('.'),
                filetype: parts[parts.length - 1]
            };
        });
    }

    loadFiles();
</script>

<div class="flex justify-center bg-gray-700">
    <div class="flex justify-between items-center w-full">
        <div class="flex-grow text-center">
            <p class="text-2xl font-bold dark:text-white">Content Library</p>
        </div>
    </div>
</div>
<Table>
    <TableHead>
        <TableHeadCell>Filename</TableHeadCell>
        <TableHeadCell>File Type</TableHeadCell>      
        <TableHeadCell>
            <span class="sr-only">Delete</span>
        </TableHeadCell>
    </TableHead>
    <TableBody tableBodyClass="divide-y">
        {#each files as file }
        <TableBodyRow>
                <TableBodyCell>{file.filename}</TableBodyCell>
                <TableBodyCell>{file.filetype}</TableBodyCell>
                <TableBodyCell>
                    <Button class="p2!" on:click={() => downloadFile(file.filename, file.filetype)}><DownloadSolid /></Button>
                    <Button class="p2!" on:click={() => deleteFile(file.filename, file.filetype)}><TrashBinSolid /></Button>
                </TableBodyCell>
        </TableBodyRow>
        {/each}
    </TableBody>
    <tfoot>
        <tr class="font-semibold text-gray-900 dark:text-white">
            <Button on:click={() => (formModal = true)}>Upload New File</Button>
        </tr>
    </tfoot>
  </Table>

  <Modal bind:open={formModal} size="xs" autoclose={false} class="w-full">
    <form class="flex flex-col space-y-6" on:submit={handleSubmit}>
        <Label class="pb-2">Upload file</Label>
        <Fileupload name="file" bind:files={selectedFile} />
        {#if isUploading}
            <Progressbar progress={uploadProgress} size="h-3" />
        {/if}
        <Button type="submit" disabled={isUploading}>
            {isUploading ? 'Uploading...' : 'Upload'}
        </Button>
    </form>
  </Modal>