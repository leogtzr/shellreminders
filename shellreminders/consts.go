package shellreminders

const (
	shellReminderMainDirectory     = ".shellreminder"
	minNumberOfRecordsInFile       = 2
	shellPresenterCommand          = "toilet"
	lessThanDays                   = 8
	recordFileSeparator            = ";"
	warningRemainingDays           = 2
	redHexColor                    = "#E88388"
	yellowHexColor                 = "#DBAB79"
	greenHexColor                  = "#A8CC8C"
	errorBuildingAppConfigFromArgs = 2
	errorRunningApplication        = 1
)

var cmdArgs = [6]string{"-f", "smblock", "-w", "900", "-F", "border"}
