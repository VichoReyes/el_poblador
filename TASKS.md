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

### ✅ Development Card Implementations - COMPLETED

**Description**: Complete implementation of remaining development cards (Monopoly, Year of Plenty, Road Building). All development cards are now fully implemented.

**Final State**: 
- ✅ Knight card fully implemented (robber placement + stealing)
- ✅ Victory Point card structure exists (no action needed - passive VP boost)
- ✅ Monopoly card - COMPLETED - player selects resource type, collects all of that type from other players
- ✅ Year of Plenty card - COMPLETED - player selects 2 resources from bank
- ✅ Road Building card - COMPLETED - player places 2 roads for free

**Implementation Details**:
- **Monopoly Phase**: `phaseMonopoly` in `game/turn.go:669-709` - resource selection menu, steals all selected resources from other players
- **Year of Plenty Phase**: `phaseYearOfPlenty` in `game/turn.go:711-767` - two-step resource selection from bank
- **Road Building Phase**: Refactored to reuse existing `phaseRoadStart/End` with new optional parameters (`isFree`, `continuation`, `helpPrefix`) - eliminates code duplication

**Key Fixes**: 
1. Fixed critical pointer issue in `PhasePlayDevelopmentCard.Confirm()` - changed from `player := p.game.players[p.game.playerTurn]` to `player := &p.game.players[p.game.playerTurn]` to ensure card state changes persist.
2. **Road Building Refactor**: Instead of duplicating road building logic, enhanced existing phases with clean public APIs:
   - `PhaseRoadStart(game, previousPhase)` - for regular road building (pays resources)
   - `PhaseRoadBuilding(game)` - for dev card road building (free, places 2 roads)
   - Internal implementation uses shared `newPhaseRoadStart()` and `newPhaseRoadEnd()` functions
   - Free road building: skips resource cost and chains first → second road → completion
   - Help text customization: shows "first free road" vs "second free road"

**Tests Added**:
- `TestMonopolyCard` - verifies resource collection from other players
- `TestYearOfPlentyCard` - verifies two-resource selection from bank  
- `TestRoadBuildingCard` - verifies card playing and phase transition

All development card functionality is now complete and tested.
