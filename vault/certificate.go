package vault

import (
	"strings"
	"time"
)

type IssueCertificateRequest struct {
	CommonName string
	Ttl        time.Duration
	AltNames   []string
	IpSans     []string
}

func (r IssueCertificateRequest) Data() map[string]interface{} {
	result := map[string]interface{}{
		"common_name": r.CommonName,
		"ttl":         r.Ttl.String(),
	}

	if len(r.AltNames) > 0 {
		result["alt_names"] = strings.Join(r.AltNames, ",")
	}

	if len(r.IpSans) > 0 {
		result["ip_sans"] = strings.Join(r.IpSans, ",")
	}

	return result
}
