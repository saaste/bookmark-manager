import { getMetadata } from "./api.js";

let baseUrl;

window.onload = () => {
    baseUrl = document.querySelector("base").attributes.getNamedItem("href").value;
    
    let scrape = document.getElementById("scrape");
    let url = document.getElementById("url");
    let title = document.getElementById("title");
    let description = document.getElementById("description");
    let loader = document.getElementById("scrape-loader");

    if (!scrape) {
        return;
    }

    scrape.addEventListener("click", (e) => {
        scrape.disabled = true;
        url.disabled = true;
        title.disabled = true;
        description.disabled = true;
        loader.classList.remove("hidden");

        getMetadata(baseUrl, url.value)
        .then(resp => {
            if (resp.title) {
                title.value = resp.title;
            }

            if (resp.description) {
                description.value = resp.description;
            }
        })
        .catch(err => {
            console.log("Error", err)
        })
        .finally(_ => {
            scrape.disabled = false;
            url.disabled = false;
            title.disabled = false;
            description.disabled = false;
            loader.classList.add("hidden");
        })
    });
}