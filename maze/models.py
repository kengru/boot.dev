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
    
    if not self._win:
      return

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
    self._cells = []
    self._x1 = x1
    self._y1 = y1
    self._num_rows = num_rows
    self._num_cols = num_cols
    self._cell_size_x = cell_size_x
    self._cell_size_y = cell_size_y
    self._win = win
    if seed:
      random.seed(seed)

    self._create_cells()
    self._break_entrance_and_exit()
    self._break_walls_r(0,0)
    self._reset_cells_visited()
  

  def _animate(self):
    if self._win is None:
      return
    self._win.redraw()
    time.sleep(0.01)


  def _break_entrance_and_exit(self):
    i = len(self._cells)
    j = len(self._cells[0])
    self._cells[0][0].has_top_wall = False
    self._cells[i - 1][j - 1].has_bottom_wall = False
    

  def _break_walls_r(self, i, j):
    self._cells[i][j]._visited = True
    while True:
      next_index_list = []

      if i > 0 and not self._cells[i - 1][j]._visited:
        next_index_list.append((i - 1, j))
      if i < self._num_cols - 1 and not self._cells[i + 1][j]._visited:
        next_index_list.append((i + 1, j))
      if j > 0 and not self._cells[i][j - 1]._visited:
        next_index_list.append((i, j - 1))
      if j < self._num_rows - 1 and not self._cells[i][j + 1]._visited:
        next_index_list.append((i, j + 1))

      if len(next_index_list) == 0:
        self._draw_cell(i, j)
        return

      # randomly choose the next direction to go
      direction_index = random.randrange(len(next_index_list))
      next_index = next_index_list[direction_index]

      if next_index[0] == i + 1:
        self._cells[i][j].has_right_wall = False
        self._cells[i + 1][j].has_left_wall = False
      if next_index[0] == i - 1:
        self._cells[i][j].has_left_wall = False
        self._cells[i - 1][j].has_right_wall = False
      if next_index[1] == j + 1:
        self._cells[i][j].has_bottom_wall = False
        self._cells[i][j + 1].has_top_wall = False
      if next_index[1] == j - 1:
        self._cells[i][j].has_top_wall = False
        self._cells[i][j - 1].has_bottom_wall = False

      self._break_walls_r(next_index[0], next_index[1])


  def _create_cells(self):
    self._cells = []
    for i in range(self._num_cols):
      li = []
      for j in range(self._num_rows):
        c = Cell(self._win)
        li.append(c)
      self._cells.append(li)
    for i in range(self._num_cols):
      for j in range(self._num_rows):
        self._draw_cell(i, j)
    

  def _draw_cell(self, i, j):
    if self._win is None:
      return
    x1 = self._x1 + self._cell_size_x * j
    y1 = self._y1 + self._cell_size_y * i
    x2 = self._x1 + self._cell_size_x * (j + 1)
    y2 = self._y1 + self._cell_size_y * (i + 1)
    self._cells[i][j].draw(x1, y1, x2, y2)
    self._animate()
      

  def _reset_cells_visited(self):
    for i in range(self._num_cols):
      for j in range(self._num_rows):
        self._cells[i][j]._visited = False


  def solve(self):
    return self._solve_r(0, 0)


  def _solve_r(self, i, j):
    from utils import get_middle_cell
    self._animate()
    self._cells[i][j]._visited = True
    if i == self._num_cols - 1 and j == self._num_rows:
      return True
    
    if i > 0 and self._cells[i - 1][j].has_right_wall and not self._cells[i - 1][j]._visited:
      move = Line(get_middle_cell(self._cells[i][j]), self._cells[i - 1][j])
      self._win.draw_line(move, "green")
      val = self._solve_r(i - 1, j)
      if val:
        return val
      self._win.draw_line(move, "red")
    if i < self._num_cols - 1 and self._cells[i + 1][j].has_left_wall and not self._cells[i + 1][j]._visited:
      move = Line(get_middle_cell(self._cells[i][j]), self._cells[i + 1][j])
      self._win.draw_line(move, "green")
      val = self._solve_r(i + 1, j)
      if val:
        return val
      self._win.draw_line(move, "red")
    if j > 0 and self._cells[i][j - 1].has_bottom_wall and not self._cells[i][j - 1]._visited:
      move = Line(get_middle_cell(self._cells[i][j]), self._cells[i][j - 1])
      self._win.draw_line(move, "green")
      val = self._solve_r(i, j - 1)
      if val:
        return val
      self._win.draw_line(move, "red")
    if j < self._num_rows - 1 and self._cells[i][j + 1].has_top_wall and not self._cells[i][j + 1]._visited:
      move = Line(get_middle_cell(self._cells[i][j]), self._cells[i][j + 1])
      self._win.draw_line(move, "green")
      val = self._solve_r(i, j + 1)
      if val:
        return val
      self._win.draw_line(move, "red")