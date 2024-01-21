function drawLetter(ctx, x, y, letter) {
  ctx.font = '48px serif';
  ctx.fillText(letter, x, y);
}

function draw() {
  const canvas = document.getElementById('canvas');
  const ctx = canvas.getContext('2d');

  for (let i = 0; i < 12; i++) {
    for (let j = 0; j < 10; j++) {
      ctx.fillStyle = `rgb(${Math.floor(255 - 42.5 * i)}, ${Math.floor(255 - 42.5 * j)}, 0)`;
      ctx.fillRect(25 + j * 50, 25 + i * 50, 50, 50);
    }
  }
}

window.addEventListener('load', draw);