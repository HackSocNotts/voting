var candidates

function load() {
    getCandidates()
}

function getCandidates() {
    var req = new XMLHttpRequest()

    req.onreadystatechange = function() {
        if (this.readyState != 4) return

        if (this.status == 200) {
            candidates = JSON.parse(req.responseText)
            renderForm()
        } else {
            // TODO: handle error
        }
    }

    req.open("GET", "/candidates/")
    req.send()
}

function renderForm() {
    for (var position of candidates) {
        var section = document.createElement("section")

        var h2 = document.createElement("h2")
        h2.innerHTML = position.role
        section.appendChild(h2)

        var ul = document.createElement("ul")
        ul.classList.add("choices")
        section.appendChild(ul)

        for (var candidate of position.candidates) {
            var li = document.createElement("li")
            li.classList.add("choice")
            if (candidate == "Re-open Nominations") {
                li.classList.add("ron")
            }
            ul.appendChild(li)

            var rank = document.createElement("span")
            rank.classList.add("rank")
            rank.innerHTML = "-"
            li.appendChild(rank)

            var name = document.createElement("span")
            name.classList.add("name")
            name.innerHTML = candidate
            li.appendChild(name)
        }

        document.getElementById("form").appendChild(section)
    }
}