// Authentication
function checkAuth() {
  const token = localStorage.getItem("token");
  if (!token) {
    window.location.href = "/login";
  }
  return token;
}

function logout() {
  localStorage.removeItem("token");
  localStorage.removeItem("refreshToken");
  window.location.href = "/login";
}

// API calls
async function apiCall(url, options = {}) {
  const token = localStorage.getItem("token");
  const defaultOptions = {
    headers: {
      Authorization: `Bearer ${token}`,
      "Content-Type": "application/json",
    },
  };

  const response = await fetch(url, { ...defaultOptions, ...options });

  if (!response.ok) {
    const data = await response.json();
    throw new Error(data.error || "API request failed");
  }

  return response.json();
}

// Form handling
function getFormData(form) {
  const formData = new FormData(form);
  const jsonObject = {};
  formData.forEach((value, key) => {
    jsonObject[key] = value;
  });
  return jsonObject;
}

// UI helpers
function showMessage(message, type = "error") {
  const messageElement = document.getElementById("error-message");
  messageElement.textContent = message;
  messageElement.classList.remove("hidden", "message-error", "message-success");
  messageElement.classList.add(`message-${type}`);
}

function hideMessage() {
  const messageElement = document.getElementById("error-message");
  messageElement.classList.add("hidden");
}

function setButtonLoading(button, isLoading) {
  button.disabled = isLoading;
  button.textContent = isLoading ? "Loading..." : button.dataset.originalText;
}

// Chirp handling
async function deleteChirp(id) {
  try {
    await apiCall(`/api/chirps/${id}`, { method: "DELETE" });
    if (typeof loadChirps === "function") {
      await loadChirps();
    }
  } catch (error) {
    showMessage(error.message);
  }
}

// Event listeners
document.addEventListener("DOMContentLoaded", function () {
  // Store original button text
  document.querySelectorAll("button").forEach((button) => {
    button.dataset.originalText = button.textContent;
  });

  // Add delete chirp event listeners
  document.addEventListener("click", function (e) {
    if (e.target.classList.contains("delete-chirp")) {
      const chirpId = e.target.dataset.chirpId;
      deleteChirp(chirpId);
    }
  });
});
