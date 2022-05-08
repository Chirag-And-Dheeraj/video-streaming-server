# Dekho
A research/study project for understanding and learning the fundamentals of audio/video streaming.

We've attempted to use [HLS](https://developer.apple.com/streaming) Protocol for building an on-demand live streaming live streaming server.

## An article that really helped us a lot
[Learning the basics of video streaming with Golang](https://www.rohitmundra.com/video-streaming-server) written by [Rohit Mundra](https://twitter.com/brohit3).

A lot of the resources that helped us are documented with links [here](https://github.com/Chirag-And-Dheeraj/video-streaming-server/blob/main/documentation/video-streaming-project-stuff/links.md).

## Architecture diagram

![image](https://user-images.githubusercontent.com/52416311/167314446-c991f74d-e579-438d-a6ad-b65b7e721e7f.png)

## Tech stack used:
- Go (Server)
- SQLite (Database)
- [FFMPEG](https://ffmpeg.org) (For breaking down the video into .ts chunks)
- [HLS.js](https://github.com/video-dev/hls.js) (Video player)
- HTML, CSS, JS

If you have any queries about the project, please reach out to me: [Dheeraj Lalwani](https://twitter.com/DhiruCodes) or [Chirag Lulla](https://twitter.com/_chiraglulla_).

You can expect a blog about the project soon. If we fail, please bug us about it so that we don't forget about it in the grand scheme of things.
