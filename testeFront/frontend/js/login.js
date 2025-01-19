const loginForm = document.getElementById("login-form");
const message = document.getElementById("message");

loginForm.addEventListener("submit", async (event) => {
    event.preventDefault();
    const username = document.getElementById("username").value;
    const password = document.getElementById("password").value;

    try {
        const response = await fetch("http://localhost:8080/login", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ usuario_id: username, senha: password }),
        });

        const data = await response.json();
        if (response.ok) {
            localStorage.setItem("token", data.token);
            const claims = JSON.parse(atob(data.token.split(".")[1]));
            if (claims.role === "admin") {
                window.location.href = "admin.html";
            } else {
                window.location.href = "user.html";
            }
        } else {
            message.innerText = data.error || "Erro no login";
        }
    } catch (error) {
        message.innerText = "Erro ao conectar ao servidor.";
    }
});
