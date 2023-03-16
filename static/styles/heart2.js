const heartIcons = document.querySelectorAll(".heart-icon");

heartIcons.forEach((heartIcon) => {
  heartIcon.addEventListener("click", () => {
    heartIcon.classList.toggle("liked");
    fetch('/profil', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ id: heartIcon.getAttribute("id") })
      });
  });
});