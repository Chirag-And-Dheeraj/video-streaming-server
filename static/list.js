window.onload = async () => {
  const response = await fetch(`${window.ENV.API_URL}/video`, {
    method: "GET",
  });

  let videos = await response.json();

  console.log(videos);

  let videoListSection = document.getElementById("video_list");

  for (let i = 0; i < videos.length; i++) {
    let videoListRow = document.createElement("section");
    let videoTitle = document.createElement("h2");
    let videoDescription = document.createElement("p");
    let videoLink = document.createElement("a");
    let videoTitleText = document.createTextNode(videos[i].title);
    let videoDescriptionText = document.createTextNode(videos[i].description);

    videoLink.setAttribute(
      "href",
      `${window.ENV.API_URL}/watch?v=${videos[i].id}`
    );
    videoLink.textContent = "Play";

    videoTitle.appendChild(videoTitleText);

    videoDescription.appendChild(videoDescriptionText);

    videoListRow.appendChild(videoTitle);
    videoListRow.appendChild(videoDescription);
    videoListRow.appendChild(videoLink);

    videoListSection.appendChild(videoListRow);
  }
};
