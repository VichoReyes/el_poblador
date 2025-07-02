## Basic look

Hexagons have the shape

```
 /‾‾‾‾‾\
/       \
\       /
 \_____/
```

Then we concatenate them so a triangle looks like this:

```
         /‾‾‾‾‾\
        /       \
 /‾‾‾‾‾\\       /
/       \\_____/
\       //‾‾‾‾‾\
 \_____//       \
        \       /
         \_____/
```

# Coordinates

All coordinate systems should aim to be internal to this package,
but it won't start like that initially.

## Intersection coordinates

The coordinates will be a tilted x-y, where X goes ↗️ and ➡️, and Y goes only ↘️.

So our triangle would have coords:

```
         /30‾‾‾\40
        /20     \41
 /10‾‾‾\\       /
/00     \\21___/31
\       //‾‾‾‾‾\
 \01___//11     \32
        \       /
         \12___/22
```

(Commas omitted here for clarity)

Notice that coords that add up to even numbers have a horizontal line to the left,
while coords that add up to odd nums have a horizontal line to the right.

## Path/edge coordinates

A path coordinate is simply a pair of intersections.
For example, the intersections "2,1" and "3,1" are connected via "2,1-3,1".
Canonical ordering is ascending.

## Tile coordinates

Each tile is represented by the intersection on the left side of its top edge.
Not all edge coordinates have corresponding tiles.

<details>
<summary>
A full Catan board, with the tile coordinates, would look like this:
</summary>

```
                 /‾‾‾‾‾\
                /  5,0  \
         /‾‾‾‾‾\\       //‾‾‾‾‾\
        /  3,0  \\_____//  6,1  \
 /‾‾‾‾‾\\       //‾‾‾‾‾\\       //‾‾‾‾‾\
/  1,0  \\_____//  4,1  \\_____//  7,2  \
\       //‾‾‾‾‾\\       //‾‾‾‾‾\\       /
 \_____//  2,1  \\_____//  5,2  \\_____/
 /‾‾‾‾‾\\       //‾‾‾‾‾\\       //‾‾‾‾‾\
/  0,1  \\_____//  3,2  \\_____//  6,3  \
\       //‾‾‾‾‾\\       //‾‾‾‾‾\\       /
 \_____//  1,2  \\_____//  4,3  \\_____/
 /‾‾‾‾‾\\       //‾‾‾‾‾\\       //‾‾‾‾‾\
/ -1,2  \\_____//  2,3  \\_____//  5,4  \
\       //‾‾‾‾‾\\       //‾‾‾‾‾\\       /
 \_____//  0,3  \\_____//  3,4  \\_____/
        \       //‾‾‾‾‾\\       /
         \_____//  1,4  \\_____/
                \       /
                 \_____/
```

</details>

## What if a coordinate doesn't exist?

Then it doesn't exist. The coordinate system assumes an infinite board,
which is then filtered.