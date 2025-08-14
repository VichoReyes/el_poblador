# Project Tasks

## Future Tasks

### Implement Knight Usage

**Description**: Implement the ability for players to use knight development cards both before rolling dice and during the idle period (after rolling).

**Current State**: 
- Knight option exists in the dice roll phase but is not implemented (panics with "Play Knight not implemented")
- Robber mechanics are not yet implemented

**Suggested Approach**: 
Give robber-related phases a "continuation" property (another phase) that allows the game to return to the pre-knight phase after the robbery is resolved.

**Implementation Details**:
1. **Phase Continuation**: Add a `continuation` field to phases that can be interrupted by robber actions
2. **Knight Usage Before Dice**: Allow players to play knights during the dice roll phase
3. **Knight Usage During Idle**: Allow players to play knights during the idle period
4. **Robber Integration**: When a robber is triggered (dice roll of 7), handle the robbery phase and then return to the appropriate continuation phase

**Benefits of This Approach**:
- Maintains game flow integrity
- Allows for flexible knight usage timing
- Provides clean phase transition management
- Enables future robber mechanics without breaking existing phase structure

**Technical Considerations**:
- Modify the `Phase` interface to support continuation phases
- Update existing phases to handle interruption and resumption
- Ensure proper state preservation during phase transitions
- Maintain test coverage for new phase behaviors

**Dependencies**:
- Robber mechanics implementation
- Development card system completion
- Phase transition testing framework

**Estimated Complexity**: Medium
**Priority**: High (blocks robber mechanics)
