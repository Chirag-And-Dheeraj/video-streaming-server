const fileForm = document.getElementById("file-form");
const divOutput = document.getElementById("divOutput");
const video = document.getElementById("video");
const titleElement = document.getElementById("title");
const descriptionElement = document.getElementById("description");
const fileError = document.getElementById("fileError");
const titleError = document.getElementById("titleError");
const descriptionError = document.getElementById("descriptionError");

fileForm.addEventListener("submit", (e) => {
    e.preventDefault()

    const title = titleElement.value;
    const description = descriptionElement.value;
    const regex = /^[a-zA-Z0-9\s\-_',.!&():]+$/

    if(!regex.test(title)) {
        titleError.textContent = "Invalid Title";
        titleError.style.display = "block";
        return;
    }

    if(!regex.test(description)) {
        descriptionError.textContent = "Invalid Description";
        descriptionError.style.display = "block";
        return;
    }

    const fileReader = new FileReader();
    const theFile = video.files[0];
    const type = theFile.type
    if(type !== "video/mp4") {
        fileError.textContent = "Only .mp4 files are supported";
        fileError.style.display = "block";
        return;
    }

    titleError.style.display = "none";
    descriptionError.style.display = "none";
    fileError.style.display = "none";

    fileReader.onload = async (ev) => {
        const CHUNK_SIZE = 50000;
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
                const response = await fetch(`${window.ENV.API_URL}/video/`, {
                    method: "POST",
                    headers: {
                        "content-type": "application/octet-stream",
                        "content-length": chunk.length,
                        "file-name": fileName,
                        "file-size": ev.target.result.byteLength,
                        "first-chunk": firstChunk,
                        title: title,
                        description: description,
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
