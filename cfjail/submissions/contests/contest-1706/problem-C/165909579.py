# infile = \
# '''
# 1
# 7
# 1000000000 1 1000000000 1 1000000000 1 1000000000
# '''
# __import__('sys').stdin = __import__('io').StringIO(infile.strip('\n'))


def solve():
    n = int(input())
    h = list(map(int, input().split()))

    set_a = 0
    for i in range(1, n-1, 2):
        set_a += max(max(h[i-1], h[i+1])-h[i] + 1, 0)

    ans = set_a

    if n % 2 == 0:
        set_b = 0
        for i in range(n-2, 1, -2):
            set_b += max(max(h[i-1], h[i+1])-h[i] + 1, 0)
            set_a -= max(max(h[i-2], h[i])-h[i-1] + 1, 0)
            ans = min(ans, set_a + set_b)

    return ans


t = int(input())
for _ in range(t):
    print(solve())