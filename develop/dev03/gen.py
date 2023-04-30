import os
import datetime

start = datetime.datetime.now()
step = lambda x: datetime.timedelta(minutes=5*x)
rows = 16

with open("sample", "w") as file:
    for i in range(rows):
        stamp = start - step(i)
        row = f'[{datetime.datetime.strftime(stamp, "%d %b %Y %H:%M")}]\tКакие то логи\t№{i+1}\n'
        file.write(row) 

print("Done")    

# ab a3
# ae a5
# ad a4
# ac a2
# aa a11 
# aa a1

