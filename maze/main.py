from window import Window
from models import Maze

def main():
  rows = 12
  cols = 16
  margin = 50
  screen_x = 800
  screen_y = 800
  cell_size_x = (screen_x - 2 * margin) / rows
  cell_size_y = (screen_y - 2 * margin) / cols

  win = Window(screen_x, screen_y)
  mz = Maze(margin, margin, rows, cols, cell_size_x, cell_size_y, win, 11)
  mz.solve()

  win.wait_for_close()

main()