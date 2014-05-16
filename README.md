YDiskHas
========

A little utility that can tell if a given file has been already uploaded by
someone to Yandex.Disk by using using its [API's deduplication
check](http://api.yandex.com/disk/doc/dg/reference/put.xml).

Potential uses are described, for example, here: ["Side channels in cloud
services, the case of deduplication in cloud
storage"](http://www.pinkas.net/PAPERS/hps.pdf).

One example:

Step 1. Send a form to Alice and ask her to fill it and save it to her
Yandex.Disk.

*form.txt*:

    Alice, do you love Bob?  

    [ ] Yes  [ ] No


Step 2. Generate two files:

*form-yes.txt*:

    Alice, do you love Bob?  

    [x] Yes  [ ] No


*form-no.txt*:

    Alice, do you love Bob?

    [ ] Yes  [x] No

Step 3. Run `./ydiskhas.sh` on both files:

    $ ./ydiskhas.sh form-yes.txt YourYDiskLogin:YourPassword

    (╯°□°)╯︵ ┻━┻
    File does not exist on Yandex.Disk.


    $ ./ydiskhas.sh form-no.txt YourYDiskLogin:YourPassword

    ˙ ͜ʟ˙
    FILE EXISTS on Yandex.Disk!


...and discover that Alice doesn't love Bob :-(
