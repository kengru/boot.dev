from tkinter import Tk, BOTH, Canvas

class Window():
  def __init__(self, width, height):
    self.__root = Tk()
    self.__root.title("Maze Solver")
    self.__root.geometry(f"{width}x{height}")
    self.__root.protocol("WM_DELETE_WINDOW", self.close)
    self.canvas = Canvas(self.__root, bg="white", height=height, width=width)
    self.canvas.pack(fill=BOTH, expand=1)
    self.running = False
  
  def redraw(self):
    self.__root.update_idletasks()
    self.__root.update()
  
  def draw_line(self, line, fill = "red"):
    line.draw(self.canvas, fill)
  
  def draw_cell(self, cell):
    cell.draw(self.canvas)
  
  def wait_for_close(self):
    self.running = True
    while self.running:
      self.redraw()
    print("window closed ...")
    
  def close(self):
    self.running = False