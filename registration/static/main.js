function register() {
    var id = document.getElementById("id").value

    if (/^\d+$/.test(id)) {
        var req = new XMLHttpRequest()

        req.onreadystatechange = function() {
            if (this.readyState != 4) return

            if (this.status == 200) {
                document.getElementById("error").innerHTML = ""

                document.getElementById("register").setAttribute("disabled", true)
                document.getElementById("register").classList.add("hidden")
                document.getElementById("register").innerHTML = "Registered"

                document.getElementById("ballot").removeAttribute("disabled")
                document.getElementById("ballot").classList.remove("hidden")

                document.getElementById("message").innerHTML = "Your ballot has been created. Click the button below to cast your votes."
            } else {
                document.getElementById("error").innerHTML = req.responseText
            }
        }

        req.open("POST", "/register/")
        req.send(id)

        document.getElementById("register").innerHTML = "Registering"
    } else {
        document.getElementById("error").innerHTML = "Invalid student ID"
    }
}

function openBallot() {
    // this will go to the URL for the ballot
}