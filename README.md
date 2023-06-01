# Trading bot

## Strategy

![visual graph](./Drawing%202023-05-14%2017.53.48.excalidraw.png)

## Accounting Spread Sheet

[spread sheet](https://docs.google.com/spreadsheets/d/1xoO0BO94zKPtBcvD7al7xl9V2PNCDV9JpapFVZaTe0Y/edit#gid=32347275)

## Blueprint

Strategy 還沒進場時，會掛兩張 (A, B) strategy 單在現價上下 10 的位置等進場。A or B 掛到後，另一張就撤回。然後開始進行 strategy。
依照 strategy 的做法:

A. 我們會對多單掛一張 stop loss 價格在入場價格的 -0.1%。
B. 掛一張空單，其價格為 A. 的價格
C. 掛一張 B. 的 take profit 價格為 B. 的 -0.22%
D. 掛一張 B. 的 stop loss 價格為 B. 的 0.07%

Maintain a state machine that indicates the current status of the strategy

## 作法

聽 websocket 價格跟 strategy 的狀態來決定下一步的動作。


## Hook responses
### Order 被 execute 後, user data push 的 response

``` json
{
    "e": "ORDER_TRADE_UPDATE",
    "T": 1685513393253,
    "E": 1685513393256,
    "o": {
        "s": "BTCUSDT",
        "c": "ZhQI9kVGyatAXRpHy8pkxR",
        "S": "BUY",
        "o": "LIMIT",
        "f": "GTC",
        "q": "0.010",
        "p": "27521",
        "ap": "27521",
        "sp": "0",
        "x": "TRADE",
        "X": "FILLED",
        "i": 3361155118,
        "l": "0.010",
        "z": "0.010",
        "L": "27521",
        "n": "0.05504200",
        "N": "USDT",
        "T": 1685513393253,
        "t": 260496182,
        "b": "0",
        "a": "0",
        "m": true,
        "R": false,
        "wt": "CONTRACT_PRICE",
        "ot": "LIMIT",
        "ps": "BOTH",
        "cp": false,
        "rp": "0",
        "pP": false,
        "si": 0,
        "ss": 0
    }
}
```

### Long order response hit take profit
```json
{
    "e": "ORDER_TRADE_UPDATE",
    "T": 1685598416696,
    "E": 1685598416700,
    "o": {
        "s": "BTCUSDT",
        "c": "BQpueAZYJLqzuqycZemlyf",
        "S": "SELL",
        "o": "MARKET",
        "f": "GTC",
        "q": "0.010",
        "p": "0",
        "ap": "26838",
        "sp": "26906.10",
        "x": "TRADE",
        "X": "FILLED",
        "i": 3362131502,
        "l": "0.010",
        "z": "0.010",
        "L": "26838",
        "n": "0.10735200",
        "N": "USDT",
        "T": 1685598416696,
        "t": 260544265,
        "b": "0",
        "a": "0",
        "m": false,
        "R": true,
        "wt": "CONTRACT_PRICE",
        "ot": "TAKE_PROFIT_MARKET",
        "ps": "LONG",
        "cp": false,
        "rp": "-0.00999999",
        "pP": false,
        "si": 0,
        "ss": 0
    }
}
```

## Account credentials

short account.

`justinhsu90@gmail.com`
`Jh8918206!`

## Transaction report format

```json
{
    "strategy_start_time": ...
    "order_id": ...
    "status": "HIT_LONG_TAKE_PROFIT"
    "realized_profit": ...
    "strategy_end_time": ...
}
```

## Misc

https://dev.binance.vision/t/error-orders-position-side-does-not-match-users-setting/5970/2