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

The coordinates will be x-y, where X goes only ➡️, but Y goes both ↘️ and ↙️.

So our triangle could have coords:

```
         /11‾‾‾\21
        /12     \22
 /02‾‾‾\\       /
/03     \\13___/23
\       //‾‾‾‾‾\
 \04___//14     \24
        \       /
         \15___/25
```

(Commas omitted here for clarity)

Notice that coords that add up to even numbers have a horizontal line to the right,
while coords that add up to odd nums have a horizontal line to the left.

## Path/edge coordinates

A path coordinate is simply a pair of intersections.
For example, the intersections "2,1" and "3,1" are connected via "2,1-3,1".
Canonical ordering is ascending.

## Tile coordinates

Each tile is represented by the intersection on its leftmost edge.
Not all edge coordinates have corresponding tiles.

<details>
<summary>
A full Catan board, with the tile coordinates, would look like this:
</summary>

```
                 /‾‾‾‾‾\
                /  2,1  \
         /‾‾‾‾‾\\       //‾‾‾‾‾\
        /  1,2  \\_____//  3,2  \
 /‾‾‾‾‾\\       //‾‾‾‾‾\\       //‾‾‾‾‾\
/  0,3  \\_____//  2,3  \\_____//  4,3  \
\       //‾‾‾‾‾\\       //‾‾‾‾‾\\       /
 \_____//  1,4  \\_____//  3,4  \\_____/
 /‾‾‾‾‾\\       //‾‾‾‾‾\\       //‾‾‾‾‾\
/  0,5  \\_____//  2,5  \\_____//  4,5  \
\       //‾‾‾‾‾\\       //‾‾‾‾‾\\       /
 \_____//  1,6  \\_____//  3,6  \\_____/
 /‾‾‾‾‾\\       //‾‾‾‾‾\\       //‾‾‾‾‾\
/  0,7  \\_____//  2,7  \\_____//  4,7  \
\       //‾‾‾‾‾\\       //‾‾‾‾‾\\       /
 \_____//  1,8  \\_____//  3,8  \\_____/
        \       //‾‾‾‾‾\\       /
         \_____//  2,9  \\_____/
                \       /
                 \_____/
```

</details>

## What if a coordinate doesn't exist?

Then it doesn't exist. The coordinate system assumes an infinite board,
which is then filtered.