import time

def sum_to_n(n):
    start = time.time()

    result = 0
    for i in range(0, n):
        result += i

    end = time.time()

    return result, end - start

def arithmetic_sum(n):
    start = time.time()
    total = n * (n + 1) // 2
    end = time.time()
    return total, end - start

def test_errors():
    try:
        y = '5'
        print('execute some code')
        raise e
    except Exception as e:
        print('executes some other code and raise an error')
        raise e
    finally:
        print('will I still get run and raise an error?')
        raise e
