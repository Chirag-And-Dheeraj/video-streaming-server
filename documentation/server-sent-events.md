# status updates for videos

- upload a video
- subscribe to SSE for that video
- keep an eye on the backend for video status and when uploaded completely, send an event to client
- if we get that event on upload page or any other page show a toast
- if we get that event on listing page, update that particular video item

# final approach for SSE

- One, single user-level EventStream

every user subscribes to that SSE

events will be separated upon name spaces

- if event is: video_upload:
  - will contain video_id
  - will contain video_processing_status
  - and client will accordingly do stuff


- https://packagemain.tech/p/implementing-server-sent-events-in-go
- https://medium.com/@rian.eka.cahya/server-sent-event-sse-with-go-10592d9c2aa1
- https://www.youtube.com/watch?v=nvijc5J-JAQ
- Cookie stuff: https://chatgpt.com/c/67ea6365-d6d4-8008-8f76-4f3f54aa3777
- https://www.freecodecamp.org/news/how-to-implement-server-sent-events-in-go/
- https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events/Using_server-sent_events
- https://threedots.tech/post/live-website-updates-go-sse-htmx/