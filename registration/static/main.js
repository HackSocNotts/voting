function register() {
    var id = document.getElementById("id").value

    if (/^\d+$/.test(id)) {
        var req = new XMLHttpRequest()

        req.onreadystatechange = function() {
            if (this.readyState != 4) return

            if (this.status == 200) {
                document.getElementById("error").innerHTML = ""
                console.log(req.responseText)
            } else {
                document.getElementById("error").innerHTML = req.responseText
            }
        }

        req.open("POST", "/register/")
        req.send(id)
    } else {
        document.getElementById("error").innerHTML = "Invalid student ID"
    }
}