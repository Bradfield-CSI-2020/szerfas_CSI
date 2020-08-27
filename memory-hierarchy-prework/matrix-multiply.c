/*
Naive code for multiplying two matrices together.

There must be a better way!
*/

#include <stdio.h>
#include <stdlib.h>

/*
  A naive implementation of matrix multiplication.

  DO NOT MODIFY THIS FUNCTION, the tests assume it works correctly, which it
  currently does
*/
void matrix_multiply(double **C, double **A, double **B, int a_rows, int a_cols,
                     int b_cols) {
  for (int i = 0; i < a_rows; i++) {
    for (int j = 0; j < b_cols; j++) {
      C[i][j] = 0;
      for (int k = 0; k < a_cols; k++)
        C[i][j] += A[i][k] * B[k][j];
    }
  }
}

//void fast_matrix_multiply(double **c, double **a, double **b, int a_height, int b_height, int b_width) {
//  return matrix_multiply(c, a, b, a_height, b_height, b_width);
//}

void fast_matrix_multiply(double **C, double **A, double **B, int a_rows, int a_cols, int b_cols) {
  for (int i = 0; i < a_rows; i++) {
    for (int j = 0; j < b_cols; j++) {
      C[i][j] = 0;
      for (int k = 0; k < a_cols; k++)
      // just testing the improvement if swap the order of the columns multiplied by one another
      // other ideas: use a temporary variable rather than constantly referring to memory, loop unroll
        C[i][j] += A[i][k] * B[i][k];
    }
  }
}
