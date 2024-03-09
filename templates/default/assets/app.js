import { getMetadata } from "./api.js";

let baseUrl;
let scrape;
let url;
let title;
let description;
let fetchingMessage;
let bookmarkForm;
let submitButton;
let requiredFields;

window.onload = () => {
    baseUrl = document.querySelector("base").attributes.getNamedItem("href").value;

    scrape = document.getElementById("scrape");

    url = document.getElementById("url");
    title = document.getElementById("title");
    description = document.getElementById("description");
    fetchingMessage = document.getElementById("fetching-metadata-message");
    bookmarkForm = document.getElementById("bookmark-form");
    
    if (bookmarkForm) {
        initializeFormValidation();
    }
    

    if (scrape) {
        scrape.addEventListener("click", handleScrapeButtonClick)
    }
}

const handleScrapeButtonClick = (e) => {
    let inputs = bookmarkForm.querySelectorAll("input, button, textarea");
    inputs.forEach((input) => {
        input.disabled = true;
    });
    
    fetchingMessage.innerHTML = "Fetching metadata";
    fetchingMessage.classList.remove("hidden");

    getMetadata(baseUrl, url.value)
        .then(resp => {
            if (resp.title || resp.description) {
                title.value = resp.title;
                description.value = resp.description;
                validateForm();
            }
        })
        .catch(err => {
            console.log("Error", err)
        })
        .finally(_ => {
            inputs.forEach((input) => {
                input.disabled = false;
            });

            fetchingMessage.classList.add("hidden");
            fetchingMessage.innerHTML = "Metadata fetched";
        })
}

const initializeFormValidation = () => {
    submitButton = bookmarkForm.querySelector("button[type=submit]");
    requiredFields = bookmarkForm.querySelectorAll("[required]");

    requiredFields.forEach((input) => {
        if (input.value == "") {
            submitButton.disabled = true;
        }

        input.addEventListener("keyup", (e) => {
            validateForm();
        });

        input.addEventListener("change", (e) => {
            validateForm();
        });
    });

    validateForm()
}

export const validateForm = () => {
    var isValid = true;
    requiredFields.forEach((input) => {
        if (input.value == "") {
            isValid = false;
            input.setAttribute("aria-invalid", "true");
            input.classList.add("invalid");
        } else {
            input.setAttribute("aria-invalid", "false");
            input.classList.remove("invalid");
        }
    })
    submitButton.disabled = !isValid
}