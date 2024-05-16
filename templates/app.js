document.addEventListener("DOMContentLoaded", function() {
    getLogo();
});

const main = document.getElementById("main")

function getLogo(){
    document.getElementById('logo').innerHTML = generateLogoTextHTML() + generateLogoImgHTML()
}
function generateLogoTextHTML() {
    return `<h1 class="inline" >Board</h1>`;
}
function generateLogoImgHTML() {
    return `<img class="inline" src="img.png" href="./img.png">`;
}
function getMainBody(){
    console.log("getMainBody");
    const body = document.createElement("div");
    body.textContent = "Treść głównego obszaru"; 
    const main = document.getElementById('main');
    main.appendChild(body);
}
