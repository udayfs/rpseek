async function search(query) {
    const result = document.getElementById("result");
    result.innerHTML = "";
    const response = await fetch("/search", {
        method: "POST",
        headers: { "Content-Type": "text/plain" },
        body: JSON.stringify({
            query: query,
        }),
    });

    const docs = await response.json();
    for (const doc of docs) {
        let item = document.createElement("span");
        item.innerHTML = `${doc.doc_id} => ${doc.rank}`;
        item.appendChild(document.createElement("br"));
        result.appendChild(item);
    }
}

let query = document.getElementById("search");
let searchBtn = document.getElementById("btn");

searchBtn.addEventListener("click", async (e) => {
    e.preventDefault();
    await search(query.value);
});
