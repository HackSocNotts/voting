function refresh() {
    var req = new XMLHttpRequest()

    req.onreadystatechange = function() {
        if (this.readyState != 4) return

        if (this.status == 200) {
            var results = JSON.parse(req.responseText)
            render(results)
        } else {
            // handle error
        }
    }

    req.open("GET", "./results")
    req.send()
}

function render(results) {
    document.getElementById("votes").innerHTML = results["num_votes"]
    document.getElementById("ballots").innerHTML = results["num_ballots"]

    var winnersUl = document.getElementById("winners")
    winnersUl.innerHTML = ""

    for (var position of Object.keys(results["winners"])) {
        var winner = results["winners"][position]
        var li = document.createElement("li")
        var h2 = document.createElement("h2")
        h2.innerHTML = position
        var span = document.createElement("span")
        span.innerHTML = winner
        li.appendChild(h2)
        li.appendChild(span)
        winnersUl.appendChild(li)
    }
}
