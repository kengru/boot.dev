def main():
  p = "books/frankenstein.txt"
  print(f"---- Report on {p} ----")
  
  text = get_text(p)
  amount = len(text.split())
  print(f"There are {amount} words in this document\n")

  d = {}
  for word in text.split():
    for l in word:
      if l.isalpha():
        comp = l.lower()
        if comp in d:
          d[comp] += 1
        else:
          d[comp] = 1
  final_dict = list(d.items())
  final_dict.sort(key=lambda x: x[1], reverse=True)

  for v in final_dict:
    print(f"The '{v[0]}' character was found {v[1]} times")
  
  print("---- End of report ---")

def get_text(p):
  return open(p).read()

main()