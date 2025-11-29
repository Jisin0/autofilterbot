package index

import (
	"fmt"
	"strings"
	"time"
)

const progressTemplate = `
%v

<b>Saved :</b>   %v
<b>Failed :</b>  %v
<b>ETA :</b>     %v
<b>PID :</b> <code>%v</code>
<b>Last Update :</b> %v

<a href='%v'><i><b>Last Indexed Message</b></i></a>
`

// buildProgressMessage builds a simple progress message for an index operation with a progress bar.
func (o *Operation) buildProgressMessage() *strings.Builder {
	var builder strings.Builder // using a builder allows external text to be added on easier

	progressBar := o.generateProgressBar()
	eta := o.calculateRemainingTime()
	now := time.Now().Format("Jan 02 15:04:05 MST")
	msgLink := fmt.Sprintf("https://t.me/c/%d/%d", o.mtprotoChannelID, o.CurrentMessageID)

	// Write to builder
	fmt.Fprintf(&builder, progressTemplate, progressBar, o.Saved, o.Failed, eta, o.ID, now, msgLink)

	return &builder
}

// calculateRemainingTime calculates the remaining time until the Estimated Time of Arrival (ETA) based on the provided parameters.
// ThankGPT!
func (o *Operation) calculateRemainingTime() string {
	var (
		total     = o.EndMessageID - o.startMessageID
		completed = o.CurrentMessageID - o.startMessageID
	)

	if completed <= 0 {
		completed = 1 // To avoid division by zero
	}

	elapsedTime := time.Since(o.startTime).Seconds()
	if elapsedTime <= 0 {
		elapsedTime = 1 // To avoid division by zero
	}

	// Calculate the rate of completion per second
	completionRate := float64(completed) / float64(elapsedTime)

	// Calculate remaining time in seconds
	remainingSeconds := float64(total-completed) / completionRate

	// Return remaining time as a duration
	remainingDuration := time.Duration(remainingSeconds) * time.Second

	// Convert remaining duration to human-readable format
	hours := int(remainingDuration.Hours())
	minutes := int(remainingDuration.Minutes()) % 60
	seconds := int(remainingDuration.Seconds()) % 60

	return fmt.Sprintf("%02d<b>h</b> %02d<b>m</b> %02d<b>s</b>", hours, minutes, seconds)
}

const (
	progressBarLength = 25 // no. of blocks in progress bar
)

// generateProgressBar generates a unicode progress bar .
// ThankGPT!.
func (o *Operation) generateProgressBar() string {
	var (
		elapsed = o.CurrentMessageID - o.StartMessageID
		total   = o.EndMessageID - o.StartMessageID
	)

	progress := float64(elapsed) / float64(total)
	barProgress := int(progress * float64(progressBarLength))

	// Build progress bar
	var progressBar strings.Builder
	progressBar.WriteString("")

	for i := 0; i < progressBarLength; i++ {
		if i < barProgress {
			progressBar.WriteString("█")
		} else {
			progressBar.WriteString("░")
		}
	}

	progressBar.WriteString("\n")

	// Add percentage
	percentage := progress * 100
	progressBar.WriteString(fmt.Sprintf("  <code>%.2f%%</code> | <code>%v</code><b>/</b><code>%v</code>", percentage, elapsed, total))

	return progressBar.String()
}
