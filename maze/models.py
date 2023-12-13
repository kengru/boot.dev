from tkinter import BOTH
import time
import random


class Point():
  def __init__(self, x, y):
    self.x = x
    self.y = y


class Line():
  def __init__(self, p1, p2):
    self.p1 = p1
    self.p2 = p2
  
  def draw(self, canvas, fill = "red"):
    canvas.create_line(
      self.p1.x, self.p1.y, self.p2.x, self.p2.y, fill=fill, width=4
    )
    canvas.pack(fill=BOTH, expand=1)


class Cell():
  def __init__(self, window = None):
    self.has_top_wall = True
    self.has_right_wall = True
    self.has_bottom_wall = True
    self.has_left_wall = True
    self._x1 = None
    self._x2 = None
    self._y1 = None
    self._y2 = None
    self._visited = False
    self._win = window
  
  def draw(self, x1, y1, x2, y2):
    self._x1 = x1
    self._x2 = x2
    self._y1 = y1
    self._y2 = y2
    
    color1 = "red" if self.has_top_wall else "white"
    line = Line(Point(x1, y1), Point(x2, y1))
    self._win.draw_line(line, color1)
    
    color2 = "red" if self.has_right_wall else "white"
    line = Line(Point(x2, y1), Point(x2, y2))
    self._win.draw_line(line, color2)
  
    color3 = "red" if self.has_bottom_wall else "white"
    line = Line(Point(x2, y2), Point(x1, y2))
    self._win.draw_line(line, color3)
  
    color4 = "red" if self.has_left_wall else "white"
    line = Line(Point(x1, y1), Point(x1, y2))
    self._win.draw_line(line, color4)
  
  def draw_move(self, to_cell, undo=False):
    from utils import get_middle_point
    fill = "gray" if undo else "red"
    p1 = get_middle_point(self._x1, self._x2, self._y1, self._y2)
    p2 = get_middle_point(to_cell._x1, to_cell._x2, to_cell._y1, to_cell._y2)
    l = Line(p1, p2)
    self._win.draw_line(l, fill)


class Maze():
  def __init__(self, x1, y1, num_rows, num_cols, cell_size_x, cell_size_y, win = None, seed = None):
    self._x1 = x1
    self._y1 = y1
    self._num_rows = num_rows
    self._num_cols = num_cols
    self._cell_size_x = cell_size_x
    self._cell_size_y = cell_size_y
    self._win = win
    self._seed = random.seed(seed) if seed else None
    self._create_cells()
  

  def _animate(self):
    self._win.redraw()
    time.sleep(0.02)


  def _break_entrance_and_exit(self):
    i = len(self._cells)
    j = len(self._cells[0])
    self._cells[0][0].has_top_wall = False
    self._cells[i - 1][j - 1].has_bottom_wall = False
    

  def _break_walls_r(self, i, j):
    current = self._cells[i][j]
    current._visited = True
    while True:
      to_visit = []
      for x in range(min(0, i - 1), max()):
        pass



  def _create_cells(self):
    self._cells = []
    for i in range(self._num_rows):
      li = []
      for j in range(self._num_cols):
        c = Cell(self._win)
        li.append(c)
      self._cells.append(li)
    for i in range(self._num_rows):
      for j in range(self._num_cols):
        self._draw_cell(self._cells[i][j], i, j)
    

  def _draw_cell(self, cell: Cell, i, j):
    x1 = self._x1 + self._cell_size_x * j
    y1 = self._y1 + self._cell_size_y * i
    x2 = self._x1 + self._cell_size_x * (j + 1)
    y2 = self._y1 + self._cell_size_y * (i + 1)
    cell.draw(x1, y1, x2, y2)
    self._animate()
      
