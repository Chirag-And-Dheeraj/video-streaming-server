const registerForm = document.getElementById("register-form");
const divOutput = document.getElementById("divOutput");
const usernameElement = document.getElementById("username");
const emailElement = document.getElementById("email");
const passwordElement = document.getElementById("password");
const confirmPasswordElement = document.getElementById("confirm_password");

registerForm.addEventListener("submit", handleSubmit);

async function handleSubmit(e) {
  e.preventDefault();

  const username = usernameElement.value.trim();
  const email = emailElement.value.trim();
  const password = passwordElement.value.trim();
  const confirmPassword = confirmPasswordElement.value.trim();

  // Client-side validation
  if (!username || !email || !password || !confirmPassword) {
    divOutput.textContent = "Please fill in all required fields";
    return;
  }

  if (password !== confirmPassword) {
    divOutput.textContent = "Passwords do not match";
    return;
  }

  try {
    const response = await fetch(`${window.ENV.API_URL}/register`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        username,
        email,
        password,
        confirm_password: confirmPassword,
      }),
    });

    let responseData;
    if (response.ok) {
      responseData = await response.json();
      divOutput.textContent = `Registration successful! Welcome, ${responseData.username}`;
      usernameElement.value = "";
      emailElement.value = "";
      passwordElement.value = "";
      confirmPasswordElement.value = "";
      window.location.href = "/login";
    } else {
      // Try to parse JSON response for errors
      responseData = await response.json().catch(() => null);
      const errorMessage =
        responseData?.error ||
        `Server responded with status ${response.status}`;
      divOutput.textContent = errorMessage;
    }
  } catch (error) {
    console.error("Error:", error);
    let errorMessage = error.message;

    if (errorMessage.includes("Failed to fetch")) {
      errorMessage = "Network error. Please check your internet connection.";
    }

    divOutput.textContent = errorMessage;
  }
}
