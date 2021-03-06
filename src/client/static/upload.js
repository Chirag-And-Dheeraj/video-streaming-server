const uploadVideoButton = document.getElementById("uploadVideoButton");
const divOutput = document.getElementById("divOutput");
const video = document.getElementById("video");
const title = document.getElementById("title");
const description = document.getElementById("description");
// const thumbnail = document.getElementById("thumbnail");

function make_id(length) {
  var result = "";
  var characters =
    "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";
  var charactersLength = characters.length;
  for (var i = 0; i < length; i++) {
    result += characters.charAt(Math.floor(Math.random() * charactersLength));
  }
  return result;
}

uploadVideoButton.addEventListener("click", () => {
  console.log(title.value);
  console.log(description.value);
  console.log("Read and upload button hit!");
  const fileReader = new FileReader();
  const theFile = video.files[0];

  fileReader.onload = async (ev) => {
    const CHUNK_SIZE = 5000000;
    const chunkCount = parseInt(ev.target.result.byteLength / CHUNK_SIZE);
    console.log(chunkCount);
    console.log("Read successfully");
    const fileName = make_id(7) + "_" + theFile.name;
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
        Math.round((sent / ev.target.result.byteLength) * 100, 0) + " %";
    }
    console.log("Successfully sent " + sent + " from the client.");
  };

  fileReader.readAsArrayBuffer(theFile);

  // const thumbnailReader = new FileReader();

  // thumbnailReader.onload = async

  // thumbnailReader.readAsDataURL(thumbnail.files[0]);
});
