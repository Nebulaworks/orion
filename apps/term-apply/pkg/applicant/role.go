package applicant

func stringRole(i int) string {
	switch i {
	case 0:
		return "Senior Software Engineer"
	case 1:
		return "Software Engineer"
	default:
		return "Unknown"
	}
}
