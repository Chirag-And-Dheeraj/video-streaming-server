class Toast {
  constructor() {
    this.initialized = false;
    this.init();
  }

  init() {
    if (this.initialized) {
      return;
    }
    this.injectCSS();
    this.injectHTML();
    this.initialized = true;
  }

  injectCSS() {
    // Check if CSS is already injected
    if (document.getElementById("toast-styles")) {
      return;
    }

    const style = document.createElement("style");
    style.id = "toast-styles";
    style.textContent = `
            .toast-container {
                position: fixed;
                bottom: 2.5rem;
                right: 2.5rem;
                z-index: 9999;
                display: flex;
                flex-direction: column;
                gap: 10px;
                pointer-events: none;
                // border: solid red 2px;
            }

            .toast {
                border-radius: 8px;
                box-shadow: 0 4px 12px rgba(0, 0, 0, 0.75);
                background: #262626;
                padding: 12px 16px;
                min-width: 300px;
                max-width: 350px;
                display: flex;
                align-items: center;
                gap: 12px;
                pointer-events: auto;
                transform: translateY(100%);
                transition: transform 0.3s ease, opacity 0.3s ease;
            }

            .toast.show {
                transform: translateY(0);
            }

            .toast-content {
                flex: 1;
            }

            .toast-title {
                font-weight: 900;
                font-size: 1rem;
                color: #f5f5f5;
                margin: 0 0 2px 0;
            }

            .toast-message {
                font-size: 0.90rem;
                color:  #f5f5f5;
                margin: 0;
            }

            .toast-icon {
                width: 20px;
                height: 20px;
                border-radius: 50%;
                display: flex;
                align-items: center;
                justify-content: center;
                font-size: 12px;
                color: white;
                font-weight: bold;
            }

            .toast.success .toast-icon {
                // background: #10b981;
                font-size: 1.5rem;
                margin: 0.5rem;
                padding: 0.5rem;
            }

            .toast.error .toast-icon {
                // background: #ef4444;
                font-size: 1.5rem;
                margin: 0.5rem;
                padding: 0.5rem;
            }

            .toast-close {
                background: none;
                border: none;
                color: #9ca3af;
                cursor: pointer;
                padding: 4px;
                border-radius: 4px;
                transition: color 0.2s ease;
                font-size: 16px;
                line-height: 1;
            }

            .toast-close:hover {
                color: #6b7280;
                background: #f3f4f6;
            }
        `;

    document.head.appendChild(style);
  }

  injectHTML() {
    if (document.getElementById("toast-container")) {
      return;
    }
    const container = document.createElement("div");
    container.id = "toast-container";
    container.className = "toast-container";
    document.body.appendChild(container);
  }

  show(title, message, type = "success", duration = 7000) {
    this.init();
    const container = document.getElementById("toast-container");
    const toast = document.createElement("div");
    toast.className = `toast ${type}`;
    const icon = document.createElement("div");
    icon.className = "toast-icon";
    icon.textContent = type === "success" ? "✅" : "❌";
    const content = document.createElement("div");
    content.className = "toast-content";
    const h4 = document.createElement("h4");
    h4.className = "toast-title";
    h4.textContent = title.length > 15 ? title.substring(0, 15) + "…" : title;
    const p = document.createElement("p");
    p.className = "toast-message";
    p.textContent = message;
    content.append(h4, p);
    const btn = document.createElement("button");
    btn.className = "toast-close";
    btn.type = "button";
    btn.textContent = "✖";
    btn.addEventListener("click", () => {
      toast.classList.remove("show");
      setTimeout(() => toast.remove(), 300);
    });
    toast.append(icon, content, btn);
    container.appendChild(toast);
    requestAnimationFrame(() => toast.classList.add("show"));
    setTimeout(() => {
      toast.classList.remove("show");
      setTimeout(() => toast.remove(), 300);
    }, duration);
  }

  close(toast) {
    toast.style.transform = "translateY(100%)";
    toast.style.opacity = "0";

    setTimeout(() => {
      if (toast.parentNode) {
        toast.parentNode.removeChild(toast);
      }
    }, 300);
  }

  success(title, message, duration) {
    return this.show(title, message, "success", duration);
  }
  error(title, message, duration) {
    return this.show(title, message, "error", duration);
  }
}
