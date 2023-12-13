from window import Window
from models import Cell, Maze

def main():
  win = Window(800, 600)
  mz = Maze(40, 40, 10, 10, 40, 40, win)
  mz._create_cells()

  win.wait_for_close()

main()