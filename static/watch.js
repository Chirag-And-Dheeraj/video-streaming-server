var video = document.getElementById("video");
const url = new URL(window.location.href)
const searchParams = new URLSearchParams(url.search);
const v = searchParams.get("v")
if (Hls.isSupported()) {
  var hls = new Hls();
  hls.loadSource(
    `http://127.0.0.1:8000/video/${v}/`
  );
  hls.attachMedia(video);
  hls.on(Hls.Events.MANIFEST_PARSED, function () {
    video.play();
  });
} else if (video.canPlayType("application/vnd.apple.mpegurl")) {
  video.src =
    `http://127.0.0.1:8000/video/${v}/`;
  video.addEventListener("loadedmetadata", function () {
    video.play();
  });
}
