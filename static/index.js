const loadConfig = async () => {
    const response = await fetch(`${window.ENV.API_URL}/config`, {
        method: "GET"
    });

    if (!response.ok) {
        throw new Error('Network response was not ok ' + response.statusText);
    }
    
    const data = await response.json()
    console.log(data)
    localStorage.setItem("FILE_SIZE_LIMIT", data["FILE_SIZE_LIMIT"])
}

loadConfig()
