const baseUrl = "http://localhost:3000"; // Adjust if hosted elsewhere

async function login() {
  const role = document.getElementById("role").value;
  const emailOrUsername = document.getElementById("emailOrUsername").value;
  const password = document.getElementById("password").value;

  const endpoint = role === "admin" ? "/login_a" : "/login_s";
  const payload = role === "admin"
    ? { username: emailOrUsername, password }
    : { email: emailOrUsername, password };

  const res = await fetch(baseUrl + endpoint, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  });

  const data = await res.json();

  if (res.ok) {
    localStorage.setItem("user", JSON.stringify(data));
  
    if (role === "admin") {
      window.location.href = "admin.html";
    } else {
      // Redirect with student ID
      console.log(data.ID)
      window.location.href = `student.html?id=${data.ID}`;
    }
  } else {
    alert(data.error || "Login failed");
  }
  
}

async function signup() {
  const student = {
    name: document.getElementById("signupName").value,
    email: document.getElementById("signupEmail").value,
    RollNumber: document.getElementById("signupRoll").value,
    phone: document.getElementById("signupPhone").value,
    password: document.getElementById("signupPassword").value,
  };

  const res = await fetch(baseUrl + "/register_s", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(student),
  });

  const data = await res.json();

  if (res.ok) {
    alert("Registration successful. Assigned Room: " + data.assignedRoom);
  } else {
    alert(data.error);
  }
}
