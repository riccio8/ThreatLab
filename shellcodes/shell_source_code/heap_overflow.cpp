void function(){
  int* p = (int*)malloc(sizeof(int) * 1000);
  // free(p)
}

int main(){
  function();
  return 0;
}
