#include <stdio.h>
#include <unistd.h>
#include <sys/reboot.h>

int main(void) {
    printf("Hello, world!\n");
    reboot(0x4321fedc);
    return 0;
}
