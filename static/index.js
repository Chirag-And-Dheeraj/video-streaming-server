const loadConfig = async () => {
  try {
    const response = await fetch(`${window.ENV.API_URL}/config`, {
      method: "GET",
    });
    if (!response.ok) {
      throw new Error(`Network response was not ok (${response.status})`);
    }
    const data = await response.json();
    localStorage.setItem("FILE_SIZE_LIMIT", data.file_size_limit);
    localStorage.setItem(
      "SUPPORTED_FILE_TYPES",
      JSON.stringify(data.supported_file_types)
    );
    console.log("Config loaded successfully:", data);
  } catch (error) {
    console.error("Error loading config:", error.message);
  }
};

loadConfig();
