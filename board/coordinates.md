Hexagons have the shape 

```
 /‾‾‾‾‾‾\
/        \
\        /
 \______/
```

Then we concatenate them so a triangle looks like this:

```
          /‾‾‾‾‾‾\
         /        \
 /‾‾‾‾‾‾\\        /
/        \\______/
\        //‾‾‾‾‾‾\
 \______//        \
         \        /
          \______/
```

The coordinates will be a tilted x-y, where X goes ↗️ and ➡️, and Y goes only ↘️.

So our triangle would have coords:

```
          /30‾‾‾‾\40
         /20      \41
 /10‾‾‾‾\\        /
/00      \\21____/31
\        //‾‾‾‾‾‾\
 \01____//11      \32
         \        /
          \12____/22
```

This coordinate system should aim to be internal to this package, but it won't start like that initially.