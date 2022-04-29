const uploadVideoButton = document.getElementById("uploadVideoButton");
const divOutput = document.getElementById("divOutput");
const video = document.getElementById("video");

uploadVideoButton.addEventListener("click", () => {
  console.log("Read and upload button hit!");
  const fileReader = new FileReader();
  const theFile = video.files[0];

  fileReader.onload = async (ev) => {
    const CHUNK_SIZE = 10000000;
    const chunkCount = parseInt(ev.target.result.byteLength / CHUNK_SIZE);
    console.log(chunkCount);
    console.log("Read successfully");
    const fileName = Math.random() * 1000 + theFile.name;
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
      // reason for await is we want to wait for server's response and not flood the backend with all requests.
      const response = await fetch("http://127.0.0.1:8000/video", {
        method: "POST",
        headers: {
          "content-type": "application/octet-stream",
          "content-length": chunk.length,
          "file-name": fileName,
          "file-size": ev.target.result.byteLength,
        },
        body: chunk,
      });

      console.log(await response.text());

      divOutput.textContent =
        Math.round((sent / ev.target.result.byteLength) * 100, 0) + " %";
    }
    console.log("Successfully sent " + sent + " from the client.");
  };

  fileReader.readAsArrayBuffer(theFile);
});
