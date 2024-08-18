# WEBSOCKET CONNECTION RULES
> Each websocket message contains [event, version, data], each value has stored in **byte** format.

### Event
**0: Connection**

> 0: Connection accepted

> 1: Connection refused

**1: Message**

**2: Bet**
> [type_bet, bet_id, option, amount]

`type_bet` - 0: Bet, 1: Cancel bet

`option` - index of the option

**254: Error**

**256: Disconnection**
