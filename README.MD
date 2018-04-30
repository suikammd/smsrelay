# SMS Relay

Used for read messages from android phone using `adb` and send them to telegram users. 2 sim cards are supported.

Usgae: `smsrelay -token YOUR_TELEGRAM_TOKEN -chat0 CHAT_ID_0 -chat1 CHAT_ID_1`

Please be attention to set `CGO_ENABLED=1` when compiling.