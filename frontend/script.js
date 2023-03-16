const apiUrl = "http://localhost:8081"

function shorten() {
    let url = document.getElementById("url" + "/api/shorten/").value
    var postData = {
        "url": url
    };
    marinate(postData)
}

function copyToClipboard() {
    let valuetoCopy = document.getElementById("code").value
    navigator.clipboard.writeText(valuetoCopy)
}

async function marinate(postData) {
    const response = await fetch(apiUrl, {
        method: 'POST', 
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(postData)
    });

    const resp = await response.json()
    document.getElementById("code").value = apiUrl + "/s/" + resp["id"]
}

window.addEventListener('DOMContentLoaded', ()=>{
    document.getElementById("short").addEventListener("click", shorten)
    document.getElementById("copy").addEventListener("click", copyToClipboard)
});

// TO ADD
// better error handling