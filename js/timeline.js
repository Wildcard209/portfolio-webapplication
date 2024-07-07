const canvas = document.getElementById('timelineCanvas');
const ctx = canvas.getContext('2d');

const dpr = window.devicePixelRatio || 1;
const canvasHeight = 200;
canvas.height = canvasHeight * dpr;
ctx.scale(dpr, dpr);

const startDate = new Date('2010-01-01');
const endDate = new Date('2020-12-31');

function calculateTotalMonths(startDate, endDate) {
    let totalMonths = (endDate.getFullYear() - startDate.getFullYear()) * 12;
    totalMonths += endDate.getMonth() - startDate.getMonth();
    return totalMonths;
}

const totalMonths = calculateTotalMonths(startDate, endDate) + 1; // +1 to include the end month

function resizeCanvas() {
    const parentWidth = canvas.parentElement.clientWidth;
    canvas.width = parentWidth * dpr;
    ctx.scale(dpr, dpr);
    drawTimeline(parentWidth);
}

function drawTimeline(parentWidth) {
    ctx.clearRect(0, 0, canvas.width, canvas.height);

    const lineY = canvasHeight / 2; // Center of the canvas for the timeline
    const padding = 10;
    const availableWidth = parentWidth - padding * 2;
    const monthSpacing = availableWidth / totalMonths;

    const endX = padding + monthSpacing * totalMonths;

    ctx.beginPath();
    ctx.moveTo(padding, lineY);
    ctx.lineTo(endX, lineY);
    ctx.strokeStyle = '#000';
    ctx.lineWidth = 2;
    ctx.stroke();

    let currentX = padding;

    drawTick(currentX, lineY, true);

    let currentDate = new Date(startDate);
    currentX += monthSpacing;

        for (let i = 1; i < totalMonths; i++) {
            if (currentDate.getMonth() === 11) {
                drawTick(currentX, lineY, true);
            } else {
                if(canvas.width > 600) {
                    drawTick(currentX, lineY, false);
                }
            }

            currentDate.setMonth(currentDate.getMonth() + 1);
            currentX += monthSpacing;
        }
    drawTick(endX, lineY, true);
}

function drawTick(x, y, isMajor) {
    ctx.beginPath();
    if (isMajor) {
        ctx.moveTo(x, y - 20); // Major tick height
        ctx.lineTo(x, y + 20);
        ctx.lineWidth = 2;
    } else {
        ctx.moveTo(x, y - 10); // Minor tick height
        ctx.lineTo(x, y + 10);
        ctx.lineWidth = 1;
    }
    ctx.strokeStyle = '#000';
    ctx.stroke();
}

resizeCanvas();

window.addEventListener('resize', resizeCanvas);