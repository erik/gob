/* The following complete B program, if compiled and put on your
  file "hstar", will act as an ascii file copy routine; the command
  at "SYSTEM?" level:

   /hstar file1 file2

  will copy file1 to file2. */

main () {
   auto j,s[20],t[20];
   reread(); /* get command line */
   getstr(s); /* put into s */
   j = getarg(t,s,0); /* skip H* name */
   j = getarg(t,s,j); /* filel */
   openr( 5,t );
   getarg(t,s,j); /* file2 */
   openw( 6,t );
   while( putchar( getchar() ) != '*e' ) ; /* copy contents */
   }
