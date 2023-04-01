const apiUrl = "http://localhost:8081"

async function shorten() {
    let postData = {
        "url": document.querySelector("#url").value
    };

    const response = await fetch(apiUrl + "/api/shorten", {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(postData)
    });

    let resp = await response.json()

    if (response.status === 200) {
        document.querySelector("#shortened-url").value = apiUrl + "/s/" + resp["id"]

        document.querySelector("button").innerText = "Copy to Clipboard"
        document.querySelector("button").onclick = copyToClipboard()

    } else {
        document.querySelector(".error").innerText = "Server has encountered some error."
        document.querySelector(".error").style.opacity = 100
    }
}

function copyToClipboard() {
    let valuetoCopy = document.querySelector("#shortened-url").value
    navigator.clipboard.writeText(valuetoCopy)
}