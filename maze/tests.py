import unittest
from unittest.mock import Mock
from models import Maze, Cell
from window import Window
from utils import get_middle_point, get_middle_cell


class Tests(unittest.TestCase):
  def test_maze_create_cells(self):
    num_rows = 6
    num_cols = 4
    m1 = Maze(0, 0, num_rows, num_cols, 10, 10, Mock(spec=Window))
    self.assertEqual(len(m1._cells), num_cols)
    self.assertEqual(len(m1._cells[0]), num_rows)


  def test_maze_create_cells_2(self):
    num_rows = 2
    num_cols = 12
    m1 = Maze(0, 0, num_rows, num_cols, 10, 10, Mock(spec=Window))
    self.assertEqual(len(m1._cells), num_cols)
    self.assertEqual(len(m1._cells[0]), num_rows)
  
  
  def test_maze_break_e_a_e(self):
    num_rows = 3
    num_cols = 4
    m1 = Maze(0, 0, num_rows, num_cols, 10, 10, Mock(spec=Window))
    m1._break_entrance_and_exit()
    self.assertEqual(m1._cells[0][0].has_top_wall, False)
    self.assertEqual(m1._cells[num_cols - 1][num_rows - 1].has_bottom_wall, False)
  
  
  def test_maze_reset_cells_visited(self):
    num_rows = 3
    num_cols = 4
    m1 = Maze(0, 0, num_rows, num_cols, 10, 10, Mock(spec=Window))
    m1._cells[1][1]._visited = True
    m1._reset_cells_visited()
    self.assertEqual(m1._cells[1][1]._visited, False)


  def test_maze_reset_cells_visited_after_recursive(self):
    num_rows = 3
    num_cols = 4
    m1 = Maze(0, 0, num_rows, num_cols, 10, 10, Mock(spec=Window))
    m1._break_walls_r(0, 0)
    m1._reset_cells_visited()
    self.assertEqual(m1._cells[1][1]._visited, False)

  
  def test_get_middle_point(self):
    x1 = 20
    x2 = 30
    y1 = 30
    y2 = 40
    p = get_middle_point(x1, x2, y1, y2)
    self.assertEqual(p.x, 25)
    self.assertEqual(p.y, 35)
  
  def test_get_middle_cell(self):
    c = Cell()
    c.draw(30, 50, 60, 120)
    p = get_middle_cell(c)
    self.assertEqual(p.x, 45)
    self.assertEqual(p.y, 85)


if __name__ == "__main__":
  unittest.main()
