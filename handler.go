package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"net/url"
	"os"
	"reflect"
	"strings"
)

// https://tech.yandex.ru/money/doc/dg/reference/notification-p2p-incoming-docpage/
// https://yoomoney.ru/transfer/myservices/http-notification
// https://yandex.ru/dev/money/apps/?_openstat=settings%3Bother%3Bapps%3Bapi

// При получении уведомления всегда проверяйте статус входящего перевода по значениям полей unaccepted и codepro.
// Если unaccepted=true, то перевод еще не зачислен на счет получателя.
// Чтобы его принять, получателю нужно совершить дополнительные действия.
// Например, освободить место на счете, если достигнут лимит доступного остатка.
// Или указать код протекции, если он необходим для получения перевода.
// Если codepro=true, то перевод защищен кодом протекции.
// Чтобы получить такой перевод, пользователю необходимо ввести код протекции.

func YandexMoneyIncomingPush(rbody, yaKey string) {
	vs, err := url.ParseQuery(rbody)
	if err != nil {
		fmt.Printf("YandexMoneyIncomingPush ParseQuery error: %s\n", err)
		// в случае ошибок здесь всегда возвращаем 200 ОК
		return
	}

	yap := &YaParams{}
	yap.ParsePostForm(vs)

	fmt.Printf("Received payment notification: %+v\n", yap)

	if err := yap.CheckSha1(yaKey); err != nil {
		fmt.Printf("CheckSha1 error: %s\n", err)
		return
	}

	m := NewMailer(
		os.Getenv("MSRV"),
		os.Getenv("MLGN"),
		os.Getenv("MPSW"),
	)

	sendmail(yap.Email, yap.Label, yap.WithDrawAmount, m)
}

func sendmail(eml, invid, amnt string, mlr *Mailer) {
	if eml != "" {
		if err := mlr.Send(eml, "Получен перевод средств по счету "+invid,
			fmt.Sprintf(`Получен перевод средств по счету %s на сумму %s`,
				invid, amnt)); err != nil {
			fmt.Printf("Sending mail for %q about payment for invoice %q error: %s\n", eml, invid, err)
			return
		}
	}
	if err := mlr.Send(os.Getenv("MLCC"), "Получен перевод средств по счету "+invid,
		fmt.Sprintf(`Получен перевод средств от имени %s по счету %s на сумму %s`,
			eml, invid, amnt)); err != nil {
		fmt.Printf("Sending mail for admin about payment for invoice %q error: %s\n", invid, err)
		return
	}
}

type YaParams struct {
	NotificationType string `form:"notification_type"`
	OperationId      string `form:"operation_id"`
	Amount           string `form:"amount"`
	Currency         string `form:"currency"`
	Datetime         string `form:"datetime"`
	Sender           string `form:"sender"`
	Codepro          string `form:"codepro"`
	Label            string `form:"label"`

	WithDrawAmount string `form:"withdraw_amount"`
	Sha1Hash       string `form:"sha1_hash"`

	Email       string `form:"email"`
	Lastname    string `form:"lastname"`
	Firstname   string `form:"firstname"`
	Fathersname string `form:"fathersname"`
}

func (yp *YaParams) ParsePostForm(data url.Values) {
	typ := reflect.TypeOf(yp).Elem()
	val := reflect.ValueOf(yp).Elem()

	for i := 0; i < typ.NumField(); i++ {
		typeField := typ.Field(i)
		structField := val.Field(i)
		if structField.CanSet() {
			inputFieldName := typeField.Tag.Get("form")
			if inputFieldName != "" {
				structField.SetString(data.Get(inputFieldName))
			}
		}
	}
}

func (yp *YaParams) CheckSha1(nsecret string) error {
	str := yp.NotificationType + "&" +
		yp.OperationId + "&" +
		yp.Amount + "&" +
		yp.Currency + "&" +
		yp.Datetime + "&" +
		yp.Sender + "&" +
		yp.Codepro + "&" +
		nsecret + "&" +
		yp.Label

	hasher := sha1.New()
	hasher.Write([]byte(str))
	sha := hex.EncodeToString(hasher.Sum(nil))

	if !strings.EqualFold(sha, yp.Sha1Hash) {
		return fmt.Errorf("not equal sign hash: %q != %q", sha, yp.Sha1Hash)
	}
	return nil
}
