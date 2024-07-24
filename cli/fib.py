def fib(n):
    if n <= 0:
        return 0
    t1 = 0
    t2 = 1
    nextTerm = t2
    for i in range(3, n + 1):
        t1 = t2
        t2 = nextTerm
        nextTerm = t1 + t2
    return nextTerm

print(fib(123))
