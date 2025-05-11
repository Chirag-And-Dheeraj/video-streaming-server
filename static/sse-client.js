const eventSource = new EventSource(`${window.ENV.API_URL}/server-events/`);

eventSource.onopen = () => {
  console.log("SSE connected");
};

eventSource.addEventListener("upload_status", function (event) {
  console.log("Update event received:", event.data);
});
