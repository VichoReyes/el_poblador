# Project Tasks

## Future Tasks

### ✅ Robber Mechanics - COMPLETED

**Description**: Robber mechanics have been fully implemented.

**Current State**: 
- ✅ Knight functionality is implemented and triggers robber placement
- ✅ Phase continuation system is in place
- ✅ `phasePlaceRobber` and `phaseStealCard` phases are implemented
- ✅ Board robber placement logic (`PlaceRobber` method) works
- ✅ Comprehensive tests for robber flow exist
- ✅ Rolling a 7 now triggers robber placement (updated `rollDice` in `game/turn.go:195`)
- ✅ Robber blocks resource generation on its tile (updated `GenerateResources` in `board/board.go:67`)
- ❌ Card discarding when players have >7 cards after rolling 7 is not implemented

**Implementation Details**:
- Fixed `rollDice` function to call `PhasePlaceRobber` when rolling 7
- Modified `GenerateResources` to skip tiles occupied by robber
- Added comprehensive tests in `board/robber_blocking_test.go`

**Remaining Work**: Card discarding for >7 cards requires new phase infrastructure
