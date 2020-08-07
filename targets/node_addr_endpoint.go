package targets

import (
	"chainflow-vitwit/config"
	"fmt"
	"os/exec"

	client "github.com/influxdata/influxdb1-client/v2"
)

// CheckMeled function to run the command get meled status and send
//alerts to telgram and email accounts
func CheckMeled(_ HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return
	}

	cmd := exec.Command("bash", "-c", "</dev/tcp/0.0.0.0/26656 &>/dev/null")
	out, err := cmd.CombinedOutput()
	if err != nil {
		_ = SendTelegramAlert("Meled on your validator instance is not running", cfg)
		_ = SendEmailAlert("Meled on your validator instance is not running", cfg)
		_ = writeToInfluxDb(c, bp, "vcf_meled_status", map[string]string{}, map[string]interface{}{"status": 0})
		return
	}

	resp := string(out)
	if resp != "" {
		_ = SendTelegramAlert(fmt.Sprintf("Meled on your validator instance is not running: \n%v", resp), cfg)
		_ = SendEmailAlert(fmt.Sprintf("Meled on your validator instance is not running: \n%v", resp), cfg)
		_ = writeToInfluxDb(c, bp, "vcf_meled_status", map[string]string{}, map[string]interface{}{"status": 0})
		return
	}

	_ = writeToInfluxDb(c, bp, "vcf_meled_status", map[string]string{}, map[string]interface{}{"status": 1})
}
