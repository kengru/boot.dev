import unittest
from unittest.mock import Mock
from models import Maze
from window import Window


class Tests(unittest.TestCase):
  def test_maze_create_cells(self):
    num_rows = 6
    num_cols = 4
    m1 = Maze(0, 0, num_rows, num_cols, 10, 10, Mock(spec=Window))
    self.assertEqual(len(m1._cells), num_rows)
    self.assertEqual(len(m1._cells[0]), num_cols)


  def test_maze_create_cells_2(self):
    num_rows = 2
    num_cols = 12
    m1 = Maze(0, 0, num_rows, num_cols, 10, 10, Mock(spec=Window))
    self.assertEqual(len(m1._cells), num_rows)
    self.assertEqual(len(m1._cells[0]), num_cols)
  
  
  def test_maze_break_e_a_e(self):
    num_rows = 3
    num_cols = 4
    m1 = Maze(0, 0, num_rows, num_cols, 10, 10, Mock(spec=Window))
    m1._break_entrance_and_exit()
    self.assertEqual(m1._cells[0][0].has_top_wall, False)
    self.assertEqual(m1._cells[num_rows - 1][num_cols - 1].has_bottom_wall, False)


if __name__ == "__main__":
  unittest.main()
