package enum

type ProgressStatus string

const (
	APPLIED      ProgressStatus = "APPLIED"
	SHORTLISTED  ProgressStatus = "SHORTLISTED"
	INTERVIEWING ProgressStatus = "INTERVIEWING"
	HIRED        ProgressStatus = "HIRED"
	REJECTED     ProgressStatus = "REJECTED"
)

// isValidProgressStatus checks if the provided status is a valid enum.ProgressStatus value.
func IsValidProgressStatus(status string) bool {
	switch ProgressStatus(status) {
	case APPLIED, SHORTLISTED, INTERVIEWING, HIRED, REJECTED:
		return true
	default:
		return false
	}
}
