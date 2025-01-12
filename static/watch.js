window.onload = async () => {
  let video = document.getElementById("video");
  const url = new URL(window.location.href);
  const searchParams = new URLSearchParams(url.search);
  const v = searchParams.get("v");

  let response = await fetch(`${window.ENV.API_URL}/video/${v}`);
  let videoDetails = await response.json();

  let title = document.getElementById("video_title");
  let description = document.getElementById("video_description");

  titleText = document.createTextNode(videoDetails.title);
  descriptionText = document.createTextNode(videoDetails.description);

  title.appendChild(titleText);
  description.appendChild(descriptionText);

  if (Hls.isSupported()) {
    var hls = new Hls();
    hls.loadSource(`${window.ENV.API_URL}/video/${v}/stream/`);
    hls.attachMedia(video);
    hls.on(Hls.Events.MANIFEST_PARSED, function () {
      video.play();
    });
  } else if (video.canPlayType("application/vnd.apple.mpegurl")) {
    video.src = `${window.ENV.API_URL}/video/${v}/`;
    video.addEventListener("loadedmetadata", function () {
      video.play();
    });
  }
};
