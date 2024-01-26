function constrain(val, low, high) {
  if (val < low) {
    return low;
  }
  if (val > high) {
    return high;
  }
  return val;
}

function Line(x) {
  this.x = x;
  this.cursor = 0;
  this.running = "on"; // on, off, dying
  this.limit = constrain(
    height * 0.1 + Math.random() * height,
    height * 0.1,
    height - separationY,
  );
  this.symbol = new Symbol(0);
  this.syms = [];

  this.draw = (ctx) => {
    // draw cursor ocasionally
    if (this.running === "on") {
      this.symbol.cursor = this.cursor;
      this.symbol.drawSymbol(ctx, this.x, 1);
      if (this.cursor % 10 === 0) {
        this.symbol.pickRandomSymbol();
      }
    }

    // draw syms
    for (const sym of this.syms) {
      // sym.alpha = 1 - sym.cursor / this.limit;
      sym.alpha -= decreasingAlpha;
      sym.drawSymbol(ctx, this.x);
    }
  };

  this.update = () => {
    // updating cursor
    this.cursor += velocity;

    // leaving a symbol behind
    if (this.cursor % separationY === 0) {
      const s = new Symbol(this.cursor);
      s.alpha = startingAlpha + this.cursor / this.limit;
      this.syms.push(s);
    }

    if (this.cursor > this.limit) {
      this.running = "dying";
    }
  };
}
