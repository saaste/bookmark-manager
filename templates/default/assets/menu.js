let menuItems;

const initialize = () => {
    let menuButton = document.getElementById("menu-icon");
    menuItems = document.getElementById("menu-items");
    menuButton.addEventListener("click", handleMenuClick);
}

const handleMenuClick = (e) => {
    e.preventDefault();
    e.stopPropagation();

    if (menuItems.classList.contains("desktop-only")) {
        menuItems.classList.remove("desktop-only");
    } else {
        menuItems.classList.add("desktop-only");
    }
}

window.addEventListener("load", initialize);