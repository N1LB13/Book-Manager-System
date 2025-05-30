// Função de logout
document.getElementById("logout").addEventListener("click", () => {
    localStorage.removeItem("token");
    window.location.href = "index.html";
});

// Função de busca por livro
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

// Função para obter recomendações
document.getElementById("genre-form").addEventListener("submit", async (event) => {
    event.preventDefault();
    console.log("Formulário enviado!");

    const token = localStorage.getItem("token");
    console.log("Token JWT:", token);

    const selectedGenres = Array.from(document.getElementById("genres").selectedOptions).map(option => option.value);
    console.log("Gêneros selecionados:", selectedGenres);

    if (selectedGenres.length !== 3) {
        alert("Você deve selecionar exatamente 3 gêneros.");
        return;
    }

    try {
        console.log("Enviando requisição...");
        const response = await fetch("http://localhost:8080/api/recommend", {
            method: "POST",
            headers: { 
                "Authorization": `Bearer ${token}`,
                "Content-Type": "application/json"
            },
            body: JSON.stringify({ generos: selectedGenres }),
        });
        console.log("Resposta recebida:", response);

        if (!response.ok) {
            throw new Error("Erro na requisição: " + response.status);
        }

        const data = await response.json();
        console.log("Dados recebidos:", data);

        const recommendationsDiv = document.getElementById("recommendations");
        recommendationsDiv.innerHTML = "<h3>Recomendações:</h3>";
        recommendationsDiv.innerHTML += data.map(book => `
            <div>
                <strong>${book.title}</strong> - ${book.author} (Rating: ${book.rating})
            </div>
        `).join("");
    } catch (error) {
        console.error("Erro ao enviar a requisição:", error);
        alert(`Erro: ${error.message}`);
    }
});
