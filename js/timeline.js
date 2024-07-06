const canvas = document.getElementById('timelineCanvas');
const ctx = canvas.getContext('2d');

// Adjust for device pixel ratio
const dpr = window.devicePixelRatio || 1;
canvas.width = 1000 * dpr;
canvas.height = 200 * dpr;
canvas.style.width = '1000px';
canvas.style.height = '200px';
ctx.scale(dpr, dpr);

const startDate = new Date('2010-03-21');
const endDate = new Date('2020-07-20');

const canvasWidth = 1000; // Width in pixels
const canvasHeight = 200; // Height in pixels
const lineY = canvasHeight / 2; // Center of the canvas for the timeline

// Function to calculate the number of months between two dates
function calculateTotalMonths(startDate, endDate) {
    let totalMonths = (endDate.getFullYear() - startDate.getFullYear()) * 12;
    totalMonths += endDate.getMonth() - startDate.getMonth();
    return totalMonths;
}

// Total number of months between start and end dates
const totalMonths = calculateTotalMonths(startDate, endDate) + 1; // +1 to include the end month

function drawTimeline() {
    ctx.clearRect(0, 0, canvas.width, canvas.height);
    ctx.beginPath();
    ctx.moveTo(50, lineY);
    ctx.lineTo(canvasWidth - 50, lineY);
    ctx.strokeStyle = '#000';
    ctx.lineWidth = 2;
    ctx.stroke();

    const monthSpacing = (canvasWidth - 100) / totalMonths; // Adjust for canvas padding
    let currentX = 50;

    // Draw start date tick with full date label
    drawFullDateTick(currentX, lineY, startDate);

    let currentDate = new Date(startDate);
    currentX += monthSpacing;

    for (let i = 1; i < totalMonths; i++) {
        if (currentDate.getMonth() === 11) {
            // Draw a year tick at the transition from December to January
            drawYearTick(currentX, lineY, currentDate.getFullYear() + 1);
        } else {
            // Draw a month tick for other months
            drawMonthTick(currentX, lineY);
        }

        // Move to the next month
        currentDate.setMonth(currentDate.getMonth() + 1);
        currentX += monthSpacing;
    }

    // Draw end date tick with full date label
    drawFullDateTick(currentX - monthSpacing, lineY, endDate);
}

function drawFullDateTick(x, y, date) {
    ctx.beginPath();
    ctx.moveTo(x, y - 20); // Major tick height
    ctx.lineTo(x, y + 20);
    ctx.strokeStyle = '#000';
    ctx.lineWidth = 2;
    ctx.stroke();

    ctx.font = '12px Arial';
    ctx.fillText(formatFullDate(date), x - 40, y + 40); // Label with full date
}

function drawYearTick(x, y, year) {
    ctx.beginPath();
    ctx.moveTo(x, y - 20); // Major tick height
    ctx.lineTo(x, y + 20);
    ctx.strokeStyle = '#000';
    ctx.lineWidth = 2;
    ctx.stroke();

    ctx.font = '12px Arial';
    ctx.fillText(year, x - 10, y + 40); // Label the year
}

function drawMonthTick(x, y) {
    ctx.beginPath();
    ctx.moveTo(x, y - 10); // Minor tick height
    ctx.lineTo(x, y + 10);
    ctx.strokeStyle = '#555';
    ctx.lineWidth = 1;
    ctx.stroke();
}

function formatFullDate(date) {
    const options = { year: 'numeric', month: 'short', day: 'numeric' };
    return date.toLocaleDateString(undefined, options);
}

drawTimeline();