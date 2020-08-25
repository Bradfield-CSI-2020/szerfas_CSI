#include <stdint.h>
#include <stdio.h>
#include <time.h>

#define TEST_LOOPS 10000000

uint64_t pagecount(uint64_t memory_size, uint64_t page_size) {
//  page_size >>= 1;
//  return memory_size>> page_size;
  return memory_size / page_size;
}


// decimal
//page_size: 100
//memory_size: 10000 <-- 4 zeroes
//
//memory_size / page_size: 100 <-- 2 zeroes
//
//11 11 11 11 11 -> 1 11 11 11 11



uint64_t empty_func(uint64_t x, uint64_t y) {
  return 0;
}

int main (int argc, char** argv) {
  clock_t baseline_start, baseline_end, test_start, test_end;
  uint64_t memory_size, page_size;
  double clocks_elapsed, time_elapsed;
  int i, ignore = 0;

  uint64_t msizes[] = {1L << 32, 1L << 40, 1L << 52};
  uint64_t psizes[] = {1L << 12, 1L << 16, 1L << 32};

  baseline_start = clock();
  for (i = 0; i < TEST_LOOPS; i++) {
    memory_size = msizes[i % 3];
    page_size = psizes[i % 3];
    uint64_t empty_count = empty_func(0, 0);
    ignore +=
        memory_size + page_size + empty_count; // so that this loop isn't just optimized away
  }
  baseline_end = clock();

  test_start = clock();
  for (i = 0; i < TEST_LOOPS; i++) {
    memory_size = msizes[i % 3];
    page_size = psizes[i % 3];
    uint64_t pcount = pagecount(memory_size, page_size);
    ignore += memory_size + page_size + pcount;
  }
  test_end = clock();

  clocks_elapsed = test_end - test_start - (baseline_end - baseline_start);
  time_elapsed = clocks_elapsed / CLOCKS_PER_SEC;

  printf("%.2fs to run %d tests (%.2fns per test)\n", time_elapsed, TEST_LOOPS,
         time_elapsed * 1e9 / TEST_LOOPS);
  return ignore;
}

/* expected compiler instructions

pagecount func:
mov rdi -> rax      move memory size to return value
rax / rsi           divide return value by page size and return

main

put test_loops constant in a register
put baseline_start, baseline_end, test_start, test_end into registers -> initialize all to zero
put memory_size, page_size into registers -> initialize all to zero
put ignore in a register and initialize to zero
(it's fine to overwrite rdi and rsi since not used)

create space on the stack for a couple of arrays, each having three slots a quad in length, so would increment stack counter by 2 * 3 * 8 = 48
run bitst or some other bitwise operation on all six array values a instructed

call clock() and store result in baseline_start
initialize a counter to zero in a register (now using nine registers, so likely will have had to save some fo the callee saved on the stack to make room)
run first loop in jump-to-middle style, something like:
retrieve msize[index] from memory and set to memory_size register
same for page_size
accumulate ignore
check test condition, jl to start of loop

call clock() and set to baseline_end

call clock() and set to test_start
enter second loop, identical to first except call pagecount after setting memory_size and page_size
call clock and set to test_end

arithmetic across registers to get clock_elapsed and time_elapsed
create string array
call to printf passing in memory reference to string array

return ignore accumulator

*/