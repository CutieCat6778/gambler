# WEBSOCKET CONNECTION RULES
> Each websocket message contains [event, version, data], each value has stored in **byte** format.

### Event
**0: Connection**

> 0: Connection accepted

> 1: Connection refused

**1: Message**

**254: Error**

**256: Disconnection**
