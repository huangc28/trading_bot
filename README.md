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


## Misc

https://dev.binance.vision/t/error-orders-position-side-does-not-match-users-setting/5970/2