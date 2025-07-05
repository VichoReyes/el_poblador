## Basic look

Hexagons have the shape

```
 /‾‾‾‾‾‾\
/        \

\        /
 \______/
```

Then we concatenate them so a triangle looks like this:

```
asd          asd ====
 \\            /‾‾‾‾‾‾\
  \\          /        \
  asd ==== asd
  ///‾‾‾‾‾‾\  \        /
 ///        \  \______/
asd          sdf
 \\\        /  /‾‾‾‾‾‾\
  \\\______/  /        \
  asd ==== asd
              \        /
               \______/
```

Which leaves extra space for roads and settlements.

## Rendering responsibilities

There are 4 entity types that can print themselves
- Tiles
- Crossings
- Paths
- Padding

### Tiles

The tiles are responsible for the area delimited by their borders (between the slashes)
For the middle row, they're responsible for a whitespace of length 10.  
So their 5 rows take up the following spaces:
- 8
- 10
- 10
- 10
- 8

### Crossings

Crossings are 3-character spaces, in a single line.

### Paths

Diagonal paths are 2x2 rhomboids, and horizontal ones are 6x1

```
 //
//

(quotes not included, notice the spaces)
" ==== "
```

### Padding

Padding is responsible for making the whole thing fit into a square.

### Overall

Iterate over crossing coordinates, left to right.
Every crossing puts itself on the line, plus the item to its right.
(can be a road or a tile).
When the item is a tile, the crossing also takes responsibility for the
paths going up and down, plus the whole tile.

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