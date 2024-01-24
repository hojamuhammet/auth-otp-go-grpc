package smpp

import (
	"auth-otp-go-grpc/internal/config"
	"log/slog"

	"github.com/fiorix/go-smpp/smpp"
	"github.com/fiorix/go-smpp/smpp/pdu/pdufield"
	"github.com/fiorix/go-smpp/smpp/pdu/pdutext"
)

type SMPPConnection struct {
	transmitter *smpp.Transmitter
	statusChan  <-chan smpp.ConnStatus
}

func NewSMPPConnection(cfg config.Config) (*SMPPConnection, error) {
	tx := &smpp.Transmitter{
		Addr:   cfg.Smpp_Address,
		User:   cfg.Smpp_User,
		Passwd: cfg.Smpp_Password,
	}

	statusChan := tx.Bind()

	for status := range statusChan {
		if status.Status() == smpp.Connected {
			slog.Info("SMPP connection established") // remove these logs if they are extra
			return &SMPPConnection{transmitter: tx, statusChan: statusChan}, nil
		}
	}

	return nil, smpp.ErrNotBound
}

func (conn *SMPPConnection) SendSMS(cfg config.Config, phoneNumber string, smsMessage string) error {
	sms := &smpp.ShortMessage{
		Src:      cfg.Smpp_Src_Phone_Number,
		Dst:      phoneNumber,
		Text:     pdutext.Raw(smsMessage), // check what provider gets: int, string or raw binary
		Register: pdufield.NoDeliveryReceipt,
	}

	_, err := conn.transmitter.Submit(sms)
	if err != nil {
		slog.Error("Failed to send SMS: %v", err)
		return err
	}

	return nil
}
