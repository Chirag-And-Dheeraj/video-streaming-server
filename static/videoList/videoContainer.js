class VideoList extends HTMLElement {
  constructor() {
    super();
    this.shadow = this.attachShadow({ mode: "open" });
    this.viewMode = localStorage.getItem("viewMode") || "list";

    const style = document.createElement("style");
    style.textContent = `
        :host {
            display: block;
            padding: 1rem;
        }
        .controls {
            display: flex;
            justify-content: center;
            margin: 0.5rem 0 1.5rem 0;
            gap: 0.5rem;
        }
        .view-toggle-button {
            background-color: #262626;
            color: white;
            border: none;
            padding: 0.5rem 0.75rem;
            border-radius: 4px;
            cursor: pointer;
            font-size: 1rem;
            transition: background-color 0.2s ease, border-color 0.2s ease;
        }
        .view-toggle-button:hover {
            background-color: #3a3a3a;
            border-color: #666;
        }
        .view-toggle-button.active {
            background-color: #9b56e1;
            color: white;
        }

        .video-items-container {

        }

        .video-items-container.list-view {
            display: grid;
            grid-template-columns: 50vw;
            gap: 1rem;
        }

        .video-items-container.grid-view {
            display: grid;
            gap: 1.5rem;
            grid-template-columns: repeat(5, 15vw);
        }

        .loading, .error {
            color: #ccc;
            text-align: center;
            padding: 2rem;
            grid-column: 1 / -1;
        }
        .error {
            color: #ff4444;
            font-weight: bold;
        }
    `;

    this.shadow.appendChild(style);

    this.controlsContainer = document.createElement("div");
    this.controlsContainer.className = "controls";
    this.shadow.appendChild(this.controlsContainer);

    this.itemsContainer = document.createElement("div");
    this.itemsContainer.className = `video-items-container ${this.viewMode}-view`;
    this.shadow.appendChild(this.itemsContainer);
  }

  async connectedCallback() {
    this.renderControls();
    await this.loadVideos();
  }

  renderControls() {
    this.controlsContainer.innerHTML = "";

    const listButton = document.createElement("button");
    listButton.className = `view-toggle-button list ${this.viewMode === "list" ? "active" : ""}`;
    listButton.title = "Switch to List View";
    listButton.textContent = "List";
    listButton.addEventListener("click", () => this.setViewMode("list"));
    this.controlsContainer.appendChild(listButton);

    const gridButton = document.createElement("button");
    gridButton.className = `view-toggle-button grid ${this.viewMode === "grid" ? "active" : ""}`;
    gridButton.title = "Switch to Grid View";
    gridButton.textContent = "Grid";
    gridButton.addEventListener("click", () => this.setViewMode("grid"));
    this.controlsContainer.appendChild(gridButton);
  }

  setViewMode(mode) {
    if (this.viewMode === mode) return;

    this.viewMode = mode;
    localStorage.setItem("viewMode", mode);

    this.itemsContainer.classList.remove("list-view", "grid-view");
    this.itemsContainer.classList.add(`${mode}-view`);

    this.shadow
      .querySelector(".view-toggle-button.list")
      .classList.toggle("active", mode === "list");
    this.shadow
      .querySelector(".view-toggle-button.grid")
      .classList.toggle("active", mode === "grid");

    this.itemsContainer.querySelectorAll("video-item").forEach((item) => {
      item.classList.remove("list-mode-item", "grid-mode-item");
      item.classList.add(mode === "grid" ? "grid-mode-item" : "list-mode-item");
    });
  }

  async loadVideos() {
    try {
      this.render(true);
      const response = await fetch(`${window.ENV.API_URL}/video`);
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      const videos = await response.json();
      this.videos = videos;
      this.render();
    } catch (error) {
      console.error("Error loading videos:", error);
      this.error = error.message;
      this.render();
    }
  }

  render(loading = false) {
    const container = this.itemsContainer;
    container.innerHTML = "";

    if (loading) {
      container.innerHTML = '<div class="loading">Loading videos...</div>';
    } else if (this.error) {
      container.innerHTML = `<div class="error">Failed to load videos: ${this.error}</div>`;
    } else if (this.videos && this.videos.length > 0) {
      this.videos.forEach((video) => {
        const videoItem = document.createElement("video-item");
        videoItem.setAttribute("name", video.title || "Untitled Video");
        videoItem.setAttribute("description", video.description || "");
        videoItem.setAttribute("thumbnail", video.thumbnail || "");
        videoItem.setAttribute("video-id", video.id);

        // TODO: status should be validated (keep a set of allowed statuses, I guess?)
        videoItem.setAttribute("status", video.status); // Pass the status here

        videoItem.classList.add(this.viewMode === "grid" ? "grid-mode-item" : "list-mode-item");

        container.appendChild(videoItem);
      });
    } else if (this.videos) {
      container.innerHTML = '<div class="loading">No videos found. Upload some!</div>';
    }
  }
}

customElements.define("video-container", VideoList);
