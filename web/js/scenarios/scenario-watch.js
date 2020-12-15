const list = document.getElementById("workflows-list")

const values = window.location.href.split("/")

const host = values[2]
const namespace = values[values.length - 2]
const name = values[values.length - 1]

const ws = new WebSocket(`ws://${host}/scenarios/watch/${namespace}/${name}`)

window.addEventListener("beforeunload", () => {
    if (ws.readyState === WebSocket.OPEN || ws.readyState === WebSocket.CONNECTING) {
        ws.close()
    }
})

ws.addEventListener("message", ev => {
    let content = JSON.stringify(JSON.parse(ev.data), null, 2)
    list.innerHTML = `${list.innerHTML}\n<li><pre>${content}</pre></li>`
})
