let menuItems;

const initialize = () => {
    let menuButton = document.getElementById("menu-icon");
    menuItems = document.getElementById("menu-items");
    menuButton.addEventListener("click", handleMenuClick);
}

const handleMenuClick = (e) => {
    e.preventDefault();
    e.stopPropagation();

    if (menuItems.classList.contains("hidden-mobile")) {
        menuItems.classList.remove("hidden-mobile");
    } else {
        menuItems.classList.add("hidden-mobile");
    }
}

window.addEventListener("load", initialize);