document.getElementById("btn").addEventListener("click", async () => {
  const username = document.getElementById("username").value.trim();
  const password = document.getElementById("password").value.trim();
  const msg = document.getElementById("msg");
  const btn = document.getElementById("btn");

  msg.textContent = "";
  msg.style.color = "#ff5252";

  // Validación básica
  if (!username || !password) {
    msg.textContent = "Completa usuario y contraseña";
    return;
  }

  // Deshabilitar botón durante la solicitud
  btn.disabled = true;
  btn.textContent = "Validando...";

  try {
    const res = await fetch("/login", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ username, password })
    });

    if (!res.ok) {
      msg.textContent = "Credenciales incorrectas";
      btn.disabled = false;
      btn.textContent = "Entrar";
      return;
    }

    // Login exitoso
    localStorage.setItem("role", username === "admin" ? "admin" : "user");
    localStorage.setItem("username", username);

    msg.style.color = "#4caf50";
    msg.textContent = "Login exitoso";

    // Redirigir al dashboard después de 1 segundo
    setTimeout(() => {
      window.location.href = "/dashboard.html";
    }, 1000);

  } catch (e) {
    msg.textContent = "Error de conexión: " + e.message;
    btn.disabled = false;
    btn.textContent = "Entrar";
  }
});
