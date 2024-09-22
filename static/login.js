const loginForm = document.getElementById("login-form");
const divOutput = document.getElementById("divOutput");
const emailElement = document.getElementById("email");
const passwordElement = document.getElementById("password");

loginForm.addEventListener("submit", handleSubmit);

async function handleSubmit(e) {
  e.preventDefault();

  const email = emailElement.value.trim();
  const password = passwordElement.value.trim();

  if (!email || !password) {
    divOutput.textContent = "Please fill in both email and password";
    return;
  }

  try {
    const response = await fetch(`${window.ENV.API_URL}/login`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ email: email, password: password }),
    });

    if (response.ok) {
      window.location.href = "/";
    } else if (!response.ok) {
      switch (response.status) {
        case 400:
          throw new Error("Invalid email or password");
        case 401:
          throw new Error("Invalid email or password");
        case 404:
          throw new Error("User does not exist");
        case 500:
          throw new Error("Internal server error");
        default:
          throw new Error(`Server responded with status ${response.status}`);
      }
    }

    const data = await response.json();
    divOutput.textContent = data.message;
    emailElement.value = "";
    passwordElement.value = "";
  } catch (error) {
    console.error("Error:", error);
    let errorMessage = error.message;

    if (errorMessage.includes("Failed to fetch")) {
      errorMessage = "Network error. Please check your internet connection.";
    }

    divOutput.textContent = errorMessage;
  }
}
