from models import Point

def get_middle_point(x1, x2, y1, y2) -> Point:
  return Point((x1 + x2) / 2, (y1 + y2) / 2)