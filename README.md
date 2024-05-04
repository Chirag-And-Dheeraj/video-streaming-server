# Dekho: A Journey into Audio/Video Streaming

> **Announcement: 📢** Dekho is back! After being a dead for nearly two years, we're reviving this project to explore the fundamentals of audio and video streaming once again. You can track our progress on the revival of the project in this milestone [here](https://github.com/Chirag-And-Dheeraj/video-streaming-server/milestone/1).

## Overview

Dekho is a research and study project aimed at understanding and mastering the intricacies of audio and video streaming. Our primary focus was on implementing the [HLS (HTTP Live Streaming)](https://developer.apple.com/streaming) protocol to build an on-demand live streaming server.

However, expect things to change in the upcoming versions while we shape this into something more refined.

## Key Resources

Our journey has been greatly enriched by the insights and guidance from various resources. A pivotal article that set us on the right path is "[Learning the basics of video streaming with Golang](https://www.rohitmundra.com/video-streaming-server)" by [Rohit Mundra](https://twitter.com/brohit3).

For a comprehensive list of resources that have been instrumental in our learning and development process, please refer to our [documentation](https://github.com/Chirag-And-Dheeraj/video-streaming-server/blob/main/documentation/video-streaming-project-stuff/links.md).

## Project Architecture

<details>
<summary>View Architecture Diagram</summary>

![Architecture Diagram](https://user-images.githubusercontent.com/52416311/167314446-c991f74d-e579-438d-a6ad-b65b7e721e7f.png)

</details>

## Technologies Used

- **Server:** Go
- **Database:** SQLite
- **Storage:** Migrating from [Deta Drive](https://deta.space/docs/en/build/reference/drive/) to [Appwrite Storage](https://appwrite.io/docs/products/storage)
- **Video Processing:** [FFMPEG](https://ffmpeg.org) for breaking down videos into .ts chunks
- **Video Player:** [HLS.js](https://github.com/video-dev/hls.js)
- **Frontend:** HTML, CSS, JS

## Getting Involved

We welcome any queries or contributions to the project. If you have any questions or suggestions, please feel free to reach out to us:

- **Dheeraj Lalwani:** [Twitter](https://twitter.com/DhiruCodes)
- **Chirag Lulla:** [Twitter](https://twitter.com/_chiraglulla_)

## Updates and Blog Posts

Stay tuned for upcoming blog posts and updates on our progress. If you haven't heard from us in a while, feel free to bug us about it!
