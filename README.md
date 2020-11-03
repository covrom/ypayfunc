# ypayfunc

The Yoomoney (Yandex Money) payment notification function for [Yandex cloud functions](https://cloud.yandex.ru/services/functions) accepts incoming payments from a credit card or Yandex money and then sends an email notifying that payment. You can change this behavior to whatever you want.

Quick start:

1. See https://yandex.ru/dev/money/apps/?_openstat=settings;other;apps;api for basic instructions.
2. Create public Yandex cloud function here: https://console.cloud.yandex.ru
3. Create archive file: `zip -r archive.zip *.go go.mod go.sum`
4. Upload archive into Yandex cloud function and create new version with:
- runtime golang114
- entrypoint main.YaPay
- memory 128m
- execution-timeout 20s
5. Set environment vars:
  - MSRV - server and port of SSL mailserver, i.e. smtp.yandex.ru:465
  - MLGN - mailserver login
  - MPSW - mailserver password
  - MLCC - your admin email for send copies of notifications
  - PUSHPSW - your secret from https://yoomoney.ru/transfer/myservices/http-notification
6. Press testing button at https://yoomoney.ru/transfer/myservices/http-notification and see logs of function in Yandex cloud
7. Create and run html file, using `testpay_example.html` template. Ensure to change your identity in `receiver` field and fill `sum` field