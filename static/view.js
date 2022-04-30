window.onload = async () => {
  console.log("Page loaded, bitch...");
  const response = await fetch("http://127.0.0.1:8000/video", {
    method: "GET",
  });
  let videos = await response.json();
  console.log(videos);
};
