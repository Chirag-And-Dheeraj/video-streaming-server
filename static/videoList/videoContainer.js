class VideoList extends HTMLElement {
  constructor() {
    super();
    this.shadow = this.attachShadow({ mode: "open" });

    const style = document.createElement("style");
    style.textContent = `  
        .video-container {
            display: grid;
            gap: 1rem;
            padding: 1rem;
        }
        .loading {
            color: #666;
            text-align: center;
        }
        .error {
            color: #ff4444;
            text-align: center;
        }
    `;

    this.shadow.appendChild(style);
    this.container = document.createElement("div");
    this.container.className = "video-container";
    this.shadow.appendChild(this.container);
  }

  async connectedCallback() {
    await this.loadVideos();
  }

  async loadVideos() {
    try {
      this.render(true);
      const response = await fetch(`${window.ENV.API_URL}/video`);
      const videos = await response.json();
      this.videos = videos;
      this.render();
    } catch (error) {
      this.error = error.message;
      this.render();
    }
  }

  render(loading = false) {
    let content = "";
    if (loading) {
      content = '<div class="loading">Loading videos...</div>';
    } else if (this.error) {
      content = `<div class="error">${this.error}</div>`;
    } else if (this.videos) {
      content = this.videos
        .map(
          (video) => `
                <video-item
                    name="${video.title}"
                    description="${video.description}"
                    thumbnail="${video.thumbnail}"
                    video-id="${video.id}"
                ></video-item>
            `
        )
        .join("");
    }
    this.container.innerHTML = content;
  }
}

customElements.define("video-container", VideoList);
