function getRandomSymbol() {
  const r1 = (11 + Math.floor(Math.random() * 5)).toString(16);
  const r2 = Math.floor(Math.random() * 16).toString(16);
  const val = `30${r1}${r2}`;
  const unicode = String.fromCharCode(parseInt(val, 16));
  return unicode;
}

function Symbol(cursor, isCursor = true) {
  this.sym = getRandomSymbol();
  this.isCursor = isCursor ? isCursor : false;
  this.alpha = 1;
  this.cursor = cursor;

  this.pickRandomSymbol = () => {
    this.sym = getRandomSymbol();
  };

  this.drawSymbol = (ctx, x) => {
    const color = `rgba(0, 255, 0, ${this.alpha})`;
    ctx.font = `${letterSize}px serif`;
    ctx.fillStyle = color;
    // introduccing random change of stables.
    if (Math.random() > stableChange) {
      this.pickRandomSymbol();
    }
    ctx.fillText(this.sym, x, this.cursor);
  };
}
