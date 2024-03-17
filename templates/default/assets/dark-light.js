const initialize = () => {
    let modeSwitchLink = document.getElementById("mode-switch");
    let modeSwitchIcon = document.getElementById("mode-switch-icon")
    let baseUrl = document.querySelector("base").attributes.getNamedItem("href").value;

    let theme = localStorage.getItem("theme")
    if (theme === null) {
        let prefersDarkMode = window.matchMedia("(prefers-color-scheme: dark)").matches;
        theme = prefersDarkMode ? "dark" : "light";
        localStorage.setItem("theme", theme)
    }

    if (theme == "dark") {
        document.documentElement.classList.remove("light")
        modeSwitchLink.setAttribute("aria-label", "Switch to light mode");
        modeSwitchIcon.setAttribute("src", `${baseUrl}assets/light.png`);
        modeSwitchIcon.setAttribute("alt", "Sun icon");
    } else {
        document.documentElement.classList.add("light");
        modeSwitchLink.setAttribute("aria-label", "Switch to dark mode");
        modeSwitchIcon.setAttribute("src", `${baseUrl}assets/dark.png`);
        modeSwitchIcon.setAttribute("alt", "Moon icon");
    }

    modeSwitchLink.addEventListener("click", (e) => {
        e.preventDefault();
        e.stopPropagation();
        let theme = localStorage.getItem("theme");
        if (theme == "dark") {
            localStorage.setItem("theme", "light");
            document.documentElement.classList.add("light");
            modeSwitchLink.setAttribute("aria-label", "Switch to dark mode");
            modeSwitchIcon.setAttribute("src", `${baseUrl}assets/dark.png`);
            modeSwitchIcon.setAttribute("alt", "Moon icon");
        } else {
            localStorage.setItem("theme", "dark");
            document.documentElement.classList.remove("light")
            modeSwitchLink.setAttribute("aria-label", "Switch to light mode");
            modeSwitchIcon.setAttribute("src", `${baseUrl}assets/light.png`);
            modeSwitchIcon.setAttribute("alt", "Sun icon");
        }
    });
}

window.addEventListener("load", initialize);
