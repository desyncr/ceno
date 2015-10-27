package main

import (
	"errors"
	"strings"
)

// Descriptions of the different classes of errors, converting from the kinds of names
// retrieved from JSON into something human readable.
var errorClassDescriptions map[string]string = map[string]string{
	"invalidUrl": "A URL for which a request could not be made, because it is invalid.",
	"noResponse": "A request that received no response, likely because the resource doesn't exist.",
	"malformed":  "A resource was received however it contains malformed content.",
}

/**
 * Convert the contents of a GET request for a report with certain restrictions
 * on the types of resources of interest or the classes of errors of interest
 * into the internal representation of those types, ready to go to the database.
 * @param reportMsg - The contents of the request for a report
 * @return An internal representation of the report to generate and an error if
 * values that aren't understood were supplied.
 */
func ConvertErrorReport(reportMsg ErrorReportMsg) (ErrorReport, error) {
	report := ErrorReport{-1, NoResources, NoErrorClasses, ""}
	reportErrorClasses := strings.Split(
		strings.Replace(reportMsg.ErrorClasses, " ", "", -1),
		",")
	reportResourceTypes := strings.Split(
		strings.Replace(reportMsg.ResourceTypes, " ", "", -1),
		",")
	for _, _class := range reportErrorClasses {
		errorClass, found := ErrorClasses[_class]
		if !found {
			return ErrorReport{}, errors.New("No such error class " + _class)
		} else {
			report.ErrorTypes |= errorClass
		}
	}
	for _, _type := range reportResourceTypes {
		resourceType, found := Resources[_type]
		if !found {
			return ErrorReport{}, errors.New("No such resource type " + _type)
		} else {
			report.ResourceTypes |= resourceType
		}
	}
	return report, nil
}

/**
 * Generate a textual report to describe a list of errors.
 * @param reports - An array of errors to describe
 */
func WriteReport(reports []ErrorReport) string {
	reportMsg := ""
	for _, report := range reports {
		// The first line looks like "Error concerns feeds, articles.\n"
		reportMsg += "Error concerns "
		resources := make([]string, 0)
		for resourceName, id := range Resources {
			if id&report.ResourceTypes != 0 {
				resources = append(resources, resourceName+"s")
			}
		}
		reportMsg += strings.Join(resources, ", ") + ".\n"
		// Now produce a human-readable list of error categories
		// e.g. converting "noResponse" into "A request that was not responded to"
		reportMsg += "The categories prescribed to the error are:\n"
		for className, id := range ErrorClasses {
			if id&report.ErrorTypes != 0 {
				reportMsg += "\t- " + errorClassDescriptions[className] + "\n"
			}
		}
		// Finally, append a line containing the error message.
		reportMsg += "Error message: " + report.ErrorMessage + "\n\n"
	}
	return reportMsg
}
