function createStars() {
    let stars = document.getElementById("stars");
    if (!stars) {
        stars = document.createElement("div");
        stars.id = "stars";
        document.body.appendChild(stars);
    }
    
    stars.innerHTML = "";
    const height = Math.max(
        document.body.scrollHeight,
        document.body.offsetHeight,
        document.documentElement.scrollHeight,
        document.documentElement.offsetHeight,
        document.documentElement.clientHeight
    );
    const width = document.documentElement.scrollWidth;
    
    stars.style.height = height + "px";
    stars.style.width = width + "px";

    for (let i = 0; i < 250; i++) {
        const star = document.createElement("div");
        star.className = "star";
        star.style.left = Math.random() * width + "px";
        star.style.top = Math.random() * height + "px";
        star.style.animationDuration = 1.5 + Math.random() * 3 + "s";
        stars.appendChild(star);
    }
}

window.addEventListener("load", function() {
    setTimeout(createStars, 100); // wait for full render
});

window.addEventListener("resize", createStars);

window.onload = function() {
    const container = document.getElementById("COMMENT-AREA");
    if (container) container.scrollTop = container.scrollHeight;
};