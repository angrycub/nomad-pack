# Output from the API client

## `nomadMembers` ([]*api.AgentMember as pretty JSON)

```json
[[ range nomadMembers ]][[ toPrettyJson . ]][[ end ]]
```

## `nomadIsEnt` (bool)

nomadIsEnt: [[with nomadIsEnt]][[.]][[end]]

### Conditional use

We love our [[ if nomadIsEnt ]]enterprise customers[[ else ]]community members[[ end ]].

## `nomadAPIClient`

[[ $c := nomadAPIClient ]]
[[ $c.Agent.Self | toPrettyJson]]
[[/* */]]
