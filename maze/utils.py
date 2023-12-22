from models import Point, Cell

def get_middle_point(x1, x2, y1, y2) -> Point:
  return Point((x1 + x2) / 2, (y1 + y2) / 2)


def get_middle_cell(cell: Cell) -> Point:
  return Point((cell._x1 + cell._x2) / 2, (cell._y1 + cell._y2) / 2)