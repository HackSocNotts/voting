function register() {
    var id = document.getElementById("id").value

    if (/^\d+$/.test(id)) {
        var req = new XMLHttpRequest()

        req.onreadystatechange = function() {
            if (this.readyState != 4) return

            console.log(this.status, req.responseText)
        }

        req.open("POST", "/register/")
        req.send(id)
    } else {
        console.error("invalid id", id)
    }
}