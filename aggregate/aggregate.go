package aggregate

// CPU_CreditUsage calculates CPU credits for a process. We use AWS formula.
// See https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/burstable-credits-baseline-concepts.html
// func CPU_CreditUsage(process CPU_Usage) (credits int) {
// 	return int(process.PCPU) * process.Duration
// }
