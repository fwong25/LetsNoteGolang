const title_p = document.getElementById("title_p");
const content_p = document.getElementById("content_p");
const title_input = document.getElementById("title_input");
const content_input = document.getElementById("content_input");
const confirmBtn = document.getElementById("confirmUpdateBtn");

confirmBtn.onclick = assignTitleAndContent;

function assignTitleAndContent() {
    title_input.value =  title_p.innerHTML;
    content_input.value =  content_p.innerHTML;
    console.log(title_p.innerHTML);
    console.log(title_input.value);
    return true;
}