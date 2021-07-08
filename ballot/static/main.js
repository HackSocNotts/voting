var candidates, ballot

function load() {
    getCandidates()
}

function getCandidates() {
    var req = new XMLHttpRequest()

    req.onreadystatechange = function() {
        if (this.readyState != 4) return

        if (this.status == 200) {
            candidates = JSON.parse(req.responseText)
            ballot = {}
            for (var i = 0; i < candidates.length; i++) {
                ballot[i] = []
            }
            renderForm()
        } else {
            // TODO: handle error
        }
    }

    req.open("GET", "/candidates/")
    req.send()
}

function renderForm() {
    for (var i = 0; i < candidates.length; i++) {
        const pos = i

        var position = candidates[i]

        var section = document.createElement("section")

        var h2 = document.createElement("h2")
        h2.innerHTML = position.role
        section.appendChild(h2)

        var ul = document.createElement("ul")
        ul.classList.add("choices")
        section.appendChild(ul)

        for (var j = 0; j < position.candidates.length; j++) {
            const cand = j

            var candidate = position.candidates[j]

            var li = document.createElement("li")
            li.classList.add("choice")
            if (candidate == "Re-open Nominations") {
                li.classList.add("ron")
            }

            li.onclick = function() { selectCandidate(pos, cand) }

            ul.appendChild(li)

            var rank = document.createElement("span")
            rank.classList.add("rank")
            rank.innerHTML = "-"
            rank.id = "rank-" + i + "-" + j
            li.appendChild(rank)

            var name = document.createElement("span")
            name.classList.add("name")
            name.innerHTML = candidate
            li.appendChild(name)
        }

        var clear = document.createElement("button")
        clear.innerHTML = "Clear choices"
        clear.classList.add("clear")
        clear.onclick = function() { clearBallot(pos) }
        section.appendChild(clear)

        document.getElementById("form").appendChild(section)
    }
}

function selectCandidate(pos, candidate) {
    if (ballot[pos].indexOf(candidate) >= 0) {
        return
    }

    ballot[pos].push(candidate)
    document.getElementById("rank-" + pos + "-" + candidate).innerHTML = ballot[pos].length

    if (ballot[pos].length == candidates[pos].candidates.length) {
        document.getElementById("rank-" + pos + "-" + candidate).parentElement.parentElement.classList.add("complete")
    }
}

function clearBallot(pos) {
    for (var vote of ballot[pos]) {
        document.getElementById("rank-" + pos + "-" + vote).innerHTML = "-"
    }

    ballot[pos] = []
    document.getElementById("rank-" + pos + "-0").parentElement.parentElement.classList.remove("complete")
}