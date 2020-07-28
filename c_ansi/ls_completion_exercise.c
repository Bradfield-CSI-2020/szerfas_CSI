#include <unistd.h>
#include <stdio.h>
#include <dirent.h>
#include <sys/stat.h>

/* TODO: handle arguments passed in */
int main ()
{
//    #define MAXPATH 10000
    /*
    get cwd
    open the cwd file
    traverse the cwd file object to get list of file names stored within the file
    */
    FILE *fp;
    DIR *dp;
    char *cwd;
    struct dirent *dir_entry;
    struct stat stbuf;
    char *file_name;
//    char c;
//    char cwd[MAXPATH];

//    fp = fopen(cwd, "r");
//    c = getc(fp);
//    printf("%d\n", c);

    cwd = getcwd(NULL, 0);
    dp = opendir(cwd);

    printf("Printing directory contents of: %s\n", getcwd(NULL, 0));
    char file_name_header[10] = "file_name";
    char file_size_header[10] = "file_size";
    printf("%30s\t\t%30s\n", file_name_header, file_size_header);
    while ((dir_entry = readdir(dp)) != NULL) {
        /* TODO: if dir_entry refers to a directory, print name without size */
        printf("%30s\t\t", dir_entry->d_name);
        stat(dir_entry->d_name, &stbuf);
        printf("%30lld\n", stbuf.st_size);
    }

//    printf("%s\n", getcwd(NULL, 0));
}
