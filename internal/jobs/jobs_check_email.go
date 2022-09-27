package jobs

import (
	"context"

	"github.com/google/uuid"
	"github.com/jordanknott/monitor/internal/utils"
	"github.com/matcornic/hermes/v2"
)

func (t *JobTasks) CheckInstallEmailReport(reportIDEncoded string) (bool, error) {
	reportID, err := uuid.Parse(reportIDEncoded)
	if err != nil {
		return false, err
	}
	ctx := context.Background()
	installID, err := t.Data.GetInstallIDForReportID(ctx, reportID)
	if err != nil {
		return false, err
	}
	install, err := t.Data.GetInstallByID(ctx, installID)
	if err != nil {
		return false, err
	}

	data := [][]hermes.Entry{}
	entries, err := t.Data.GetReportEntriesForReport(ctx, reportID)
	if err != nil {
		return false, err
	}
	if len(entries) == 0 {
		return true, nil
	}
	for _, entry := range entries {
		data = append(data, []hermes.Entry{
			{
				Key:   "Filepath",
				Value: entry.Filepath,
			},
		})
	}

	h := hermes.Hermes{
		Product: hermes.Product{
			Name: "Monitor",
			Link: t.AppConfig.Email.SiteURL,
			Logo: "https://digitalmarketingformanufacturers.com/wp-content/uploads/2022/03/DD-logo-cropped.svg",
		},
	}

	email := hermes.Email{
		Body: hermes.Body{
			Name: "Jordan Knott",
			Intros: []string{
				"Some files have changed on the install " + install.Nicename + "!",
			},
			Table: hermes.Table{
				Data: data,
			},
			Actions: []hermes.Action{
				{
					Instructions: "To view the report, click here",
					Button: hermes.Button{
						Color:     "#7367F0", // Optional action button color
						TextColor: "#FFFFFF",
						Text:      "View report",
						Link:      t.AppConfig.Email.SiteURL + "/",
					},
				},
			},
			Outros: []string{
				"Need help, or have questions? Just reply to this email, we'd love to help.",
			},
		},
	}

	emailBody, err := h.GenerateHTML(email)
	if err != nil {
		return false, err
	}
	_, err = h.GeneratePlainText(email)
	if err != nil {
		return false, err
	}
	err = utils.SendMail(t.AppConfig.Email, utils.Email{
		To:   "jordan@drivendigital.us",
		HTML: emailBody,
	}, "["+install.Nicename+"] Some files have changed!")
	if err != nil {
		return false, err
	}
	return true, nil
}
