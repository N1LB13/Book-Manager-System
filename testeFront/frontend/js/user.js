document.getElementById("logout").addEventListener("click", () => {
    localStorage.removeItem("token");
    window.location.href = "index.html";
});

document.getElementById("search-book-form").addEventListener("submit", async (event) => {
    event.preventDefault();
    const token = localStorage.getItem("token");
    const bookId = document.getElementById("book-id").value;

    try {
        const response = await fetch(`http://localhost:8080/api/livros/${bookId}`, {
            method: "GET",
            headers: { "Authorization": `Bearer ${token}` },
        });

        const data = await response.json();
        if (response.ok) {
            document.getElementById("book-details").innerText = JSON.stringify(data, null, 2);
        } else {
            alert("Livro não encontrado.");
        }
    } catch (error) {
        alert("Erro ao conectar ao servidor.");
    }
});
