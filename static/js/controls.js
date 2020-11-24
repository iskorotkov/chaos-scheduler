stepsControl = document.getElementById("steps-control")
stepsIndicator = document.getElementById("steps-indicator")

stepsControl.oninput = function () {
    stepsIndicator.textContent = stepsControl.value
}
stepsIndicator.textContent = stepsControl.value
