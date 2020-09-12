#include "vec.h"
#include "stdio.h"

// 9ms with length optimization
// 3ms after changing pointer memory references to avodi memory aliases
// could I change the function signature to prevent use of multiple pointers passed in such a way that could lead to memory aliasing
// loop unrolling



//data_t dotproduct(vec_ptr u, vec_ptr v) {
//   data_t sum = 0, u_val, v_val;
//
//   for (long i = 0; i < vec_length(u); i++) { // we can assume both vectors are same length
//        get_vec_element(u, i, &u_val);
//        get_vec_element(v, i, &v_val);
//        sum += u_val * v_val;
//   }
//   return sum;
//}


data_t dotproduct(vec_ptr u, vec_ptr v) {
   data_t sum0 = 0;
   data_t sum1 = 0;
   data_t sum2 = 0;
   data_t sum3 = 0;

   data_t *up = get_vec_start(u);
   data_t *vp = get_vec_start(v);

   long length = vec_length(u);

   long i = 0;
   for (; i < length; i += 4) { // we can assume both vectors are same length
        sum0 += *(up+i) * *(vp+i);
        sum1 += *(up+i+1) * *(vp+i+1);
        sum2 += *(up+i+2) * *(vp+i+2);
        sum3 += *(up+i+3) * *(vp+i+3);
   }

   for (; i < length; i++) {
      sum0 += *(up+i) * *(vp+i);
    }

   return sum0 + sum1 + sum2 + sum3;
}


//data_t dotproduct(vec_ptr u, vec_ptr v) {
//   data_t u_val, v_val;
//   long length = vec_length(u);
//
//   data_t sum0 = 0;
//   data_t sum1 = 0;
//   data_t sum2 = 0;
//   data_t sum3 = 0;
//   data_t sum4 = 0;
//   data_t sum5 = 0;
//   data_t sum6 = 0;
//
//   long i = 0;
//
//   for (; i < length - 6; i += 7) { // we can assume both vectors are same length
//        u_val = u->data[i];
//        v_val = v->data[i];
//        sum0 += u_val * v_val;
//
//        u_val = u->data[i+1];
//        v_val = v->data[i+1];
//        sum1 += u_val * v_val;
//
//        u_val = u->data[i+2];
//        v_val = v->data[i+2];
//        sum2 += u_val * v_val;
//
//        u_val = u->data[i+3];
//        v_val = v->data[i+3];
//        sum3 += u_val * v_val;
//
//        u_val = u->data[i+4];
//        v_val = v->data[i+4];
//        sum4 += u_val * v_val;
//
//        u_val = u->data[i+5];
//        v_val = v->data[i+5];
//        sum5 += u_val * v_val;
//
//        u_val = u->data[i+6];
//        v_val = v->data[i+6];
//        sum5 += u_val * v_val;
//
//   }
//
//   for (; i < length; i++) {
//       sum0 += u->data[i] * v->data[i];
//     }
//
//   return sum0 + sum1 + sum2 + sum3 + sum4 + sum5 + sum6;
//}
//
//data_t dotproduct(vec_ptr u, vec_ptr v) {
//   data_t sum = 0, u_val, v_val;
//   long *up = &u_val;
//   long *vp = &v_val;
//   long length = vec_length(u);
//   get_vec_element(u, 0, up);
//   get_vec_element(v, 0, vp);
//   sum += *up * *vp;
//
//   printf("up is %p\n", (void *) up);
//   printf("vp is %p\n", (void *) vp);
//
//
////   u_vec_element
////   v_vec_elemnt
//
//   for (long i = 1; i < length; i++) { // we can assume both vectors are same length
//        up++;
//        vp++;
//        if (i < 10) {
//            printf("up is %p\n", (void *) up);
//            printf("vp is %p\n", (void *) vp);
//        }
//        sum += *up * *vp;
//   }
//   return sum;
//}
//
