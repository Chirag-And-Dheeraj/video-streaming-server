const uploadVideoButton = document.getElementById("uploadVideoButton");
const divOutput = document.getElementById("divOutput");
const video = document.getElementById("video");
const title = document.getElementById("title");
const description = document.getElementById("description");
const error = document.getElementById("error");

uploadVideoButton.addEventListener("click", () => {
    console.log(title.value);
    console.log(description.value);
    console.log("Read and upload button hit!");
    const fileReader = new FileReader();
    const theFile = video.files[0];

    console.log(theFile.type);
    const type = theFile.type;
    if(type !== "video/mp4") {
        error.textContent = "Only .mp4 files are supported";
        error.style.display = "block";
        return;
    }

    error.style.display = "none";

    fileReader.onload = async (ev) => {
        const CHUNK_SIZE = 20000000;
        const chunkCount = parseInt(ev.target.result.byteLength / CHUNK_SIZE);
        console.log(chunkCount);
        console.log("Read successfully");

        import("https://jspm.dev/uuid").then(async (uuid) => {
            const uuidv4 = uuid.v4;
            const fileName = uuidv4();
            console.log(fileName);
            let sent = 0;
            for (let chunkID = 0; chunkID < chunkCount + 1; chunkID++) {
                console.log(chunkID);
                let chunk;
                if (chunkID == chunkCount) {
                    chunk = ev.target.result.slice(chunkID * CHUNK_SIZE);
                } else {
                    chunk = ev.target.result.slice(
                        chunkID * CHUNK_SIZE,
                        chunkID * CHUNK_SIZE + CHUNK_SIZE
                    );
                }
                console.log("Chunk byteLength: ", chunk.byteLength);
                sent += chunk.byteLength;
                firstChunk = false;
                if (chunkID == 0) {
                    firstChunk = true;
                }
                // reason for await is we want to wait for server's response and not flood the backend with all requests.
                const response = await fetch("http://127.0.0.1:8000/video/", {
                    method: "POST",
                    headers: {
                        "content-type": "application/octet-stream",
                        "content-length": chunk.length,
                        "file-name": fileName,
                        "file-size": ev.target.result.byteLength,
                        "first-chunk": firstChunk,
                        title: title.value,
                        description: description.value,
                    },
                    body: chunk,
                });

                console.log(await response.text());

                divOutput.textContent =
                    Math.round((sent / ev.target.result.byteLength) * 100, 0) +
                    " %";
            }
            console.log("Successfully sent " + sent + " from the client.");
        });
    };

  fileReader.readAsArrayBuffer(theFile);
});
