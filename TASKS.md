# Project Tasks

## Future Tasks

### Implement Robber Mechanics

**Description**: Implement the robber mechanics that are triggered when a player rolls a 7 or plays a knight card.

**Current State**: 
- Knight functionality is implemented and ready
- Phase continuation system is in place
- Robber mechanics are not yet implemented

**Suggested Approach**: 
Use the existing phase continuation system to handle robber phases and return to the appropriate continuation phase after robbery resolution.
Notice how the Phase's BoardCursor can be used to render a cursor on the tiles of the map when selecting a tile to place the robber.

**Estimated Complexity**: Medium
**Priority**: High (next major feature to implement)
