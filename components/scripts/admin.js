import { getMetadata, getTags } from "./api.js";

let baseUrl;
let scrape;
let url;
let title;
let description;
let fetchingMessage;
let bookmarkForm;
let submitButton;
let requiredFields;
let tagsInput;
let allTags;
let tagSuggestions;

const initialize = () => {
    baseUrl = document.querySelector("base").attributes.getNamedItem("href").value;
    scrape = document.getElementById("scrape");
    url = document.getElementById("url");
    title = document.getElementById("title");
    description = document.getElementById("description");
    fetchingMessage = document.getElementById("fetching-metadata-message");
    bookmarkForm = document.getElementById("bookmark-form");
    tagsInput = document.querySelector("input#tags");
    tagSuggestions = document.getElementById("tag-suggestions");

    if (bookmarkForm) {
        initializeFormValidation();
        let inputs = bookmarkForm.querySelectorAll("input, textarea, button");
        inputs.forEach((input) => {
            input.addEventListener("focus", handleFormInputFocus);
        })
    }

    if (scrape) {
        scrape.addEventListener("click", handleScrapeButtonClick);
    }

    if (tagsInput) {
        tagsInput.addEventListener("keyup", handleTagsKeyUp);
        window.addEventListener("resize", handleWindowResize);

        getTags(baseUrl).then(resp => {
            allTags = resp.tags;
        })
            .catch(err => {
                console.log("Error fetching tags", err);
            })

        handleWindowResize();
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
            console.log("Error fetching metadata", err)
        })
        .finally(_ => {
            inputs.forEach((input) => {
                input.disabled = false;
            });

            fetchingMessage.classList.add("hidden");
            fetchingMessage.innerHTML = "Metadata fetched";
        })
}

const handleTagsKeyUp = (e) => {
    let lastSpaceIndex = tagsInput.value.lastIndexOf(" ") + 1;
    let currentTag = tagsInput.value.substr(lastSpaceIndex, tagsInput.value.length);
    let matchingTags = allTags.filter((tag) => {
        if (!currentTag) {
            return false;
        }
        return tag.startsWith(currentTag) && !tagsInput.value.includes(tag);
    })

    let tagOptions = [];
    matchingTags.forEach((tag) => {
        let option = document.createElement("li");
        option.innerHTML = tag;
        option.setAttribute("tabindex", "0");
        option.addEventListener("keyup", handleTagSuggestionKeyUp);
        option.addEventListener("click", handleTagSuggestionClick);
        tagOptions.push(option);
    })

    tagSuggestions.replaceChildren(...tagOptions);

    if (e.key == "ArrowDown" && tagOptions.length > 0) {
        tagSuggestions.children[0].focus();
        return;
    }
    if (e.key == "Escape") {
        tagSuggestions.replaceChildren();
    }

}

const handleTagSuggestionKeyUp = (e) => {
    let currentSuggestion = e.target;
    if (e.key == "ArrowDown" && currentSuggestion.nextSibling) {
        currentSuggestion.nextSibling.focus();
    }
    if (e.key == "ArrowUp") {
        if (currentSuggestion.previousSibling) {
            currentSuggestion.previousSibling.focus();
        } else {
            tagsInput.focus();
        }
    }
    if (e.key == "Enter") {
        let lastSpaceIndex = tagsInput.value.lastIndexOf(" ") + 1;
        let inputValueWithoutTagPrefix = tagsInput.value.substr(0, lastSpaceIndex);
        tagsInput.value = inputValueWithoutTagPrefix + e.target.innerHTML + " ";
        tagsInput.focus();
        tagSuggestions.replaceChildren();
    }
    if (e.key == "Escape") {
        tagsInput.focus();
        tagSuggestions.replaceChildren();
    }
}

const handleTagSuggestionClick = (e) => {
    let lastSpaceIndex = tagsInput.value.lastIndexOf(" ") + 1;
    let inputValueWithoutTagPrefix = tagsInput.value.substr(0, lastSpaceIndex);
    tagsInput.value = inputValueWithoutTagPrefix + e.target.innerHTML + " ";
    tagsInput.focus();
    tagSuggestions.replaceChildren();
}

const handleFormInputFocus = (e) => {
    if (e.target.id == "tags") {
        return;
    }

    tagSuggestions.replaceChildren();
}

const handleWindowResize = (e) => {
    let coords = tagsInput.getBoundingClientRect();
    tagSuggestions.style.top = `${coords.bottom + 2 + window.scrollY}px`;
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

const validateForm = () => {
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

window.addEventListener("load", initialize)