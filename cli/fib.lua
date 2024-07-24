
function fib(n)
    if n <= 0 then
        return 0
    end
    local t1 = 0
    local t2 = 1
    local nextTerm = t2
    for i = 3, n do
        t1 = t2
        t2 = nextTerm
        nextTerm = t1 + t2
    end
    return nextTerm
end

print(fib(123))
