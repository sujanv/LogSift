package compliance

import (
	"testing"

	"github.com/sujanv/LogSift/pkg/model"
)

func TestComplianceScanner(t *testing.T) {
	r1 := &model.Record{Message: "Uploaded image with credentials AKIAIOSFODNN7EXAMPLE to S3"}
	f1 := ScanRecord(1, r1)
	if len(f1) == 0 || f1[0].ViolationType != ViolationAWSAccessKey {
		t.Errorf("expected AWS key detection")
	}

	r2 := &model.Record{Message: "Generated key: -----BEGIN RSA PRIVATE KEY----- MIIEow..."}
	f2 := ScanRecord(2, r2)
	if len(f2) == 0 || f2[0].ViolationType != ViolationPrivateSSHKey {
		t.Errorf("expected private key detection")
	}

	r3 := &model.Record{Message: "Fetched patient prescription details for patient_id 9942"}
	f3 := ScanRecord(3, r3)
	if len(f3) == 0 || f3[0].ViolationType != ViolationHIPAAData {
		t.Errorf("expected HIPAA detection")
	}

	report := GenerateReport(10, append(f1, f2...))
	if report.Passed {
		t.Errorf("audit report must fail on violations")
	}
}
