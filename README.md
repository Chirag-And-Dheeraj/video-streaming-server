# Dekho: A Journey into Audio/Video Streaming

> **Announcement: 📢** Dekho is back! After being a dead for nearly two years, we're reviving this project to explore the fundamentals of audio and video streaming once again. You can track our progress on the revival of the project in this milestone [here](https://github.com/Chirag-And-Dheeraj/video-streaming-server/milestone/1).

## Overview

Dekho is a research and study project aimed at understanding and mastering the intricacies of audio and video streaming. Our primary focus was on implementing the [HLS (HTTP Live Streaming)](https://developer.apple.com/streaming) protocol to build an on-demand live streaming server.

However, expect things to change in the upcoming versions while we shape this into something more refined.

## Setup Instructions

- Install the `make` utility because we have created a **_[Makefile](Makefile)_** to ease the setup
- Clone the repository
- Create an **_[Appwrite storage bucket](https://appwrite.io/docs/products/storage)_**
  - Make sure that you make a note of `APPWRITE_KEY`, `APPWRITE_PROJECT_ID` and the `BUCKET_ID`.
- Now we are moving to the dependency installation steps.
- This project needs `go version go1.22.2`, `ffmpeg` utility and `psql (PostgreSQL) 16.3` database.
- To install these, and the Go dependencies, run `make install`.
  - The `make install` takes care of installing go, ffmpeg and psql dependencies.
  - Run `make start-postgres` to start the postgres service.
  - Create a database user by running this command `sudo -u postgres createuser -s username_here -P`
  - `-P` will prompt for a password.
  - Enter the `psql` shell by running this command `sudo -u postgres psql` and create a database by running `CREATE DATABASE <database_name>;`.
  - Create a `.env` file using **_[template.env](template.env)_** as a reference.
- Once the dependencies are installed successfully, run `make cleanstart`.
- This command will create all the necessary folders and start the server on `http://127.0.0.1:8000`
- If you just want to run the server, run: `make start`
- If you just want to clean up, run: `make clean`

## Technologies Used

- **Server:** Go
- **Database:** SQLite
- **Storage:** [Appwrite Storage](https://appwrite.io/docs/products/storage)
- **Video Processing:** [FFMPEG](https://ffmpeg.org) for breaking down videos into .ts chunks
- **Video Player:** [HLS.js](https://github.com/video-dev/hls.js)
- **Frontend:** HTML, CSS, JS

## Project Architecture

<details>
<summary>View Architecture Diagram</summary>

![Architecture Diagram](https://user-images.githubusercontent.com/52416311/167314446-c991f74d-e579-438d-a6ad-b65b7e721e7f.png)

</details>

## Key Resources

Our journey has been greatly enriched by the insights and guidance from various resources. A pivotal article that set us on the right path is "[Learning the basics of video streaming with Golang](https://www.rohitmundra.com/video-streaming-server)" by [Rohit Mundra](https://twitter.com/brohit3).

For a comprehensive list of resources that have been instrumental in our learning and development process, please refer to our [documentation](https://github.com/Chirag-And-Dheeraj/video-streaming-server/blob/main/documentation/video-streaming-project-stuff/links.md).

## Getting Involved

We welcome any queries or contributions to the project. If you have any questions or suggestions, please feel free to reach out to us:

- **Dheeraj Lalwani:** [Twitter](https://twitter.com/DhiruCodes)
- **Chirag Lulla:** [Twitter](https://twitter.com/_chiraglulla_)

## Updates and Blog Posts

Stay tuned for upcoming blog posts and updates on our progress. If you haven't heard from us in a while, feel free to bug us about it!
