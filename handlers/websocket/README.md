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

**3: GetBetWinRate**
> [bet_id, amount_a, amount_b, input]

`amount_a` - first value after comma

`amount_b` - second value after comma

`input` - input value

**254: Error**

**256: Disconnection**
