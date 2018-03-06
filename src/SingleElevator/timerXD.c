#include<stdio.h>
#include<sys/time.h>

 static double get_wall_time(void){
    struct timeval time;
    gettimeofday(&time, NULL);
    return (double)time.tv_sec + (double)time.tv_usec * .000001;
}

void main() {
    double hei = get_wall_time();
 
   printf("%f\n", hei);
}
