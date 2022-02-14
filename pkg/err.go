// Copyright 2021-2022 myl7
// SPDX-License-Identifier: Apache-2.0

package pkg

import (
	"errors"
	"fmt"
	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
	"github.com/myl7/tgchan2tw/pkg/cfg"
	"log"
	"strings"
)

func reportErr(errMsg string) []error {
	var errs []error

	log.Print(errMsg)

	if cfg.Cfg.OnErrEmail != "" && !cfg.Cfg.DisableSmtp {
		err := reportErrViaEmail(errMsg)
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}

func reportErrViaEmail(errMsg string) error {
	auth := sasl.NewPlainClient("", cfg.Cfg.SmtpUsername, cfg.Cfg.SmtpPassword)
	to := []string{cfg.Cfg.OnErrEmail}
	msg := strings.NewReader(fmt.Sprintf("To: %s\r\n", cfg.Cfg.OnErrEmail) +
		fmt.Sprintf("Subject: Error of tgchan2tw instance %s\r\n", cfg.Cfg.Name) +
		"\r\n" +
		fmt.Sprintf("%s\r\n", errMsg))
	err := smtp.SendMail(fmt.Sprintf("%s:%d", cfg.Cfg.SmtpHost, cfg.Cfg.SmtpPort), auth, cfg.Cfg.SmtpSender, to, msg)
	if err != nil {
		return errors.New(fmt.Sprint("failed to report error via email:", err))
	}
	return nil
}
