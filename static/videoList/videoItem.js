class VideoItem extends HTMLElement {
  constructor() {
    super();
    this.shadow = this.attachShadow({ mode: "open" });
    const style = document.createElement("style");

    style.textContent = `
          :host(.list-mode-item) .video-item {
              display: flex;
              align-items: center;
              justify-content: space-between;
              width: 100%;
              background: #262626;
              margin: 6px;
              padding: 12px;
              border-radius: 12px;
          }
           :host(.list-mode-item) .thumbnail-container {
              width: 200px;
              height: 120px;
              flex-shrink: 0;
              border-radius: 12px;
          }
          :host(.list-mode-item) .content {
              flex-grow: 1;
              text-align: left;
              padding: 0 1rem;
              overflow: hidden;
          }
          :host(.list-mode-item) .actions {
              margin-left: auto;
              padding: 1em;
              flex-shrink: 0;
          }

          :host(.grid-mode-item) .video-item {
              display: flex;
              flex-direction: column;
              align-items: stretch;
              height: 100%;
              overflow: hidden;
              background-color: #262626;
              border-radius: 8px;
          }
          :host(.grid-mode-item) .thumbnail-container {
              width: 100%;
              aspect-ratio: 16 / 9;
              height: auto;
              width: 250px;
              flex-shrink: 0;
              border-radius: 8px 8px 0 0;
          }
           :host(.grid-mode-item) .thumbnail-container img {
                object-fit: cover;
                border-radius: 8px 8px 0 0;
           }
          :host(.grid-mode-item) .content {
              padding: 0.75rem 0.75rem 0.5rem 0.75rem;
              text-align: left;
              overflow: hidden;
              flex-grow: 1;
              display: flex;
              flex-direction: column;
          }
           :host(.grid-mode-item) .content .name {
               white-space: nowrap;
               overflow: hidden;
               text-overflow: ellipsis;
               margin-bottom: 0.25rem;
               font-size: 1.0em;
           }
           :host(.grid-mode-item) .content .description {
               font-size: 0.85em;
               color: #ccc;
               display: -webkit-box;
               -webkit-line-clamp: 2;
               -webkit-box-orient: vertical;
               overflow: hidden;
               text-overflow: ellipsis;
               line-height: 1.3;
               margin-bottom: 0.5rem;
               flex-grow: 1;
           }
          :host(.grid-mode-item) .actions {
              padding: 0rem 0.75rem 0.75rem 0.75rem;
              display: flex;
              justify-content: flex-end;
              margin-top: auto;
              flex-shrink: 0;
          }
           :host(.grid-mode-item) .action-button {
               padding: 0.4rem 0.6rem;
               font-size: 0.9em;
           }

          .video-item {
              color: white;
              cursor: pointer;
              position: relative;
              transition: transform 0.2s ease-in-out, box-shadow 0.2s ease-in-out;
          }
          .video-item:hover {
               transform: translateY(-3px);
               box-shadow: 0 5px 15px rgba(0,0,0,0.4);
          }

          .thumbnail-container {
              position: relative;
              overflow: hidden;
              background-color: #1a1a1a;
          }

          .thumbnail-container img {
              display: block;
              width: 100%;
              height: 100%;
          }

          .content .name {
              margin: 0 0 0.25rem 0;
              font-size: 1.1em;
              font-weight: 600;
              color: #ffffff;
          }
          .content .description {
              margin: 0;
              font-size: 0.9em;
              color: #b3b3b3;
              line-height: 1.4;
          }

          .actions {
              display: flex;
              gap: 0.5rem;
          }

          .action-button {
              padding: 0.5rem 0.75rem;
              border-radius: 4px;
              text-decoration: none;
              color: white;
              font-weight: 500;
              cursor: pointer;
              border: none;
              background-color: #3f3f3f;
              transition: background-color 0.2s ease;
          }
           .action-button:hover {
                background-color: #555;
           }

          .action-button.delete-modal,
          .action-button.delete {
              background-color: #d9534f;
              border-color: #d43f3a;
          }
          .action-button.delete-modal:hover,
          .action-button.delete:hover {
              background-color: #c9302c;
              border-color: #ac2925;
          }

          .video-item .thumbnail-container::before {
              content: '';
              position: absolute;
              top: 0;
              left: 0;
              width: 100%;
              height: 100%;
              background-color: rgba(0, 0, 0, 0);
              transition: background-color 0.3s ease;
              pointer-events: none;
              border-radius: inherit;
          }
          .video-item:hover .thumbnail-container::before {
              background-color: rgba(0, 0, 0, 0.4);
          }

          .play-button {
              position: absolute;
              top: 50%;
              left: 50%;
              transform: translate(-50%, -50%) scale(0.7);
              opacity: 0;
              transition: all 0.3s cubic-bezier(0.175, 0.885, 0.32, 1.275);
              border: none;
              background: rgba(0, 0, 0, 0.6);
              border-radius: 50%;
              width: 55px;
              height: 55px;
              display: flex;
              align-items: center;
              justify-content: center;
              cursor: pointer;
              pointer-events: none;
              box-shadow: 0 2px 5px rgba(0,0,0,0.3);
          }
          .play-button::before {
              content: '▶';
              color: white;
              font-size: 22px;
              margin-left: 4px;
          }
          .video-item:hover .play-button {
              opacity: 1;
              transform: translate(-50%, -50%) scale(1);
          }

          .modal {
              position: fixed;
              top: 0;
              left: 0;
              width: 100%;
              height: 100%;
              background-color: rgba(0, 0, 0, 0.75);
              z-index: 1000;
              display: flex;
              justify-content: center;
              align-items: center;
              opacity: 0;
              transition: opacity 0.2s ease-in-out;
              pointer-events: none;
          }
          .modal[style*="display: flex"] {
              opacity: 1;
              pointer-events: auto;
          }

          .modal-content {
              background-color: #2b2b2b;
              padding: 2rem;
              border-radius: 8px;
              text-align: center;
              color: #eee;
              max-width: 420px;
              width: 90%;
              box-shadow: 0 5px 15px rgba(0,0,0,0.5);
              transform: scale(0.95);
              transition: transform 0.2s ease-in-out;
          }
          .modal[style*="display: flex"] .modal-content {
              transform: scale(1);
          }

          .modal-actions {
              display: flex;
              justify-content: center;
              gap: 1rem;
              margin-top: 1.5rem;
          }
           .modal-content p {
               margin-bottom: 0.5rem;
               color: #ccc;
           }
           .modal-content .name {
                font-weight: bold;
                margin: 0.5rem 0 1rem 0;
                color: white;
                word-break: break-all;
                background-color: #3f3f3f;
                padding: 0.5rem;
                border-radius: 4px;
           }

           .action-button.cancel {
               background-color: #6c757d;
               border-color: #6c757d;
           }
            .action-button.cancel:hover {
               background-color: #5a6268;
               border-color: #545b62;
           }
      `;

    const template = document.createElement("template");

    template.innerHTML = `
      <div class="video-item">
          <div class="thumbnail-container">
              <img class="thumbnail" alt="thumbnail" loading="lazy"/>
              <button class="play-button" aria-label="Play video"></button>
          </div>
          <div class="content">
              <h3 class="name"></h3>
              <p class="description"></p>
          </div>
          <div class="actions">
              <button class="action-button delete-modal" title="Delete Video">Delete</button>
          </div>
      </div>
    `;

    const modalTemplate = document.createElement("template");

    modalTemplate.innerHTML = `
      <div class="modal" id="deleteConfirmModal" style="display: none;">
          <div class="modal-content">
              <p>Are you sure you want to delete:</p>
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

    this.videoItemElement = this.shadow.querySelector('.video-item');
    this.modalElement = this.shadow.querySelector('#deleteConfirmModal');

    this.initialize();
    this.initializeModal();
  }


  initialize() {
       this.videoItemElement.addEventListener('click', (e) => {
          if (e.target.closest('.actions')) {
              return;
          }
           e.preventDefault();
           this.handlePlay();
       });

        const deleteModalButton = this.shadow.querySelector('.delete-modal');
        if (deleteModalButton) {
             deleteModalButton.addEventListener('click', (e) => {
                  e.stopPropagation();
                  const fileName = this.getAttribute("name") || "this video";
                  const fileNameElement = this.modalElement.querySelector(".name");
                  fileNameElement.textContent = fileName;
                  this.modalElement.style.display = "flex";
             });
        }
  }

  initializeModal() {
    const cancelButton = this.modalElement.querySelector(".cancel");
    const deleteConfirmButton = this.modalElement.querySelector(".delete");

    cancelButton.addEventListener("click", (e) => {
        e.stopPropagation();
        this.modalElement.style.display = "none";
    });

    deleteConfirmButton.addEventListener("click", (e) => {
        e.stopPropagation();
        this.modalElement.style.display = "none";
        this.handleDelete();
    });

    this.modalElement.addEventListener("click", (e) => {
      if (e.target === this.modalElement) {
        this.modalElement.style.display = "none";
      }
    });
  }

  handlePlay() {
    const videoId = this.getAttribute("video-id");
    if (videoId) {
      try {
          const watchUrl = new URL(`${window.ENV.API_URL}/watch`);
          watchUrl.searchParams.set('v', videoId);
          window.location.href = watchUrl.toString();
      } catch(e) {
          console.error("Error creating watch URL:", e);
          window.location.href = `/watch?v=${videoId}`;
      }
    } else {
        console.error("Video ID attribute is missing, cannot play.");
    }
  }

  handleDelete() {
    const videoId = this.getAttribute("video-id");
    if (!videoId) {
        console.error("Video ID missing, cannot delete.");
        alert("Cannot delete video: ID is missing.");
        return;
    }

    const deleteButton = this.shadowRoot.querySelector(".delete-modal");
    const originalText = deleteButton.textContent;
    const originalTitle = deleteButton.title;
    deleteButton.textContent = "Deleting...";
    deleteButton.title = "Deleting...";
    deleteButton.disabled = true;


    fetch(`${window.ENV.API_URL}/video/${videoId}`, {
      method: "DELETE",
      headers: {
           'Accept': 'application/json'
       }
    })
      .then(async (response) => {
        if (response.ok) {
          this.remove();
           this.dispatchEvent(new CustomEvent('item-deleted', { bubbles: true, composed: true, detail: { id: videoId } }));
        } else {
           let errorMsg = `HTTP ${response.status} ${response.statusText}`;
           try {
                const errData = await response.json();
                errorMsg = errData.detail || errData.message || errorMsg;
           } catch (e) {

           }
           console.error("Error deleting video:", errorMsg);
           alert(`Failed to delete video: ${errorMsg}`);
           deleteButton.textContent = "Error";
           deleteButton.title = "Deletion failed";
           setTimeout(() => {
               deleteButton.textContent = originalText;
               deleteButton.title = originalTitle;
               deleteButton.disabled = false;
           }, 2500);
        }
      })
      .catch((error) => {
        console.error("Network error during delete:", error);
        alert(`Failed to delete video: Network error. ${error.message}`);
        deleteButton.textContent = "Error";
        deleteButton.title = "Deletion failed";
        setTimeout(() => {
            deleteButton.textContent = originalText;
            deleteButton.title = originalTitle;
            deleteButton.disabled = false;
        }, 2500);
      });
  }

  static get observedAttributes() {
    return ["name", "description", "thumbnail", "video-id"];
  }

  attributeChangedCallback(name, oldValue, newValue) {
    if (name === 'thumbnail') {
      const element = this.shadow.querySelector('.thumbnail');
      if (element) element.src = newValue || '';
    } else if (name === 'name') {
      const element = this.shadow.querySelector('.content .name');
       if (element) element.textContent = newValue || 'Untitled Video';
    } else if (name === 'description') {
      const element = this.shadow.querySelector('.content .description');
       if (element) element.textContent = newValue || '';
    }
  }

  disconnectedCallback() {

  }
}

customElements.define("video-item", VideoItem);