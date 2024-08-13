#include <iostream>
#include <string.h>
using namespace std;

void function(char *input){
    char buffer[1];
    strcpy(buffer, input); 
}

int main(){
    const char *string = "hello world";
    function((char *)string); 
    return 0;
}
