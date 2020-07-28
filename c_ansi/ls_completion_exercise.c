#include <unistd.h>
#include <stdio.h>
#include <dirent.h>
#include <sys/stat.h>

/* Approach to ls implementation
If arguments, parse arguments for target directory
Else, get cwd to use as target directory
Open the target directory
Traverse the target directory to get list of directory entries;
For each directory entry, get the file_name and a piece of metadata, like it's size
Print the file name and its size in a reasonable format
*/

/* TODO: handle arguments passed in */
int main ()
{
    FILE *fp;
    DIR *dp;
    char *cwd;
    struct dirent *dir_entry;
    struct stat stbuf;
    char *file_name;

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
}
/* TODOs:
- break out into a few functions
- check if need to set function prototype and variable declarations separately after importing headers
*/