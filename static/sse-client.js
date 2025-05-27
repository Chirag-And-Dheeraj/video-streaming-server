const eventSource = new EventSource(`${window.ENV.API_URL}/server-events/`);

eventSource.onopen = () => {
  console.log("SSE connected");
};

eventSource.addEventListener("upload_status", function (event) {
  console.log("Update event received:", event.data);
});

let toaster = new Toast();

eventSource.addEventListener("upload_status", (event) => {
  data = JSON.parse(event.data);
  if (data.upload_status === 1) {
    toaster.success(data.video_title, "Processing completed.");
  } else {
    toaster.error(data.video_title, "Processing failed.");
  }
});
