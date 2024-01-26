const width = 800;
const height = 600;
const separation = width / 15;
const separationY = 30;
const velocity = 2;
const startingAlpha = 0.3;
const decreasingAlpha = 0.008;
const letterSize = 15;
const generationSpeed = 0.97; // less while closest to 0
const stableChange = 0.9995; // changes already stable symbols

let linesArr = [];

function setup() {
  // fill the initial array.
  for (let i = 0; i < separation; i++) {
    linesArr.push(undefined);
  }

  window.requestAnimationFrame(draw);
}

function draw() {
  const ctx = document.getElementById("canvas").getContext("2d");
  ctx.clearRect(0, 0, 800, 600); // clearing the canvas

  // Random generation of lines
  for (let i = 0; i < separation; i++) {
    if (Math.random() > generationSpeed) {
      if (linesArr[i] == undefined) {
        const posX = i * (width / separation);
        linesArr[i] = new Line(posX);
      }
    }
  }
  // console.log(
  //   linesArr.reduce((previous, current) => {
  //     return current ? previous + 1 : previous;
  //   }, 0),
  // );

  // updating and drawing all existing lines
  for (let i = 0; i < separation; i++) {
    const line = linesArr[i];
    if (line !== undefined) {
      if (line.running !== "off") {
        if (line.running !== "dying") {
          line.update();
        }
      } else {
        linesArr[i] = undefined;
      }
      line.draw(ctx);
      if (
        line.syms.length > 0 &&
        line.syms.every((s) => {
          return s.alpha < 0.01;
        })
      ) {
        line.running = "off";
      }
    }
  }
  // sending a new request for animation frame
  window.requestAnimationFrame(draw);
}

setup();
