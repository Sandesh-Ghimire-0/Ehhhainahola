window.addEventListener('DOMContentLoaded', function() {
    const stars = document.createElement("div");
    const STAR_COUNT = 250;
    document.body.appendChild(stars);
    stars.id = "stars";

    function createStars() {
        stars.innerHTML = "";
        const height = document.documentElement.scrollHeight;    
        const width = document.documentElement.scrollWidth;
        stars.style.height = height + "px";
        stars.style.width = width + "px";

        for (let i = 0; i < STAR_COUNT; i++) {
            const star = document.createElement("div");
            star.className = "star";
            star.style.left = Math.random() * width + "px";
            star.style.top = Math.random() * height + "px";
            star.style.animationDuration = 1.5 + Math.random() * 3 + "s";
            stars.appendChild(star);
        }
    }

    window.addEventListener("resize", createStars);
    window.addEventListener("load", createStars);
});

window.onload = function () {
    const container = document.getElementById("COMMENT-AREA");
    if (container) container.scrollTop = container.scrollHeight;
};