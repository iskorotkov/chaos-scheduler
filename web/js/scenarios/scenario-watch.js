const url = window.location.href
const values = url.split("/")

const namespace = values[values.length - 2]
const name = values[values.length - 1]

const list = document.getElementById("workflows-list")

const socket = new WebSocket(`ws://localhost:8811/scenarios/watch/${namespace}/${name}`)

window.addEventListener("unload", () => {
    if (socket.OPEN || socket.CONNECTING) {
        socket.close()
    }
})

socket.addEventListener("message", ev => {
    let content = JSON.stringify(JSON.parse(ev.data), null, 2)
    list.innerHTML = `${list.innerHTML}\n<li><pre>${content}</pre></li>`
})
