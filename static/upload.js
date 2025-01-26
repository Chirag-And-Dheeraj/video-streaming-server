const fileForm = document.getElementById("file-form");
const divOutput = document.getElementById("divOutput");
const video = document.getElementById("video");
const titleElement = document.getElementById("title");
const descriptionElement = document.getElementById("description");
const fileError = document.getElementById("fileError");
const titleError = document.getElementById("titleError");
const descriptionError = document.getElementById("descriptionError");
const uploadVideoButton = document.getElementById("uploadVideoButton");

function checkFileType(file) {
  const supportedTypes = JSON.parse(
    localStorage.getItem("SUPPORTED_FILE_TYPES")
  );

  if (
    !supportedTypes.some(
      (supportedType) => supportedType.file_type === file.type
    )
  ) {
    console.log(file.type);
    const supportedExtensions = supportedTypes
      .map((type) => type.file_extension)
      .join(", ");

    fileError.textContent = `Only ${supportedExtensions} files are supported`;
    fileError.style.display = "block";
    return false;
  }
  return true;
}

fileForm.addEventListener("submit", async (e) => {
  e.preventDefault();

  const title = titleElement.value;
  const description = descriptionElement.value;
  const regex = /^[a-zA-Z0-9\s\-_',.!&():]+$/;

  if (!regex.test(title)) {
    titleError.textContent = "Invalid Title";
    titleError.style.display = "block";
    return;
  }

  if (!regex.test(description)) {
    descriptionError.textContent = "Invalid Description";
    descriptionError.style.display = "block";
    return;
  }

  const fileReader = new FileReader();
  const theFile = video.files[0];
  const size = theFile.size;

  console.log(`file type = ${theFile.type}`);

  if (!checkFileType(theFile)) {
    return;
  }

  const sizeLimit = localStorage.getItem("FILE_SIZE_LIMIT");
  if (size > sizeLimit) {
    fileError.textContent = `File size is greater than ${(
      sizeLimit /
      (1024 * 1024)
    ).toFixed(2)} MB`;
    fileError.style.display = "block";
    return;
  }

  titleError.style.display = "none";
  descriptionError.style.display = "none";
  fileError.style.display = "none";

  uploadVideoButton.disabled = true;

  fileReader.onload = async (ev) => {
    const CHUNK_SIZE = 50000;
    const chunkCount = Math.ceil(ev.target.result.byteLength / CHUNK_SIZE);
    console.log(chunkCount);
    console.log("Read successfully");

    import("https://jspm.dev/uuid").then(async (uuid) => {
      const uuidv4 = uuid.v4;
      const fileName = uuidv4();
      console.log(fileName);
      let sent = 0;
      let chunkID;
      for (chunkID = 0; chunkID < chunkCount; chunkID++) {
        console.log(chunkID);
        let chunk;
        if (chunkID == chunkCount - 1) {
          chunk = ev.target.result.slice(chunkID * CHUNK_SIZE);
        } else {
          chunk = ev.target.result.slice(
            chunkID * CHUNK_SIZE,
            chunkID * CHUNK_SIZE + CHUNK_SIZE
          );
        }
        console.log("Chunk byteLength: ", chunk.byteLength);
        sent += chunk.byteLength;
        const firstChunk = chunkID === 0;

        const response = await fetch(`${window.ENV.API_URL}/video/`, {
          method: "POST",
          headers: {
            "content-type": "application/octet-stream",
            "content-length": chunk.length,
            "file-name": fileName,
            "file-size": ev.target.result.byteLength,
            "first-chunk": firstChunk.toString(),
            title: title,
            description: description,
          },
          body: chunk,
        });

        console.log(await response.text());

        divOutput.textContent = `${Math.round(
          (sent / ev.target.result.byteLength) * 100,
          0
        )} %`;
      }

      if (chunkID >= chunkCount) {
        divOutput.append(
          ". Your video will be available in few minutes on List Files page."
        );
      }
      console.log(`Successfully sent ${sent} bytes from the client.`);
    });
  };

  fileReader.readAsArrayBuffer(theFile);
});
