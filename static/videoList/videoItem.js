class VideoItem extends HTMLElement {
  constructor() {
    super();
    this.shadow = this.attachShadow({
      mode: "open"
    });
    const style = document.createElement("style");
    style.textContent = `
          .video-item {
              display: flex;
              align-items: center;
              justify-content: space-between;
              padding: 1rem;
              border-radius: 8px;
              background-color:#262626;
              color: white;
              margin-bottom: 1rem;
              width: 100%;
              cursor: pointer;
          }

          .thumbnail-container {
              width: 200px;
              height: 150px;
              overflow: hidden;
          }

          .thumbnail-container img {
              width: 100%;
              height: 100%;
              object-fit: content;
          }

          .content {
              flex: 0 0 auto; /* Prevent growing/shrinking */
              text-align: left;
              padding: 0 1rem;
          }

          .actions {
              display: flex;
              gap: 0.5rem;
              margin-left: auto; /* Push to right */
              padding: 1em;
          }

          .action-button {
              padding: 0.5rem 0.75rem;
              border-radius: 4px;
              text-decoration: none;
              color: white;
              font-weight: 500;
              cursor: pointer;
              border: none;
              background-color: #141414;
          }

          .delete-modal, .delete {
              background-color: #ff4444;
          }

          .modal {
              display: none;
              position: fixed;
              top: 0;
              left: 0;
              width: 100%;
              height: 100%;
              background-color: rgba(0, 0, 0, 0.5);
              z-index: 1000;
          }

          .modal-content {
              position: absolute;
              top: 50%;
              left: 50%;
              transform: translate(-50%, -50%);
              background-color: #262626;
              padding: 2rem;
              border-radius: 8px;
              text-align: center;
              color: white;
          }

          .modal-actions {
              display: flex;
              justify-content: center;
              gap: 1rem;
              margin-top: 1rem;
          }

          .thumbnail-container {
              position: relative; /* Added for overlay positioning */
              width: 200px;
              height: 150px;
              overflow: hidden;
          }

          .thumbnail-container img {
              width: 100%;
              height: 100%;
              object-fit: content;
          }

          .video-item:hover .thumbnail-container::before {
              content: '';
              position: absolute;
              top: 0;
              left: 0;
              width: 100%;
              height: 100%;
              background-color: rgba(0, 0, 0, 0.5);
              opacity: 1;
              transition: opacity 0.3s ease;
          }

          .thumbnail-container:hover::before {
              opacity: 1;
          }

          .play-button {
              position: absolute;
              top: 50%;
              left: 50%;
              transform: translate(-50%, -50%) scale(0.8);
              opacity: 0;
              transition: all 0.3s ease;
              border: none;
              background: none;
              width: 60px;
              height: 60px;
              display: flex;
              align-items: center;
              justify-content: center;
              cursor: pointer;
              text-align: center;
          }

          .play-button::before {
              content: '▶';
              color: white;
              font-size: 25px;
          }

          .video-item:hover .play-button {
              opacity: 1;
              transform: translate(-50%, -50%) scale(1);
          }
      `;

    const template = document.createElement("template");

    template.innerHTML = `
      <div class="video-item">
          <div class="thumbnail-container">
              <img class="thumbnail" alt="thumbnail"/>
              <button class="play-button"></button>
          </div>
          <div class="content">
              <h3 class="name"></h3>
              <p class="description"></p>
          </div>
          <div class="actions">
              <button class="action-button delete-modal">Delete</button>
          </div>
      </div>
  `;

    const modalTemplate = document.createElement("template");
    modalTemplate.innerHTML = `
      <div class="modal" id="deleteConfirmModal">
          <div class="modal-content">
              <p>Are you sure you want to delete the following file:</p>
              <h3 class="name"></h3>
              <div class="modal-actions">
                  <button class="action-button cancel">Cancel</button>
                  <button class="action-button delete">Delete</button>
              </div>
          </div>
      </div>
  `;

    this.shadow.appendChild(style);

    this.shadow.appendChild(template.content.cloneNode(true));
    this.shadow.appendChild(modalTemplate.content.cloneNode(true));
    this.initializeModal();
    this.initialize();
  }

  initialize() {
    this.shadowRoot.addEventListener("click", (e) => {
      const target = e.target;
      if (target.classList.contains("delete-modal")) {
        e.stopPropagation();
        return;
      }
      if (target.closest(".video-item")) {
        e.preventDefault();
        this.handlePlay();
      }
    });
  }

  initializeModal() {
    const modal = this.shadow.querySelector("#deleteConfirmModal");
    const cancelButton = modal.querySelector(".cancel");
    const deleteButton = modal.querySelector(".delete");
    const fileNameElement = modal.querySelector(".name");

    this.shadowRoot.addEventListener("click", (e) => {
      if (e.target.classList.contains("delete-modal")) {
        e.stopPropagation();
        const fileName = this.getAttribute("name") || "This video";
        fileNameElement.textContent = fileName;
        modal.style.display = "block";
      }
    });

    cancelButton.addEventListener("click", () => {
      modal.style.display = "none";
    });

    deleteButton.addEventListener("click", () => {
      modal.style.display = "none";
      this.handleDelete();
    });

    modal.addEventListener("click", (e) => {
      if (e.target === modal) {
        modal.style.display = "none";
      }
    });
  }

  handlePlay() {
    const videoId = this.getAttribute("video-id");
    window.location.href = `${window.ENV.API_URL}/watch?v=${videoId}`;
  }

  handleDelete() {
    const videoId = this.getAttribute("video-id");
    const deleteButton = this.shadowRoot.querySelector(".delete");

    deleteButton.textContent = "Deleting...";

    fetch(`${window.ENV.API_URL}/video/${videoId}`, {
      method: "DELETE",
    })
      .then((response) => {
        if (response.ok) {
          this.remove();
        } else {
          deleteButton.textContent = "Error";
        }
      })
      .catch((error) => {
        deleteButton.textContent = "Error";
        console.error("Error deleting video:", error);
      });
  }

  static get observedAttributes() {
    return ["name", "description", "thumbnail", "video-id"];
  }

  attributeChangedCallback(name, oldValue, newValue) {
    const element = this.shadow.querySelector(`.${name}`);
    if (element) {
      if (name === "thumbnail") {
        element.src = newValue;
      } else if (name === "name") {
        element.textContent = newValue;
      } else if (name === "description") {
        element.textContent = newValue;
      }
    }
  }

  disconnectedCallback() {
    this.shadowRoot.removeEventListener(
      "click",
      this.shadowRoot.lastEventCallback
    );
  }
}

customElements.define("video-item", VideoItem);
