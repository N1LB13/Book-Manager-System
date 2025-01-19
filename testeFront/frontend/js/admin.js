document.getElementById("logout").addEventListener("click", () => {
    localStorage.removeItem("token");
    window.location.href = "index.html";
});

document.getElementById("add-book-form").addEventListener("submit", async (event) => {
    event.preventDefault();
    const token = localStorage.getItem("token");
    const title = document.getElementById("title").value;
    const author = document.getElementById("author").value;
    const genre = document.getElementById("genre").value;
    const rating = document.getElementById("rating").value;

    try {
        const response = await fetch("http://localhost:8080/api/livros", {
            method: "POST",
            headers: { "Authorization": `Bearer ${token}`, "Content-Type": "application/json" },
            body: JSON.stringify({ title, author, main_genre: genre, rating: parseFloat(rating) }),
        });

        if (response.ok) {
            alert("Livro adicionado com sucesso!");
        } else {
            alert("Erro ao adicionar livro.");
        }
    } catch (error) {
        alert("Erro ao conectar ao servidor.");
    }
});

document.getElementById("delete-book-form").addEventListener("submit", async (event) => {
    event.preventDefault();
    const token = localStorage.getItem("token");
    const bookId = document.getElementById("book-id").value;

    try {
        const response = await fetch(`http://localhost:8080/api/livros/${bookId}`, {
            method: "DELETE",
            headers: { "Authorization": `Bearer ${token}` },
        });

        if (response.ok) {
            alert("Livro removido com sucesso!");
        } else {
            alert("Erro ao remover livro.");
        }
    } catch (error) {
        alert("Erro ao conectar ao servidor.");
    }
});
