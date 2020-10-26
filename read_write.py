f = open('/dev/urandom', 'rb')
res = f.read(2)
print(str(res))
f.close()

